package conversion

import (
	"strconv"

	"github.com/acexy/golang-toolkit/logger"
)

func parseInt(value string, bit int) (int64, error) {
	return strconv.ParseInt(value, 10, bit)
}

func parseUint(value string, bit int) (uint64, error) {
	return strconv.ParseUint(value, 10, bit)
}

func parseFloat(value string, bit int) (float64, error) {
	return strconv.ParseFloat(value, bit)
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

// ParseIntPanic 将字符串转换为int 异常将触发panic
func ParseIntPanic(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseIntError 将字符串转换为int 可能返回错误
func ParseIntError(value string) (int, error) {
	return strconv.Atoi(value)
}

// ParseUint 将字符串转换为uint
func ParseUint(value string) uint {
	v, err := parseUint(value, 0)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint(v)
}

// ParseUintPanic 将字符串转换为uint 异常将触发panic
func ParseUintPanic(value string) uint {
	v, err := parseUint(value, 0)
	if err != nil {
		panic(err)
	}
	return uint(v)
}

// ParseUintError 将字符串转换为uint 可能返回错误
func ParseUintError(value string) (uint, error) {
	v, err := parseUint(value, 0)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ParseInt8 将字符串转换为int8
func ParseInt8(value string) int8 {
	v, err := parseInt(value, 8)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int8(v)
}

// ParseInt8Panic 将字符串转换为int8 异常将触发panic
func ParseInt8Panic(value string) int8 {
	v, err := parseInt(value, 8)
	if err != nil {
		panic(err)
	}
	return int8(v)
}

// ParseInt8Error 将字符串转换为int8 可能返回错误
func ParseInt8Error(value string) (int8, error) {
	v, err := parseInt(value, 8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

// ParseUint8 将字符串转换为uint8
func ParseUint8(value string) uint8 {
	v, err := parseUint(value, 8)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint8(v)
}

// ParseUint8Panic 将字符串转换为uint8 异常将触发panic
func ParseUint8Panic(value string) uint8 {
	v, err := parseUint(value, 8)
	if err != nil {
		panic(err)
	}
	return uint8(v)
}

// ParseUint8Error 将字符串转换为uint8 可能返回错误
func ParseUint8Error(value string) (uint8, error) {
	v, err := parseUint(value, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

// ParseInt16 将字符串转换为int16
func ParseInt16(value string) int16 {
	v, err := parseInt(value, 16)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int16(v)
}

// ParseInt16Panic 将字符串转换为int16 异常将触发panic
func ParseInt16Panic(value string) int16 {
	v, err := parseInt(value, 16)
	if err != nil {
		panic(err)
	}
	return int16(v)
}

// ParseInt16Error 将字符串转换为int16
func ParseInt16Error(value string) (int16, error) {
	v, err := parseInt(value, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

// ParseUint16 将字符串转换为uint16
func ParseUint16(value string) uint16 {
	v, err := parseUint(value, 16)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint16(v)
}

// ParseUint16Panic 将字符串转换为uint16 异常将触发panic
func ParseUint16Panic(value string) uint16 {
	v, err := parseUint(value, 16)
	if err != nil {
		panic(err)
	}
	return uint16(v)
}

// ParseUint16Error 将字符串转换为uint16
func ParseUint16Error(value string) (uint16, error) {
	v, err := parseUint(value, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

// ParseInt32 将字符串转换为int32
func ParseInt32(value string) int32 {
	v, err := parseInt(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int32(v)
}

// ParseInt32Panic 将字符串转换为int32 异常将触发panic
func ParseInt32Panic(value string) int32 {
	v, err := parseInt(value, 32)
	if err != nil {
		panic(err)
	}
	return int32(v)
}

// ParseInt32Error 将字符串转换为int32
func ParseInt32Error(value string) (int32, error) {
	v, err := parseInt(value, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// ParseUint32 将字符串转换为uint32
func ParseUint32(value string) uint32 {
	v, err := parseUint(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint32(v)
}

// ParseUint32Panic 将字符串转换为uint32 异常将触发panic
func ParseUint32Panic(value string) uint32 {
	v, err := parseUint(value, 32)
	if err != nil {
		panic(err)
	}
	return uint32(v)
}

// ParseUint32Error 将字符串转换为uint32
func ParseUint32Error(value string) (uint32, error) {
	v, err := parseUint(value, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

// ParseInt64 将字符串转换为int64
func ParseInt64(value string) int64 {
	v, err := parseInt(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseInt64Panic 将字符串转换为int64 异常将触发panic
func ParseInt64Panic(value string) int64 {
	v, err := parseInt(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseInt64Error 将字符串转换为int64
func ParseInt64Error(value string) (int64, error) {
	v, err := parseInt(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// ParseUint64 将字符串转换为uint64
func ParseUint64(value string) uint64 {
	v, err := parseUint(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseUint64Panic 将字符串转换为uint64 异常将触发panic
func ParseUint64Panic(value string) uint64 {
	v, err := parseUint(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseUint64Error 将字符串转换为uint64
func ParseUint64Error(value string) (uint64, error) {
	v, err := parseUint(value, 64)
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}

// ParseFloat32 将字符串转换为float32
func ParseFloat32(value string) float32 {
	v, err := parseFloat(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return float32(v)
}

// ParseFloat32Panic 将字符串转换为float32 异常将触发panic
func ParseFloat32Panic(value string) float32 {
	v, err := parseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return float32(v)
}

// ParseFloat32Error 将字符串转换为float32
func ParseFloat32Error(value string) (float32, error) {
	v, err := parseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

// ParseFloat64 将字符串转换为float64
func ParseFloat64(value string) float64 {
	v, err := parseFloat(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseFloat64Panic 将字符串转换为float64 异常将触发panic
func ParseFloat64Panic(value string) float64 {
	v, err := parseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseFloat64Error 将字符串转换为float64
func ParseFloat64Error(value string) (float64, error) {
	v, err := parseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
