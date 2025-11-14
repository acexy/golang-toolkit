package caching

import (
	"context"
	"errors"
	"time"

	logrus "github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/util/gob"
	"github.com/allegro/bigcache/v3"
)

type logger struct {
}

func (l logger) Printf(format string, v ...interface{}) {
	logrus.Logrus().Debugln(format, v)
}

type BigCacheBucket struct {
	cache *bigcache.BigCache
}

func NewBigCacheByConfig(config bigcache.Config) *BigCacheBucket {
	config.Logger = logger{}
	cache, _ := bigcache.New(context.Background(), config)
	return &BigCacheBucket{cache: cache}
}

func NewSimpleBigCache(duration time.Duration) *BigCacheBucket {
	c := bigcache.DefaultConfig(duration)
	c.CleanWindow = 5 * time.Second
	c.MaxEntrySize = 1024 * 1024 * 10
	c.StatsEnabled = false
	c.Logger = logger{}
	cache, _ := bigcache.New(context.Background(), c)
	return &BigCacheBucket{cache: cache}
}

func (b *BigCacheBucket) Get(key MemCacheKey, result any, keyAppend ...interface{}) error {
	bs, err := b.GetBytes(key, keyAppend...)
	if err != nil {
		return err
	}
	return gob.Decode(bs, result)
}
func (b *BigCacheBucket) GetBytes(key MemCacheKey, keyAppend ...interface{}) ([]byte, error) {
	bytes, err := b.cache.Get(key.RawKeyString(keyAppend...))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, CacheMiss
		}
		return nil, err
	}
	return bytes, nil
}

func (b *BigCacheBucket) Put(key MemCacheKey, data any, keyAppend ...interface{}) error {
	bs, err := gob.Encode(data)
	if err != nil {
		return err
	}
	err = b.cache.Set(key.RawKeyString(keyAppend...), bs)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Evict(key MemCacheKey, keyAppend ...interface{}) error {
	return b.cache.Delete(key.RawKeyString(keyAppend...))
}
