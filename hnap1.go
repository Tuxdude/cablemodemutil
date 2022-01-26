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

// Generates the key which contains the response for the SOAP action in the payload.
func actionResponseKey(action string) string {
	return fmt.Sprintf("%sResponse", action)
}

// Generates the key which contains the result of the SOAP action in the payload.
func actionResultKey(action string) string {
	return fmt.Sprintf("%sResult", action)
}

// Generates the private key using the public key, challenge and the clear password.
func genPrivateKey(publicKey string, challenge string, clearPassword string) string {
	return genHMACMD5(publicKey+clearPassword, challenge)
}

// Generates the hashed password using the private key, challenge and the clear password.
func genHashedPassword(privateKey string, challenge string, clearPassword string) string {
	return genHMACMD5(privateKey, challenge)
}
