package sha256

import (
	"encoding/hex"
	"testing"
)

func TestHexSha256(t *testing.T) {
	t.Log(HexSha256("hello world"))
}

func TestBytesSha256(t *testing.T) {
	t.Log(hex.EncodeToString(BytesSha256([]byte("hello world"))))
}
