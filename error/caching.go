package error

import "errors"

var (
	// ErrCacheMiss 表示缓存未命中
	ErrCacheMiss = errors.New("cache miss")

	// ErrBadBucketName 表示缓存桶名称无效或未注册
	ErrBadBucketName = errors.New("bad bucket name")
)
