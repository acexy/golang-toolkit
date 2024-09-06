package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/acexy/golang-toolkit/math/conversion"
)

func HexSha256(data string) string {
	return hex.EncodeToString(BytesSha256(conversion.ParseBytes(data)))
}

func BytesSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
