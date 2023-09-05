package cache

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

type BigCache struct {
	CachingBucket
}

func (b *BigCache) BucketName() string {
	return ""
}
func (b *BigCache) Get(key string) (any, error) {
	return nil, nil
}
func (b *BigCache) Put(key string, data any) error {
	return nil
}
func (b *BigCache) Evict(key string) error {
	return nil
}

func newBigCache() *CachingBucket {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	return &CachingBucket{
		cache: cache,
	}
}
