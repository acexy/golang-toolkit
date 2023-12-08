package cache

import (
	"errors"
	"github.com/acexy/golang-toolkit/logger"
	"sync"
)

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

	cachingManager = &CachingManager{
		caches: make(map[string]CachingBucket),
	}
	cachingManager.caches[bucketName] = bucket
	return cachingManager
}

func (c *CachingManager) AddBucket(bucketName string, bucket CachingBucket) {
	if _, flag := c.caches[bucketName]; !flag {
		c.caches[bucketName] = bucket
	} else {
		logger.Logrus().Warnln("duplicate bucketName")
	}
}

func (c *CachingManager) GetBucket(bucketName string) CachingBucket {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.caches[bucketName]
}

func (c *CachingManager) Get(bucketName, key string, result any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Get(key, result)
}

func (c *CachingManager) Put(bucketName, key string, data any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Put(key, data)
}

func (c *CachingManager) Evict(bucketName, key string) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Evict(key)
}
