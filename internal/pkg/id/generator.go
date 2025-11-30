package id

import (
	"crypto/rand"
	"encoding/hex"
)

// New generates a random identifier.
func New() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
