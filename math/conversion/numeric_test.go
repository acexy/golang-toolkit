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
			if got := ToInt(tt.value); got != tt.want {
				t.Errorf("ToInt(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToUint(tt.value); got != tt.want {
				t.Errorf("ToUint(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToInt8(tt.value); got != tt.want {
				t.Errorf("ToInt8(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToUint8(tt.value); got != tt.want {
				t.Errorf("ToUint8(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToInt16(tt.value); got != tt.want {
				t.Errorf("ToInt16(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToUint16(tt.value); got != tt.want {
				t.Errorf("ToUint16(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToInt32(tt.value); got != tt.want {
				t.Errorf("ToInt32(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToUint32(tt.value); got != tt.want {
				t.Errorf("ToUint32(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToInt64(tt.value); got != tt.want {
				t.Errorf("ToInt64(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToUint64(tt.value); got != tt.want {
				t.Errorf("ToUint64(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToFloat32(tt.value); got != tt.want {
				t.Errorf("ToFloat32(%v) = %v, want %v", tt.value, got, tt.want)
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
			if got := ToFloat64(tt.value); got != tt.want {
				t.Errorf("ToFloat64(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
