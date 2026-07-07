package conversion

import (
	"math"
	"testing"
)

func TestFromInt(t *testing.T) {
	tests := []struct {
		value    int
		expected string
	}{
		{10, "10"},
		{-5, "-5"},
		{0, "0"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := FromInt(tt.value)
			if actual != tt.expected {
				t.Errorf("FromInt(%d) = %s; expected %s", tt.value, actual, tt.expected)
			}
		})
	}
}

func TestFromUint64Max(t *testing.T) {
	actual := FromUint64(math.MaxUint64)
	expected := "18446744073709551615"
	if actual != expected {
		t.Fatalf("FromUint64(MaxUint64) = %s; expected %s", actual, expected)
	}
}

func TestFromUint(t *testing.T) {
	tests := []struct {
		value    uint
		expected string
	}{
		{10, "10"},
		{0, "0"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := FromUint(tt.value)
			if actual != tt.expected {
				t.Errorf("FromUint(%d) = %s; expected %s", tt.value, actual, tt.expected)
			}
		})
	}
}
