package cache

import (
	"errors"
	"github.com/acexy/golang-toolkit/log"
	"sync"
)

var cacheOnce sync.Once

type CachingManager struct {
	rw     sync.RWMutex
	caches map[string]CachingBucket
}

type CachingBucket interface {

	// Get 获取指定key对应的值
	// result 值类型指针
	Get(key string, result any) error

	// Put 设置key对应值
	Put(key string, data any) error

	// Evict 清除缓存
	Evict(key string) error
}

func NewCacheBucketManager(bucketName string, bucket CachingBucket) *CachingManager {
	var cachingManager *CachingManager
	cacheOnce.Do(func() {
		cachingManager = &CachingManager{
			caches: make(map[string]CachingBucket),
		}
		cachingManager.caches[bucketName] = bucket
	})
	log.Logrus().Warnln("bigCache repeated initialization")
	return cachingManager
}

func (c *CachingManager) GetBucket(bucketName string) CachingBucket {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.caches[bucketName]
}

func (c *CachingManager) Get(bucketName, key string, result any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Get(key, result)
}
