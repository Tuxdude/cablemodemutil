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
		{
			Name:   "uid",
			Value:  tok.uid,
			Path:   "/",
			Secure: true,
		},
		{
			Name:   "PrivateKey",
			Value:  tok.privateKey,
			Path:   "/",
			Secure: true,
		},
	}
}

// Adds the necessary headers and cookies to the HTTP request for the specified SOAP action.
func addHeadersAndCookies(req *http.Request, action string, tok *token) error {
	auth, err := genHNAPAuth(tok.privateKey, action)
	if err != nil {
		return err
	}

	req.Header.Add(contentTypeHeader, contentTypeHeaderValue)
	req.Header.Add(actionHeader, actionURI(action))
	req.Header.Add(hnapAuthHeader, auth)
	c := genCookies(tok)
	for _, cookie := range c {
		req.AddCookie(cookie)
	}
	return nil
}

type httpClient struct {
	client *http.Client
	url    string
	debug  RetrieverDebug
}

func newHTTPClient(url string, skipVerifyCert bool, debug *RetrieverDebug) *httpClient {
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
	c.debug = *debug
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

	err = addHeadersAndCookies(req, action, tok)
	if err != nil {
		return nil, err
	}

	if c.debug.Debug {
		fmt.Printf("Dumping token before the request: %s\n", action)
		debugToken(tok)
	}

	if c.debug.DebugReq {
		debugHTTPRequest(req)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed.\nresp: %s\nreason: %w", action, prettyPrintJSON(resp), err)
	}
	defer resp.Body.Close()
	if c.debug.DebugResp {
		debugHTTPResponse(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed while reading the response body.\nreason: %w", action, err)
	}

	if resp.StatusCode != 200 {
		// TODO: Convert this case into a specific error type that allows the
		// caller to retry explicitly.
		if resp.StatusCode == 404 && action == queryAction {
			return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed with 404 status code possibly due to credentials having expired.\nbody:%s", action, string(body))
		}
		return nil, fmt.Errorf("HTTP POST request for SOAP action %q failed due to non-success status code: %d\nbody:%s", action, resp.StatusCode, string(body))
	}

	return &body, nil
}
