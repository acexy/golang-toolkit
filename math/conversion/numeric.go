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

func ToInt(value string) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		logger.Logrus().Warningln("parse string to numeric error. string:", value)
		return 0
	}
	return v
}

func ToUint(value string) uint {
	v := parseUint(value, 0)
	return uint(v)
}

func ToInt8(value string) int8 {
	v := parseInt(value, 8)
	return int8(v)
}

func ToUint8(value string) uint8 {
	v := parseUint(value, 8)
	return uint8(v)
}

func ToInt16(value string) int16 {
	v := parseInt(value, 16)
	return int16(v)
}

func ToUint16(value string) uint16 {
	v := parseUint(value, 16)
	return uint16(v)
}

func ToInt32(value string) int32 {
	v := parseInt(value, 32)
	return int32(v)
}

func ToUint32(value string) uint32 {
	v := parseUint(value, 32)
	return uint32(v)
}

func ToInt64(value string) int64 {
	return parseInt(value, 64)
}

func ToUint64(value string) uint64 {
	return parseUint(value, 64)
}

func ToFloat32(value string) float32 {
	return float32(parseFloat(value, 32))
}

func ToFloat64(value string) float64 {
	return parseFloat(value, 64)
}