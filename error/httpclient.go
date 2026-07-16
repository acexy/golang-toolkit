package error

import "errors"

var (
	// ErrUnsupportedHTTPMethod 表示不支持的 HTTP 请求方法
	ErrUnsupportedHTTPMethod = errors.New("unsupported http method")
)
