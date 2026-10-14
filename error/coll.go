package error

import "errors"

var (
	// ErrSliceIndexOutOfRange indicates slice index out of range.
	ErrSliceIndexOutOfRange = errors.New("slice index out of range")
)
