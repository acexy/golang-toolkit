package conversion

import (
	"testing"
)

func TestToInt(t *testing.T) {
	tests := []struct {
		value string
		want  int
	}{
		{"123", 123},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseInt(tt.value); got != tt.want {
				t.Errorf("ParseInt(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToUint(t *testing.T) {
	tests := []struct {
		value string
		want  uint
	}{
		{"123", 123},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseUint(tt.value); got != tt.want {
				t.Errorf("ParseUint(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToInt8(t *testing.T) {
	tests := []struct {
		value string
		want  int8
	}{
		{"123", 123},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseInt8(tt.value); got != tt.want {
				t.Errorf("ParseInt8(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToUint8(t *testing.T) {
	tests := []struct {
		value string
		want  uint8
	}{
		{"123", 123},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseUint8(tt.value); got != tt.want {
				t.Errorf("ParseUint8(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToInt16(t *testing.T) {
	tests := []struct {
		value string
		want  int16
	}{
		{"12345", 12345},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseInt16(tt.value); got != tt.want {
				t.Errorf("ParseInt16(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToUint16(t *testing.T) {
	tests := []struct {
		value string
		want  uint16
	}{
		{"12345", 12345},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseUint16(tt.value); got != tt.want {
				t.Errorf("ParseUint16(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		value string
		want  int32
	}{
		{"123456789", 123456789},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseInt32(tt.value); got != tt.want {
				t.Errorf("ParseInt32(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToUint32(t *testing.T) {
	tests := []struct {
		value string
		want  uint32
	}{
		{"123456789", 123456789},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseUint32(tt.value); got != tt.want {
				t.Errorf("ParseUint32(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		value string
		want  int64
	}{
		{"1234567890123456789", 1234567890123456789},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseInt64(tt.value); got != tt.want {
				t.Errorf("ParseInt64(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToUint64(t *testing.T) {
	tests := []struct {
		value string
		want  uint64
	}{
		{"1234567890123456789", 1234567890123456789},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseUint64(tt.value); got != tt.want {
				t.Errorf("ParseUint64(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		value string
		want  float32
	}{
		{"123.456", 123.456},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseFloat32(tt.value); got != tt.want {
				t.Errorf("ParseFloat32(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		value string
		want  float64
	}{
		{"123.456", 123.456},
		{"abc", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ParseFloat64(tt.value); got != tt.want {
				t.Errorf("ParseFloat64(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
