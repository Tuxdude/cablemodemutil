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
