package cablemodemutil

import (
	"crypto/hmac"
	"crypto/md5" // nolint:gosec
	"fmt"
	"io"
)

// Generates HMAC-MD5 using the specified key and message strings.
func genHMACMD5(key string, msg string) (string, error) {
	h := hmac.New(md5.New, []byte(key))
	_, err := io.WriteString(h, msg)
	if err != nil {
		return "", fmt.Errorf("HMAC MD5 generation failed, reason: %w", err)
	}
	return fmt.Sprintf("%X", h.Sum(nil)), nil
}
