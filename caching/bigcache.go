package caching

import (
	"context"
	"github.com/acexy/golang-toolkit/util/gob"
	"github.com/allegro/bigcache/v3"
	"time"
)

type BigCacheBucket struct {
	cache *bigcache.BigCache
}

func NewBigCacheByConfig(config bigcache.Config) *BigCacheBucket {
	cache, _ := bigcache.New(context.Background(), config)
	return &BigCacheBucket{cache: cache}
}

func NewSimpleBigCache(duration time.Duration) *BigCacheBucket {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(duration))
	return &BigCacheBucket{cache: cache}
}

func (b *BigCacheBucket) Get(key MemCacheKey, result any, keyAppend ...interface{}) error {
	bs, err := b.cache.Get(originKeyString(key.KeyFormat, keyAppend...))
	if err != nil {
		return err
	}
	return gob.Decode(bs, result)
}

func (b *BigCacheBucket) Put(key MemCacheKey, data any, keyAppend ...interface{}) error {
	bs, err := gob.Encode(data)
	if err != nil {
		return err
	}
	err = b.cache.Set(originKeyString(key.KeyFormat, keyAppend...), bs)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Evict(key MemCacheKey, keyAppend ...interface{}) error {
	return b.cache.Delete(originKeyString(key.KeyFormat, keyAppend...))
}
