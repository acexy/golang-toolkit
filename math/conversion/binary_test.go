package conversion

import (
	"fmt"
	"testing"
)

func TestNewFromRawBinary(t *testing.T) {
	b, _ := NewFromRawBinary("11111111")
	fmt.Println(b.ToBit())
	fmt.Println(b.ToHex())
}

func TestNewFromRawHex(t *testing.T) {
	b, _ := NewFromRawHex("ff1")
	fmt.Println(b.Value())
}

func TestNewFormBytes(t *testing.T) {
	fmt.Println(NewFormBytes([]byte{255, 0, 3, 45, 251}).ToBit(","))
}
