package cablemodemutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// Logs token debug information.
func debugToken(tok *token) {
	fmt.Printf("Current Time: %q\n", time.Now())
	fmt.Printf("Token: {expiry:%q uid:%q privateKey:%q}\n\n\n", tok.expiry, tok.uid, tok.privateKey)
}

// Dumps the HTTP request for the purpose of debugging.
func debugHTTPRequest(req *http.Request) {
	writeDebugOutput(httputil.DumpRequestOut(req, true))
}

// Dumps the HTTP response for the purpose of debugging.
func debugHTTPResponse(resp *http.Response) {
	writeDebugOutput(httputil.DumpResponse(resp, true))
}

// Logs the specified HTTP payload.
func writeDebugOutput(data []byte, err error) {
	if err != nil {
		log.Fatalf("%s\n\n", err)
	}
	fmt.Printf("%s\n\n", data)
}

// Returns the JSON formatted string representation of the specified object.
func prettyPrintJSON(x interface{}) string {
	p, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", x)
	}
	return string(p)
}
