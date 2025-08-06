// File: /quicklynks/backend/internal/utils/random.go
package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// StringWithCharset generates a random string of a given length from a charset.
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateShortCode creates a random string to be used as a short URL code.
func GenerateShortCode() string {
	// Using 7 characters gives 62^7 possible combinations, which is over 3.5 trillion.
	// This is sufficient to avoid collisions for a very long time.
	return StringWithCharset(7, charset)
}
