package conversion

import (
	"strconv"
	"unsafe"
)

func fromNumber(value int64) string {
	return strconv.FormatInt(value, 10)
}
func fromFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// FromInt Int转字符串
func FromInt(value int) string {
	return fromNumber(int64(value))
}

// FromUint Uint转字符串
func FromUint(value uint) string {
	return fromNumber(int64(value))
}

// FromInt8 Int8转字符串
func FromInt8(value int8) string {
	return fromNumber(int64(value))
}

// FromUint8 Uint8转字符串
func FromUint8(value uint8) string {
	return fromNumber(int64(value))
}

// FromInt16 Int16转字符串
func FromInt16(value int16) string {
	return fromNumber(int64(value))
}

// FromUint16 Uint16转字符串
func FromUint16(value uint16) string {
	return fromNumber(int64(value))
}

// FromInt32 Int32转字符串
func FromInt32(value int32) string {
	return fromNumber(int64(value))
}

// FromUint32 Uint32转字符串
func FromUint32(value uint32) string {
	return fromNumber(int64(value))
}

// FromInt64 Int64转字符串
func FromInt64(value int64) string {
	return fromNumber(value)
}

// FromUint64 Uint64转字符串
func FromUint64(value uint64) string {
	return fromNumber(int64(value))
}

// FromFloat32 Float32转字符串
func FromFloat32(value float32) string {
	return fromFloat(float64(value))
}

// FromFloat64 Float64转字符串
func FromFloat64(value float64) string {
	return fromFloat(value)
}

// FromBytes 通过go 1.22新函数直接通过字节数组转换成字符串 不进行内存分配高性能 存在不安全性，共享底层数据
func FromBytes(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}
