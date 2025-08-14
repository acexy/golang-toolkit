package caching

import (
	"errors"
	"fmt"

	"github.com/acexy/golang-toolkit/logger"
)

var (
	CacheMiss = errors.New("cache miss")
)

type MemCacheKey struct {
	// 最终key值的格式化格式 将使用 fmt.Sprintf(key.KeyFormat, keyAppend) 进行处理
	KeyFormat string
}

// NewNemCacheKey 创建一个缓存key
func NewNemCacheKey(keyFormat string) MemCacheKey {
	return MemCacheKey{KeyFormat: keyFormat}
}

// RawKeyString 获取原始的key字符串
func (m MemCacheKey) RawKeyString(keyAppend ...interface{}) string {
	if len(keyAppend) > 0 {
		return fmt.Sprintf(m.KeyFormat, keyAppend...)
	}
	return m.KeyFormat
}

type CacheManager struct {
	caches map[string]CacheBucket
}

type CacheBucket interface {

	// Get 获取指定key对应的值
	// result 值类型指针 如果未能查到内容应当返还
	Get(key MemCacheKey, result any, keyAppend ...interface{}) error

	// GetBytes 获取指定key对应的值
	GetBytes(key MemCacheKey, keyAppend ...interface{}) ([]byte, error)

	// Put 设置key对应值
	Put(key MemCacheKey, data any, keyAppend ...interface{}) error

	// Evict 清除缓存
	Evict(key MemCacheKey, keyAppend ...interface{}) error
}

func NewCacheBucketManager(bucketName string, bucket CacheBucket) *CacheManager {
	manager := NewEmptyCacheBucketManager()
	manager.caches[bucketName] = bucket
	return manager
}

func NewEmptyCacheBucketManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]CacheBucket),
	}
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

func (c *CacheManager) Get(bucketName string, key MemCacheKey, result any, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Get(key, result, keyAppend...)
}

func (c *CacheManager) Put(bucketName string, key MemCacheKey, data any, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Put(key, data, keyAppend...)
}

func (c *CacheManager) Evict(bucketName string, key MemCacheKey, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Warnln("bad bucketName", bucketName)
		return errors.New("bad bucketName " + bucketName)
	}
	return bucket.Evict(key, keyAppend...)
}
