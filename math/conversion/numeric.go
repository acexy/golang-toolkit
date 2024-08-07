package conversion

import (
	"github.com/acexy/golang-toolkit/logger"
	"strconv"
)

func parseInt(value string, bit int) int64 {
	v, err := strconv.ParseInt(value, 10, bit)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

func parseUint(value string, bit int) uint64 {
	v, err := strconv.ParseUint(value, 10, bit)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

func parseFloat(value string, bit int) float64 {
	v, err := strconv.ParseFloat(value, bit)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseInt 将字符串转换为int
func ParseInt(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseUint 将字符串转换为uint
func ParseUint(value string) uint {
	v := parseUint(value, 0)
	return uint(v)
}

// ParseInt8 将字符串转换为int8
func ParseInt8(value string) int8 {
	v := parseInt(value, 8)
	return int8(v)
}

// ParseUint8 将字符串转换为uint8
func ParseUint8(value string) uint8 {
	v := parseUint(value, 8)
	return uint8(v)
}

// ParseInt16 将字符串转换为int16
func ParseInt16(value string) int16 {
	v := parseInt(value, 16)
	return int16(v)
}

// ParseUint16 将字符串转换为uint16
func ParseUint16(value string) uint16 {
	v := parseUint(value, 16)
	return uint16(v)
}

// ParseInt32 将字符串转换为int32
func ParseInt32(value string) int32 {
	v := parseInt(value, 32)
	return int32(v)
}

// ParseUint32 将字符串转换为uint32
func ParseUint32(value string) uint32 {
	v := parseUint(value, 32)
	return uint32(v)
}

// ParseInt64 将字符串转换为int64
func ParseInt64(value string) int64 {
	return parseInt(value, 64)
}

// ParseUint64 将字符串转换为uint64
func ParseUint64(value string) uint64 {
	return parseUint(value, 64)
}

// ParseFloat32 将字符串转换为float32
func ParseFloat32(value string) float32 {
	return float32(parseFloat(value, 32))
}

// ParseFloat64 将字符串转换为float64
func ParseFloat64(value string) float64 {
	return parseFloat(value, 64)
}
