package conversion

import "testing"

func TestParseBytesFromHex(t *testing.T) {
	b, err := ParseBytesFromHex("ff00")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 2 || b[0] != 255 || b[1] != 0 {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestParseBytesFromHexError(t *testing.T) {
	if _, err := ParseBytesFromHex("bad-hex"); err == nil {
		t.Fatal("expected hex decode error")
	}
}
