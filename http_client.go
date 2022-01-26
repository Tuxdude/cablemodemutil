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
	connectionTimeout      = 15
	contentTypeHeader      = "Content-Type"
	contentTypeHeaderValue = "application/json; charset=UTF-8"
)

// Adds the necessary headers to the HTTP request for the specified SOAP action.
func addHeaders(req *http.Request, action string) {
	req.Header.Add(contentTypeHeader, contentTypeHeaderValue)
	req.Header.Add(actionHeader, actionURI(action))
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
func (c *httpClient) sendPOST(action string, payload *bytes.Buffer) (*[]byte, error) {
	req, err := http.NewRequest("POST", c.url, payload)
	if err != nil {
		return nil, err
	}

	addHeaders(req, action)

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
