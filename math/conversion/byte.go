package conversion

import (
	"encoding/hex"
	"unsafe"
)

// ParseBytes returns a read-only byte view without copying; do not modify it.
// Use []byte(s) when an independent writable slice is required.
func ParseBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// ParseBytesFromHex decodes a hexadecimal string into bytes.
func ParseBytesFromHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
