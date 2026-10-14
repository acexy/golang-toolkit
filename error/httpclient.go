package error

import "errors"

var (
	// ErrUnsupportedHTTPMethod indicates unsupported http method.
	ErrUnsupportedHTTPMethod = errors.New("unsupported http method")
)
