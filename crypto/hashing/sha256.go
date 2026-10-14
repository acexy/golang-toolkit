package hashing

import (
	"crypto/sha256"

	"github.com/acexy/golang-toolkit/math/conversion"
)

// Sha256Hex returns the SHA256 digest of a string as lowercase hex.
func Sha256Hex(data string) string {
	return hexString(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Base64 returns the SHA256 digest of a string as Base64.
func Sha256Base64(data string) string {
	return base64String(Sha256Bytes(conversion.ParseBytes(data)))
}

// Sha256Bytes returns the SHA256 digest of bytes.
func Sha256Bytes(data []byte) []byte {
	return hashBytes(sha256.New(), data)
}

// Sha256FileHex streams a file and returns its SHA256 digest as lowercase hex.
func Sha256FileHex(absFilePath string) (string, error) {
	hash, err := hashFile(sha256.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return hexString(hash), nil
}

// Sha256FileBase64 streams a file and returns its SHA256 digest as Base64.
func Sha256FileBase64(absFilePath string) (string, error) {
	hash, err := hashFile(sha256.New(), absFilePath)
	if err != nil {
		return "", err
	}
	return base64String(hash), nil
}
