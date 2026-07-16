package hashing

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func hashBytes(hashFunc hash.Hash, data []byte) []byte {
	hashFunc.Write(data)
	return hashFunc.Sum(nil)
}

func hashFile(hashFunc hash.Hash, absFilePath string) ([]byte, error) {
	file, err := os.Open(absFilePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	if _, err = io.Copy(hashFunc, file); err != nil {
		return nil, err
	}
	return hashFunc.Sum(nil), nil
}

func hexString(data []byte) string {
	return hex.EncodeToString(data)
}

func base64String(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
