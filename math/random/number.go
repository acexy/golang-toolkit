package random

import (
	"math/rand/v2"
)

// RandInt returns a pseudo-random integer in [0, max], or -1 for a negative max.
func RandInt(max int) int {
	if max < 0 {
		return -1
	}
	return int(rand.Uint64N(uint64(max) + 1))
}

// RandRangeInt returns a pseudo-random integer in [min, max], or -1 for an invalid range.
func RandRangeInt(min, max int) int {
	if min > max {
		return -1
	}

	// Compute width using unsigned arithmetic to avoid signed boundary overflow.
	width := uint64(max) - uint64(min) + 1
	if width == 0 {
		// A full 64-bit int range has width 2^64, represented here by zero.
		return int(rand.Uint64())
	}
	return int(uint64(min) + rand.Uint64N(width))
}
