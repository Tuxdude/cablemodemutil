package cablemodemutil

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"io"
)

// Generates HMAC-MD5 using the specified key and message strings.
func genHMACMD5(key string, msg string) string {
	h := hmac.New(md5.New, []byte(key))
	io.WriteString(h, msg)
	return fmt.Sprintf("%X", h.Sum(nil))
}
