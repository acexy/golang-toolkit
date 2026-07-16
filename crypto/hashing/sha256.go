package hashing

import (
	"crypto/sha256"

	"github.com/acexy/golang-toolkit/math/conversion"
)

// Sha256Hex 对字符串执行 SHA256 运算，并返回 Hex 字符串
func Sha256Hex(data string) string {
	return hexString(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Base64 对字符串执行 SHA256 运算，并返回 Base64 字符串
func Sha256Base64(data string) string {
	return base64String(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Bytes 对字节数组执行 SHA256 运算
func Sha256Bytes(data []byte) []byte {
	return hashBytes(sha256.New(), data)
}

// Sha256FileHex 计算文件 SHA256，并返回 Hex 字符串
func Sha256FileHex(absFilePath string) (string, error) {
	hash, err := hashFile(sha256.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return hexString(hash), nil
}

// Sha256FileBase64 计算文件 SHA256，并返回 Base64 字符串
func Sha256FileBase64(absFilePath string) (string, error) {
	hash, err := hashFile(sha256.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return base64String(hash), nil
}
