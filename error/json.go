package error

import "errors"

var (
	// ErrUnsupportedTimestampType indicates unsupported timestamp type.
	ErrUnsupportedTimestampType = errors.New("unsupported timestamp type")
)
