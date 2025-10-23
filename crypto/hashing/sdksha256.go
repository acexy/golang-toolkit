package hashing

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/acexy/golang-toolkit/math/conversion"
)

// Sha256Hex 返回Hex结果
func Sha256Hex(data string) string {
	return hex.EncodeToString(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Base64 返回Base64结果
func Sha256Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Bytes 对 byte 数组执行 hash256 运算
func Sha256Bytes(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
