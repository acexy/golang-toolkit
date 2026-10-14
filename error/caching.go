package error

import "errors"

var (
	// ErrCacheMiss indicates cache miss.
	ErrCacheMiss = errors.New("cache miss")

	// ErrBadBucketName indicates an unknown or unregistered cache bucket.
	ErrBadBucketName = errors.New("bad bucket name")
)
