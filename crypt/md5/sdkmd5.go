package md5

import (
	sdkmd5 "crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// HexMd5 Md5字符串返回HexString
func HexMd5(data string) string {
	hash := BytesMd5([]byte(data))
	return hex.EncodeToString(hash[:])
}

// BytesMd5 对字节做Md5
func BytesMd5(bytes []byte) [sdkmd5.Size]byte {
	return sdkmd5.Sum(bytes)
}

// FileHexMd5 计算文件md5返回Hex
func FileHexMd5(absFilePath string) (string, error) {
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
