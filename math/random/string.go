package random

import (
	"github.com/google/uuid"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersLen = 62

// RandString 生成指定长度的随机字符串
func RandString(length int) string {
	result := make([]byte, length)
	for i := range result {
		index := RandInt(lettersLen - 1)
		result[i] = letters[index]
	}
	return string(result)
}

// UUID 生成UUID
func UUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
