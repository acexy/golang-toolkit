package random

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersLen = 62

func RandString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(lettersLen-1)]
	}
	return string(result)
}
