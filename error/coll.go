package error

import "errors"

var (
	// ErrSliceIndexOutOfRange 表示切片索引越界
	ErrSliceIndexOutOfRange = errors.New("slice index out of range")
)
