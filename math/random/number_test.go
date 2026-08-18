package random

import (
	"math"
	"testing"
)

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
	if got := RandInt(math.MaxInt); got < 0 {
		t.Fatalf("RandInt(math.MaxInt) out of range: %d", got)
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
	if got := RandRangeInt(math.MinInt, math.MinInt); got != math.MinInt {
		t.Fatalf("unexpected minimum boundary result: %d", got)
	}
	if got := RandRangeInt(math.MaxInt, math.MaxInt); got != math.MaxInt {
		t.Fatalf("unexpected maximum boundary result: %d", got)
	}
	// 完整整数区间不应因宽度溢出而触发 panic。
	_ = RandRangeInt(math.MinInt, math.MaxInt)
}
