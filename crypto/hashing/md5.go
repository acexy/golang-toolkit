package hashing

import (
	"crypto/md5"

	"github.com/acexy/golang-toolkit/math/conversion"
)

// Md5Hex 对字符串执行 MD5 运算，并返回 Hex 字符串
func Md5Hex(data string) string {
	hash := Md5Bytes(conversion.ParseBytes(data))
	return hexString(hash[:])
}

// Md5Base64 对字符串执行 MD5 运算，并返回 Base64 字符串
func Md5Base64(data string) string {
	hash := Md5Bytes(conversion.ParseBytes(data))
	return base64String(hash[:])
}

// Md5Bytes 对字节数组执行 MD5 运算
func Md5Bytes(data []byte) [md5.Size]byte {
	return md5.Sum(data)
}

// Md5FileHex 计算文件 MD5，并返回 Hex 字符串
func Md5FileHex(absFilePath string) (string, error) {
	hash, err := hashFile(md5.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return hexString(hash), nil
}

// Md5FileBase64 计算文件 MD5，并返回 Base64 字符串
func Md5FileBase64(absFilePath string) (string, error) {
	hash, err := hashFile(md5.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return base64String(hash), nil
}
