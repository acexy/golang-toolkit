package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/acexy/golang-toolkit/math/conversion"
)

// HexSha256 对值得字符串执行 hash256 运算
func HexSha256(data string) string {
	return hex.EncodeToString(BytesSha256(conversion.ParseBytes(data)))
}

// BytesSha256 对 byte 数组执行 hash256 运算
func BytesSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
