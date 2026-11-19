package random

import (
	"encoding/hex"
	"math/rand/v2"
	"uuid"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersLen = 62

// RandString returns a pseudo-random alphanumeric string, or empty for nonpositive lengths. It is not a secret generator.
func RandString(length int) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := range result {
		index := rand.IntN(lettersLen)
		result[i] = letters[index]
	}
	return string(result)
}

// UUID returns a UUID v4 as 32 lowercase hexadecimal characters without hyphens.
func UUID() string {
	id := uuid.NewV4()
	return hex.EncodeToString(id[:])
}
