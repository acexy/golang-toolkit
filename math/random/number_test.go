package random

import "testing"

func TestRandInt(t *testing.T) {
	if got := RandInt(-1); got != -1 {
		t.Fatalf("RandInt(-1) = %d; expected -1", got)
	}
	if got := RandInt(0); got != 0 {
		t.Fatalf("RandInt(0) = %d; expected 0", got)
	}
	for range 20 {
		got := RandInt(2)
		if got < 0 || got > 2 {
			t.Fatalf("RandInt(2) out of range: %d", got)
		}
	}
}

func TestRandRangeInt(t *testing.T) {
	if got := RandRangeInt(2, 1); got != -1 {
		t.Fatalf("RandRangeInt(2, 1) = %d; expected -1", got)
	}
	for range 20 {
		got := RandRangeInt(1, 2)
		if got < 1 || got > 2 {
			t.Fatalf("RandRangeInt(1, 2) out of range: %d", got)
		}
	}
}
