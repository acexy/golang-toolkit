package cache

import (
	"context"
	"github.com/acexy/golang-toolkit/util"
	"github.com/allegro/bigcache/v3"
	"time"
)

type BigCacheBucket struct {
	cache *bigcache.BigCache
}

func (b *BigCacheBucket) Get(key string, result any) error {
	v, err := b.cache.Get(key)
	if err != nil {
		return err
	}
	err = util.ParseJsonError(string(v), result)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Put(key string, data any) error {
	bytes, err := util.ToJsonBytesError(data)
	if err != nil {
		return err
	}
	err = b.cache.Set(key, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (b *BigCacheBucket) Evict(key string) error {
	return b.cache.Delete(key)
}

func NewBigCacheByConfig(config bigcache.Config) *BigCacheBucket {
	cache, _ := bigcache.New(context.Background(), config)
	return &BigCacheBucket{cache: cache}
}

func NewSimpleBigCache(duration time.Duration) *BigCacheBucket {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(duration))
	return &BigCacheBucket{cache: cache}
}
