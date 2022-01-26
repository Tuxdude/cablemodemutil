package cablemodemutil

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	actionHeader           = "SOAPAction"
	hnapAuthHeader         = "HNAP_AUTH"
	connectionTimeout      = 15
	contentTypeHeader      = "Content-Type"
	contentTypeHeaderValue = "application/json; charset=UTF-8"
)

// Generate the list of cookies to be included in the request.
func genCookies(tok *token) []*http.Cookie {
	if tok.uid == "" {
		// If uid is empty, it means we have not yet received a valid
		// uid and the private key is the dummy initial value 'withoutLoginKey'.
		return nil
	}

	return []*http.Cookie{
		&http.Cookie{
			Name:   "uid",
			Value:  tok.uid,
			Path:   "/",
			Secure: true,
		},
		&http.Cookie{
			Name:   "PrivateKey",
			Value:  tok.privateKey,
			Path:   "/",
			Secure: true,
		},
	}
}

// Adds the necessary headers and cookies to the HTTP request for the specified SOAP action.
func addHeadersAndCookies(req *http.Request, action string, tok *token) {
	req.Header.Add(contentTypeHeader, contentTypeHeaderValue)
	req.Header.Add(actionHeader, actionURI(action))
	req.Header.Add(hnapAuthHeader, genHNAPAuth(tok.privateKey, action))
	c := genCookies(tok)
	for _, cookie := range c {
		req.AddCookie(cookie)
	}
}

type httpClient struct {
	client *http.Client
	url    string
}

func newHttpClient(url string, skipVerifyCert bool) *httpClient {
	c := httpClient{}
	c.client = &http.Client{
		Timeout: connectionTimeout * time.Second,
	}
	if skipVerifyCert {
		c.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	c.url = url
	return &c
}

// Sends the HTTP POST request for the specified SOAP action containing the specified payload.
func (c *httpClient) sendPOST(action string, payload *bytes.Buffer, tok *token) (*[]byte, error) {
	req, err := http.NewRequest("POST", c.url, payload)
	if err != nil {
		return nil, err
	}
	// Needed to avoid EOF errors.
	// See https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
	req.Close = true

	addHeadersAndCookies(req, action, tok)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed.\nresp: %s\nreason: %w", action, prettyPrintJSON(resp), err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed while reading the response body.\nreason: %w", action, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed due to non-success status code: %d\nbody:%s", action, resp.StatusCode, string(body))
	}

	return &body, nil
}
