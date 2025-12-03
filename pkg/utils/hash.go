package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// ToSha256 returns the sha256 hash of a given string
func ToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
