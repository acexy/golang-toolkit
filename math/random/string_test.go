package random

import "testing"

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
	if got := UUID(); len(got) != 32 {
		t.Fatalf("unexpected uuid length: %d", len(got))
	}
}
