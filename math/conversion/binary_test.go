package conversion

import (
	"fmt"
	"testing"
)

func TestNewFromRawBinary(t *testing.T) {
	b, _ := NewFromRawBinary("11111111")
	fmt.Println(b.To8Bit())
	fmt.Println(b.To2Hex())
}

func TestNewFromRawHex(t *testing.T) {
	b, _ := NewFromRawHex("ff1")
	fmt.Println(b.Value())
}

func TestNewFormBytes(t *testing.T) {
	bs := NewFormBytes([]byte{255, 0, 3, 45, 251})
	fmt.Println(bs.To8Bits(","))
	fmt.Println(bs.To2HexString(","))
}
