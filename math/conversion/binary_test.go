package conversion

import (
	"errors"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func TestNewFromRawBinary(t *testing.T) {
	b, err := NewFromRawBinary("11111111")
	if err != nil {
		t.Fatal(err)
	}
	bit, err := b.To8Bit()
	if err != nil {
		t.Fatal(err)
	}
	if bit != "11111111" {
		t.Fatalf("unexpected bit value: %s", bit)
	}
	hex, err := b.To2Hex()
	if err != nil {
		t.Fatal(err)
	}
	if hex != "ff" {
		t.Fatalf("unexpected hex value: %s", hex)
	}
}

func TestBinaryOutOfByteRange(t *testing.T) {
	b := NewFromDecimal(256)
	if _, err := b.To8Bit(); !errors.Is(err, toolkitError.ErrBinaryValueOutOfByteRange) {
		t.Fatalf("expected ErrBinaryValueOutOfByteRange, got %v", err)
	}
	if _, err := b.To2Hex(); !errors.Is(err, toolkitError.ErrBinaryValueOutOfByteRange) {
		t.Fatalf("expected ErrBinaryValueOutOfByteRange, got %v", err)
	}
}

func TestNewFromRawHex(t *testing.T) {
	b, err := NewFromRawHex("0xff")
	if err != nil {
		t.Fatal(err)
	}
	if b.Value() != "11111111" {
		t.Fatalf("unexpected binary value: %s", b.Value())
	}
}

func TestNewFromBytes(t *testing.T) {
	bs := NewFromBytes([]byte{255, 0, 3})
	bits, err := bs.To8Bits(",")
	if err != nil {
		t.Fatal(err)
	}
	if bits != "11111111,00000000,00000011" {
		t.Fatalf("unexpected bits: %s", bits)
	}
	if hex := bs.To2Hex(","); hex != "ff,00,03" {
		t.Fatalf("unexpected hex: %s", hex)
	}
}
