package conversion

import (
	"strconv"
	"unsafe"
)

func fromNumber(value int64) string {
	return strconv.FormatInt(value, 10)
}
func fromUnsignedNumber(value uint64) string {
	return strconv.FormatUint(value, 10)
}
func fromFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// FromInt formats a int as a decimal string.
func FromInt(value int) string {
	return fromNumber(int64(value))
}

// FromUint formats a uint as a decimal string.
func FromUint(value uint) string {
	return fromUnsignedNumber(uint64(value))
}

// FromInt8 formats a int8 as a decimal string.
func FromInt8(value int8) string {
	return fromNumber(int64(value))
}

// FromUint8 formats a uint8 as a decimal string.
func FromUint8(value uint8) string {
	return fromUnsignedNumber(uint64(value))
}

// FromInt16 formats a int16 as a decimal string.
func FromInt16(value int16) string {
	return fromNumber(int64(value))
}

// FromUint16 formats a uint16 as a decimal string.
func FromUint16(value uint16) string {
	return fromUnsignedNumber(uint64(value))
}

// FromInt32 formats a int32 as a decimal string.
func FromInt32(value int32) string {
	return fromNumber(int64(value))
}

// FromUint32 formats a uint32 as a decimal string.
func FromUint32(value uint32) string {
	return fromUnsignedNumber(uint64(value))
}

// FromInt64 formats a int64 as a decimal string.
func FromInt64(value int64) string {
	return fromNumber(value)
}

// FromUint64 formats a uint64 as a decimal string.
func FromUint64(value uint64) string {
	return fromUnsignedNumber(value)
}

// FromFloat32 formats a float32 as a decimal string using float64 precision.
func FromFloat32(value float32) string {
	return fromFloat(float64(value))
}

// FromFloat64 formats a float64 as a decimal string.
func FromFloat64(value float64) string {
	return fromFloat(value)
}

// FromBytes returns a string view without copying; do not modify the backing bytes while it is in use.
// Use string(bytes) when an independent string is required.
func FromBytes(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
