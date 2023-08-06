package md5

import (
	sdkmd5 "crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// Md5Hex Md5字符串返回HexString
func Md5Hex(data string) string {
	hash := Md5([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Md5 对字节做Md5
func Md5(bytes []byte) [sdkmd5.Size]byte {
	return sdkmd5.Sum(bytes)
}

// Md5FileHex 计算文件md5返回Hex
func Md5FileHex(absFilePath string) (string, error) {
	file, err := os.Open(absFilePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	hash := sdkmd5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
