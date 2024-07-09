package conversion

import (
	"fmt"
	"unsafe"
)

func fromNumber(value int64) string {
	return fmt.Sprintf("%d", value)
}
func fromFloat(value float64) string {
	return fmt.Sprintf("%f", value)
}

func FromInt(value int) string {
	return fromNumber(int64(value))
}

func FromUint(value uint) string {
	return fromNumber(int64(value))
}

func FromInt8(value int8) string {
	return fromNumber(int64(value))
}

func FromUint8(value uint8) string {
	return fromNumber(int64(value))
}

func FromInt16(value int16) string {
	return fromNumber(int64(value))
}

func FromUint16(value uint16) string {
	return fromNumber(int64(value))
}

func FromInt32(value int32) string {
	return fromNumber(int64(value))
}

func FromUint32(value uint32) string {
	return fromNumber(int64(value))
}

func FromInt64(value int64) string {
	return fromNumber(value)
}

func FromUint64(value uint64) string {
	return fromNumber(int64(value))
}

func FromFloat32(value float32) string {
	return fromFloat(float64(value))
}

func FromFloat64(value float64) string {
	return fromFloat(value)
}

// FromBytes 通过go 1.22新函数直接通过字节数组转换成字符串 不进行内存分配高性能
func FromBytes(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
