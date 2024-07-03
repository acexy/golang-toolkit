package caching

import (
	"errors"
	"github.com/acexy/golang-toolkit/logger"
	"sync"
)

type CacheManager struct {
	caches map[string]CacheBucket
}

type CacheBucket interface {

	// Get 获取指定key对应的值
	// result 值类型指针
	Get(key string, result any) error

	// Put 设置key对应值
	Put(key string, data any) error

	// Evict 清除缓存
	Evict(key string) error
}

var once sync.Once
var cachingManager *CacheManager

func NewCacheBucketManager(bucketName string, bucket CacheBucket) *CacheManager {
	if cachingManager == nil {
		NewEmptyCacheBucketManager()
	}
	cachingManager.caches[bucketName] = bucket
	return cachingManager
}

func NewEmptyCacheBucketManager() *CacheManager {
	once.Do(func() {
		cachingManager = &CacheManager{
			caches: make(map[string]CacheBucket),
		}
	})
	return cachingManager
}

func (c *CacheManager) AddBucket(bucketName string, bucket CacheBucket) {
	if _, flag := c.caches[bucketName]; !flag {
		c.caches[bucketName] = bucket
	} else {
		logger.Logrus().Warnln("duplicate bucketName")
	}
}

func (c *CacheManager) GetBucket(bucketName string) CacheBucket {
	return c.caches[bucketName]
}

func (c *CacheManager) Get(bucketName, key string, result any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Get(key, result)
}

func (c *CacheManager) Put(bucketName, key string, data any) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Put(key, data)
}

func (c *CacheManager) Evict(bucketName, key string) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Evict(key)
}
