package random

import (
	"math/rand/v2"
	"strings"

	"github.com/google/uuid"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const lettersLen = 62

// RandString 生成指定长度的伪随机字符串，不适用于密钥、令牌等安全场景
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

// UUID 生成UUID
func UUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
