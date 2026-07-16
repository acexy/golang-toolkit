package error

import "errors"

var (
	// ErrUnsupportedTimestampType 表示不支持的时间戳类型
	ErrUnsupportedTimestampType = errors.New("unsupported timestamp type")
)
