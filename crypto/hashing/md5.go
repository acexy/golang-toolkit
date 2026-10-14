package hashing

import (
	"crypto/md5"

	"github.com/acexy/golang-toolkit/math/conversion"
)

// Md5Hex returns the MD5 digest of a string as lowercase hex.
func Md5Hex(data string) string {
	hash := Md5Bytes(conversion.ParseBytes(data))
	return hexString(hash[:])
}

// Md5Base64 returns the MD5 digest of a string as Base64.
func Md5Base64(data string) string {
	hash := Md5Bytes(conversion.ParseBytes(data))
	return base64String(hash[:])
}

// Md5Bytes returns the MD5 digest of bytes.
func Md5Bytes(data []byte) [md5.Size]byte {
	return md5.Sum(data)
}

// Md5FileHex streams a file and returns its MD5 digest as lowercase hex.
func Md5FileHex(absFilePath string) (string, error) {
	hash, err := hashFile(md5.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return hexString(hash), nil
}

// Md5FileBase64 streams a file and returns its MD5 digest as Base64.
func Md5FileBase64(absFilePath string) (string, error) {
	hash, err := hashFile(md5.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return base64String(hash), nil
}
