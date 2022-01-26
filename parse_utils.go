package cablemodemutil

import (
	"fmt"
)

// Parses the value of the specified key as a string in the specified status information.
func parseString(data actionResponseBody, key string, desc string) (string, error) {
	s, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("unable to find key %q for parsing %q.data=%v", key, desc, data)
	}
	return s, nil
}

// Parses the value of the specified key as a bool in the specified status information.
func parseBool(data actionResponseBody, key string, trueVal string, desc string) (bool, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return false, err
	}
	return s == trueVal, nil
}
