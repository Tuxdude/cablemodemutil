package cablemodemutil

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	connectionTimeout = 15
)

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
