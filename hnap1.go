package cablemodemutil

import (
	"fmt"
	"time"
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
func genPrivateKey(publicKey string, challenge string, clearPassword string) (string, error) {
	return genHMACMD5(publicKey+clearPassword, challenge)
}

// Generates the hashed password using the private key and challenge.
func genHashedPassword(privateKey string, challenge string) (string, error) {
	return genHMACMD5(privateKey, challenge)
}

// Generates the HNAP auth for the request.
func genHNAPAuth(privateKey string, soapAction string) (string, error) {
	currTime := time.Now().UnixMilli()
	authMsg := fmt.Sprintf("%d%s", currTime, actionURI(soapAction))
	auth, err := genHMACMD5(privateKey, authMsg)
	if err != nil {
		return "", fmt.Errorf("HNAP auth generation failed, reason: %w", err)
	}
	return fmt.Sprintf("%s %d", auth, currTime), nil
}
