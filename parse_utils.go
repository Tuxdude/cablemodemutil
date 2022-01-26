package cablemodemutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Parses the specified string as an uint64 after stripping the suffix if required.
func parseUint64(str string, hasSuffix bool, suffix string, desc string) (uint64, error) {
	if hasSuffix {
		if !strings.HasSuffix(str, suffix) {
			return 0, fmt.Errorf("expected %s with %q suffix, but not available in %q", desc, suffix, str)
		}
		str = strings.TrimSuffix(str, suffix)
	}

	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %q to uint64: %w", str, err)
	}
	return res, nil
}

// Parses the specified string as a channel frequency after stripping the ' Hz' suffix if required.
func parseFreqStr(str string, hasHzSuffix bool, desc string) (uint32, error) {
	f, err := parseUint64(str, hasHzSuffix, " Hz", desc)
	if err != nil {
		return 0, err
	}
	return uint32(f), nil
}

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

// Parses the value of the specified key as a channel frequency in the specified status information.
func parseFreq(data actionResponseBody, key string, hasHzSuffix bool, desc string) (uint32, error) {
	s, err := parseString(data, key, desc)
	if err != nil {
		return 0, err
	}
	return parseFreqStr(s, hasHzSuffix, desc)
}
