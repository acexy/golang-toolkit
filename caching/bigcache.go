package caching

import (
	"context"
	"errors"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
	logrus "github.com/acexy/golang-toolkit/logger"
	"github.com/allegro/bigcache/v3"
)

type log struct {
}

// Printf routes BigCache diagnostics to the toolkit debug logger.
func (l log) Printf(format string, v ...any) {
	logrus.Logrus().Debugln(format, v)
}

// BigCacheBucket adapts BigCache to byte-oriented cache operations.
type BigCacheBucket struct {
	cache *bigcache.BigCache
}

// NewBigCacheByConfig creates a BigCache bucket with the supplied configuration and toolkit logging.
func NewBigCacheByConfig(config bigcache.Config) (*BigCacheBucket, error) {
	config.Logger = log{}
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &BigCacheBucket{cache: cache}, nil
}

// NewSimpleBigCache creates a bucket with the given lifetime and a five-second cleanup interval.
func NewSimpleBigCache(duration time.Duration) (*BigCacheBucket, error) {
	c := bigcache.DefaultConfig(duration)
	c.CleanWindow = 5 * time.Second
	c.StatsEnabled = false
	c.Logger = log{}
	cache, err := bigcache.New(context.Background(), c)
	if err != nil {
		return nil, err
	}
	return &BigCacheBucket{cache: cache}, nil
}

// GetBytes reads cached bytes and returns ErrCacheMiss for a missing entry.
func (b *BigCacheBucket) GetBytes(key CacheKey, keyAppend ...any) ([]byte, error) {
	bytes, err := b.cache.Get(key.RawKeyString(keyAppend...))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, toolkitError.ErrCacheMiss
		}
		return nil, err
	}
	return bytes, nil
}

// PutBytes stores bytes under the formatted key.
func (b *BigCacheBucket) PutBytes(key CacheKey, data []byte, keyAppend ...any) error {
	err := b.cache.Set(key.RawKeyString(keyAppend...), data)
	if err != nil {
		return err
	}
	return nil
}

// Evict removes the formatted key from the bucket.
func (b *BigCacheBucket) Evict(key CacheKey, keyAppend ...any) error {
	return b.cache.Delete(key.RawKeyString(keyAppend...))
}
