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

// ParseInt parses a string as int, logging and returning zero on failure.
func ParseInt(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseIntPanic parses a string as int and panics on conversion failure.
func ParseIntPanic(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseIntError parses a string as int and returns conversion errors.
func ParseIntError(value string) (int, error) {
	return strconv.Atoi(value)
}

// ParseUint parses a string as uint, logging and returning zero on failure.
func ParseUint(value string) uint {
	v, err := parseUint(value, 0)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint(v)
}

// ParseUintPanic parses a string as uint and panics on conversion failure.
func ParseUintPanic(value string) uint {
	v, err := parseUint(value, 0)
	if err != nil {
		panic(err)
	}
	return uint(v)
}

// ParseUintError parses a string as uint and returns conversion errors.
func ParseUintError(value string) (uint, error) {
	v, err := parseUint(value, 0)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ParseInt8 parses a string as int8, logging and returning zero on failure.
func ParseInt8(value string) int8 {
	v, err := parseInt(value, 8)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int8(v)
}

// ParseInt8Panic parses a string as int8 and panics on conversion failure.
func ParseInt8Panic(value string) int8 {
	v, err := parseInt(value, 8)
	if err != nil {
		panic(err)
	}
	return int8(v)
}

// ParseInt8Error parses a string as int8 and returns conversion errors.
func ParseInt8Error(value string) (int8, error) {
	v, err := parseInt(value, 8)
	if err != nil {
		return 0, err
	}
	return int8(v), nil
}

// ParseUint8 parses a string as uint8, logging and returning zero on failure.
func ParseUint8(value string) uint8 {
	v, err := parseUint(value, 8)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint8(v)
}

// ParseUint8Panic parses a string as uint8 and panics on conversion failure.
func ParseUint8Panic(value string) uint8 {
	v, err := parseUint(value, 8)
	if err != nil {
		panic(err)
	}
	return uint8(v)
}

// ParseUint8Error parses a string as uint8 and returns conversion errors.
func ParseUint8Error(value string) (uint8, error) {
	v, err := parseUint(value, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

// ParseInt16 parses a string as int16, logging and returning zero on failure.
func ParseInt16(value string) int16 {
	v, err := parseInt(value, 16)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int16(v)
}

// ParseInt16Panic parses a string as int16 and panics on conversion failure.
func ParseInt16Panic(value string) int16 {
	v, err := parseInt(value, 16)
	if err != nil {
		panic(err)
	}
	return int16(v)
}

// ParseInt16Error parses a string as int16 and returns conversion errors.
func ParseInt16Error(value string) (int16, error) {
	v, err := parseInt(value, 16)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

// ParseUint16 parses a string as uint16, logging and returning zero on failure.
func ParseUint16(value string) uint16 {
	v, err := parseUint(value, 16)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint16(v)
}

// ParseUint16Panic parses a string as uint16 and panics on conversion failure.
func ParseUint16Panic(value string) uint16 {
	v, err := parseUint(value, 16)
	if err != nil {
		panic(err)
	}
	return uint16(v)
}

// ParseUint16Error parses a string as uint16 and returns conversion errors.
func ParseUint16Error(value string) (uint16, error) {
	v, err := parseUint(value, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

// ParseInt32 parses a string as int32, logging and returning zero on failure.
func ParseInt32(value string) int32 {
	v, err := parseInt(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return int32(v)
}

// ParseInt32Panic parses a string as int32 and panics on conversion failure.
func ParseInt32Panic(value string) int32 {
	v, err := parseInt(value, 32)
	if err != nil {
		panic(err)
	}
	return int32(v)
}

// ParseInt32Error parses a string as int32 and returns conversion errors.
func ParseInt32Error(value string) (int32, error) {
	v, err := parseInt(value, 32)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// ParseUint32 parses a string as uint32, logging and returning zero on failure.
func ParseUint32(value string) uint32 {
	v, err := parseUint(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return uint32(v)
}

// ParseUint32Panic parses a string as uint32 and panics on conversion failure.
func ParseUint32Panic(value string) uint32 {
	v, err := parseUint(value, 32)
	if err != nil {
		panic(err)
	}
	return uint32(v)
}

// ParseUint32Error parses a string as uint32 and returns conversion errors.
func ParseUint32Error(value string) (uint32, error) {
	v, err := parseUint(value, 32)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

// ParseInt64 parses a string as int64, logging and returning zero on failure.
func ParseInt64(value string) int64 {
	v, err := parseInt(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseInt64Panic parses a string as int64 and panics on conversion failure.
func ParseInt64Panic(value string) int64 {
	v, err := parseInt(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseInt64Error parses a string as int64 and returns conversion errors.
func ParseInt64Error(value string) (int64, error) {
	v, err := parseInt(value, 64)
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// ParseUint64 parses a string as uint64, logging and returning zero on failure.
func ParseUint64(value string) uint64 {
	v, err := parseUint(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseUint64Panic parses a string as uint64 and panics on conversion failure.
func ParseUint64Panic(value string) uint64 {
	v, err := parseUint(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseUint64Error parses a string as uint64 and returns conversion errors.
func ParseUint64Error(value string) (uint64, error) {
	v, err := parseUint(value, 64)
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}

// ParseFloat32 parses a string as float32, logging and returning zero on failure.
func ParseFloat32(value string) float32 {
	v, err := parseFloat(value, 32)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return float32(v)
}

// ParseFloat32Panic parses a string as float32 and panics on conversion failure.
func ParseFloat32Panic(value string) float32 {
	v, err := parseFloat(value, 32)
	if err != nil {
		panic(err)
	}
	return float32(v)
}

// ParseFloat32Error parses a string as float32 and returns conversion errors.
func ParseFloat32Error(value string) (float32, error) {
	v, err := parseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

// ParseFloat64 parses a string as float64, logging and returning zero on failure.
func ParseFloat64(value string) float64 {
	v, err := parseFloat(value, 64)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

// ParseFloat64Panic parses a string as float64 and panics on conversion failure.
func ParseFloat64Panic(value string) float64 {
	v, err := parseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// ParseFloat64Error parses a string as float64 and returns conversion errors.
func ParseFloat64Error(value string) (float64, error) {
	v, err := parseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
