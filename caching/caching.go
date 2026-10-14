package caching

import (
	"fmt"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/util/gob"
)

// CacheKey defines a cache key format with optional fmt.Sprintf arguments.
type CacheKey struct {
	// KeyFormat is the format used to construct the final cache key.
	KeyFormat string
}

// BucketName identifies a registered cache bucket.
type BucketName string

// Codec encodes cache values and decodes them into caller-provided destinations.
type Codec interface {
	// Encode serializes a value to bytes.
	Encode(data any) ([]byte, error)

	// Decode decodes bytes into a writable destination.
	Decode(bs []byte, result any) error
}

// GobCodec implements Codec with Go GOB serialization.
type GobCodec struct {
}

// Encode serializes a value using GOB.
func (g GobCodec) Encode(data any) ([]byte, error) {
	return gob.Encode(data)
}

// Decode decodes GOB bytes into the destination.
func (g GobCodec) Decode(bs []byte, result any) error {
	return gob.Decode(bs, result)
}

func defaultCodec(codec Codec) Codec {
	if codec != nil {
		return codec
	}
	return GobCodec{}
}

// NewCacheKey creates a formatted cache key.
func NewCacheKey(keyFormat string) CacheKey {
	return CacheKey{KeyFormat: keyFormat}
}

// NewBucketName creates a bucket identifier.
func NewBucketName(bucketName string) BucketName {
	return BucketName(bucketName)
}

// RawKeyString formats the key when arguments are supplied, otherwise returns the raw format.
func (m CacheKey) RawKeyString(keyAppend ...any) string {
	if len(keyAppend) > 0 {
		return fmt.Sprintf(m.KeyFormat, keyAppend...)
	}
	return m.KeyFormat
}

// CacheManager manages named buckets and a shared value codec.
type CacheManager struct {
	lock   sync.RWMutex
	caches map[BucketName]CacheBucket
	codec  Codec
}

// CacheBucket defines byte storage operations for a cache bucket.
type CacheBucket interface {

	// GetBytes reads bytes under a formatted key.
	GetBytes(key CacheKey, keyAppend ...any) ([]byte, error)

	// PutBytes stores bytes under a formatted key.
	PutBytes(key CacheKey, data []byte, keyAppend ...any) error

	// Evict removes a formatted key.
	Evict(key CacheKey, keyAppend ...any) error
}

// NewCacheManager creates a manager using the first codec or GOB by default.
func NewCacheManager(codec ...Codec) *CacheManager {
	var selectedCodec Codec
	if len(codec) > 0 {
		selectedCodec = codec[0]
	}
	return &CacheManager{
		caches: make(map[BucketName]CacheBucket),
		codec:  defaultCodec(selectedCodec),
	}
}

// AddBucket registers a bucket; duplicate names are logged and keep the existing bucket.
func (c *CacheManager) AddBucket(bucketName BucketName, bucket CacheBucket) *CacheManager {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, flag := c.caches[bucketName]; !flag {
		c.caches[bucketName] = bucket
	} else {
		logger.Logrus().Errorln("caching: duplicate bucketName", bucketName)
	}
	return c
}

// GetBucket returns the named bucket, or nil if it is not registered.
func (c *CacheManager) GetBucket(bucketName BucketName) CacheBucket {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.caches[bucketName]
}

// Get reads and decodes a value into result; unknown buckets return ErrBadBucketName.
func (c *CacheManager) Get(bucketName BucketName, key CacheKey, result any, keyAppend ...any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	bs, err := bucket.GetBytes(key, keyAppend...)
	if err != nil {
		return err
	}
	return c.codec.Decode(bs, result)
}

// GetBytes reads raw bytes from a registered bucket.
func (c *CacheManager) GetBytes(bucketName BucketName, key CacheKey, keyAppend ...any) ([]byte, error) {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return nil, toolkitError.ErrBadBucketName
	}
	return bucket.GetBytes(key, keyAppend...)
}

// Put encodes and stores a value in a registered bucket.
func (c *CacheManager) Put(bucketName BucketName, key CacheKey, data any, keyAppend ...any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	bs, err := c.codec.Encode(data)
	if err != nil {
		return err
	}
	return bucket.PutBytes(key, bs, keyAppend...)
}

// PutBytes stores raw bytes in a registered bucket.
func (c *CacheManager) PutBytes(bucketName BucketName, key CacheKey, data []byte, keyAppend ...any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	return bucket.PutBytes(key, data, keyAppend...)
}

// Evict removes a key from a registered bucket.
func (c *CacheManager) Evict(bucketName BucketName, key CacheKey, keyAppend ...any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	return bucket.Evict(key, keyAppend...)
}
