package cache

import (
	"github.com/acexy/golang-toolkit/log"
	"sync"
)

var cacheOnce sync.Once
var caches map[string]*CachingBucket
var rw sync.RWMutex

type CachingBucket struct {
	cache any
}

type Cache interface {
	BucketName() string
	Get(key string) (any, error)
	Put(key string, data any) error
	Evict(key string) error
}

func NewBigCache(bucketName string) *CachingBucket {
	var cache *CachingBucket
	cacheOnce.Do(func() {
		cache = &CachingBucket{}
		caches = make(map[string]*CachingBucket, 1)
		caches[bucketName] = newBigCache()
	})
	log.Logrus().Warnln("bigCache repeated initialization")
	return cache
}

func (c *CachingBucket) GetBucket(bucketName string) *CachingBucket {
	rw.RLock()
	defer rw.RUnlock()
	return caches[bucketName]
}
