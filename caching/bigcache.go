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

func (l log) Printf(format string, v ...interface{}) {
	logrus.Logrus().Debugln(format, v)
}

type BigCacheBucket struct {
	cache *bigcache.BigCache
}

func NewBigCacheByConfig(config bigcache.Config) (*BigCacheBucket, error) {
	config.Logger = log{}
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &BigCacheBucket{cache: cache}, nil
}

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

func (b *BigCacheBucket) GetBytes(key CacheKey, keyAppend ...interface{}) ([]byte, error) {
	bytes, err := b.cache.Get(key.RawKeyString(keyAppend...))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, toolkitError.ErrCacheMiss
		}
		return nil, err
	}
	return bytes, nil
}

func (b *BigCacheBucket) PutBytes(key CacheKey, data []byte, keyAppend ...interface{}) error {
	err := b.cache.Set(key.RawKeyString(keyAppend...), data)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Evict(key CacheKey, keyAppend ...interface{}) error {
	return b.cache.Delete(key.RawKeyString(keyAppend...))
}
