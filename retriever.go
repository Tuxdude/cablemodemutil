package cablemodemutil

import (
	"fmt"
)

// Retriever is used to retrieve the current status of the Cable Modem.
type Retriever struct {
	username      string
	clearPassword string
}

// RetrieverInput is used to specify the input for building a Retriever.
type RetrieverInput struct {
	// User name for authenticating with the cable modem.
	Username string
	// Password for authenticating with the cable modem.
	ClearPassword string
}

// Builds a new Retriever object to query the Cable Modem.
func NewStatusRetriever(input *RetrieverInput) *Retriever {
	r := Retriever{}
	r.username = input.Username
	r.clearPassword = input.ClearPassword
	return &r
}

// Retrieves the current detailed raw status from the cable modem.
func (r *Retriever) RawStatus() (CableModemRawStatus, error) {
	return nil, fmt.Errorf("unimplemented!")
}
