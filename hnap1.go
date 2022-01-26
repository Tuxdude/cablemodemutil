package cablemodemutil

import (
	"fmt"
)

const (
	soapNamespace = "http://purenetworks.com/HNAP1"
)

// Generates the SOAP Action URI.
func actionURI(action string) string {
	return fmt.Sprintf("\"%s/%s\"", soapNamespace, action)
}
