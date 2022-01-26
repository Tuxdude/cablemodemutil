package cablemodemutil

import (
	"fmt"
)

const (
	urlFormat = "%s://%s/HNAP1/"
)

// Retriever is used to retrieve the current status of the Cable Modem.
type Retriever struct {
	client        *httpClient
	username      string
	clearPassword string
}

// RetrieverInput is used to specify the input for building a Retriever.
type RetrieverInput struct {
	// The host name or IP address of the cable modem device.
	Host string
	// The protocol used to connect to the cable modem, either "http" or "https".
	Protocol string
	// If true skips verifying the cable modem's SSL certificate, false otherwise.
	SkipVerifyCert bool
	// User name for authenticating with the cable modem.
	Username string
	// Password for authenticating with the cable modem.
	ClearPassword string
}

// Builds a new Retriever object to query the Cable Modem.
func NewStatusRetriever(input *RetrieverInput) *Retriever {
	url := fmt.Sprintf(urlFormat, input.Protocol, input.Host)
	r := Retriever{}
	r.client = newHttpClient(url, input.SkipVerifyCert)
	r.username = input.Username
	r.clearPassword = input.ClearPassword
	return &r
}

// Retrieves the current detailed raw status from the cable modem.
func (r *Retriever) RawStatus() (CableModemRawStatus, error) {
	return nil, fmt.Errorf("unimplemented!")
}
