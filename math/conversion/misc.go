package conversion

import "strings"

const (
	hexPrefix  = "0x"
	zeroString = "0"
	maxByte    = 255
)

func appendLeftZero(value string, length int) string {
	if len(value) < length {
		return strings.Repeat(zeroString, length-len(value)) + value
	}
	return value
}
