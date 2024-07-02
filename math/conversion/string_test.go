package conversion

import "testing"

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
