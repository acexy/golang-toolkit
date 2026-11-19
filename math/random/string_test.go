package random

import (
	"encoding/hex"
	"testing"
)

func TestRandString(t *testing.T) {
	if got := RandString(5); len(got) != 5 {
		t.Fatalf("unexpected string length: %d", len(got))
	}
	if got := RandString(0); got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
	if got := RandString(-1); got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
}

func TestUUID(t *testing.T) {
	got := UUID()
	decoded, err := hex.DecodeString(got)
	if err != nil || len(decoded) != 16 || hex.EncodeToString(decoded) != got {
		t.Fatalf("UUID() = %q, want 32 lowercase hexadecimal characters", got)
	}
	if decoded[6]>>4 != 4 || decoded[8]>>6 != 2 {
		t.Fatalf("UUID() = %q, want an RFC 9562 version 4 UUID", got)
	}
}
