package caching

import (
	"context"
	"errors"
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
	bs, err := b.GetBytes(key, keyAppend...)
	if err != nil {
		return err
	}
	return gob.Decode(bs, result)
}
func (b *BigCacheBucket) GetBytes(key MemCacheKey, keyAppend ...interface{}) ([]byte, error) {
	bytes, err := b.cache.Get(OriginKeyString(key.KeyFormat, keyAppend...))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, SourceNotFound
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
	err = b.cache.Set(OriginKeyString(key.KeyFormat, keyAppend...), bs)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Evict(key MemCacheKey, keyAppend ...interface{}) error {
	return b.cache.Delete(OriginKeyString(key.KeyFormat, keyAppend...))
}
