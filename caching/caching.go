package caching

import (
	"fmt"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/util/gob"
)

type CacheKey struct {
	// 最终key值的格式化格式 将使用 fmt.Sprintf(key.KeyFormat, keyAppend) 进行处理
	KeyFormat string
}

type BucketName string

type Codec interface {
	// Encode 将缓存值编码为字节数组
	Encode(data any) ([]byte, error)

	// Decode 将缓存字节数组解码到 result 指针
	Decode(bs []byte, result any) error
}

type GobCodec struct {
}

func (g GobCodec) Encode(data any) ([]byte, error) {
	return gob.Encode(data)
}

func (g GobCodec) Decode(bs []byte, result any) error {
	return gob.Decode(bs, result)
}

func defaultCodec(codec Codec) Codec {
	if codec != nil {
		return codec
	}
	return GobCodec{}
}

// NewCacheKey 创建一个缓存key
func NewCacheKey(keyFormat string) CacheKey {
	return CacheKey{KeyFormat: keyFormat}
}

// NewBucketName 创建一个缓存桶名称
func NewBucketName(bucketName string) BucketName {
	return BucketName(bucketName)
}

// RawKeyString 获取原始的key字符串
func (m CacheKey) RawKeyString(keyAppend ...interface{}) string {
	if len(keyAppend) > 0 {
		return fmt.Sprintf(m.KeyFormat, keyAppend...)
	}
	return m.KeyFormat
}

type CacheManager struct {
	lock   sync.RWMutex
	caches map[BucketName]CacheBucket
	codec  Codec
}

type CacheBucket interface {

	// GetBytes 获取指定key对应的值
	GetBytes(key CacheKey, keyAppend ...interface{}) ([]byte, error)

	// PutBytes 设置key对应的字节数组
	PutBytes(key CacheKey, data []byte, keyAppend ...interface{}) error

	// Evict 清除缓存
	Evict(key CacheKey, keyAppend ...interface{}) error
}

func NewCacheManager(codec ...Codec) *CacheManager {
	var selectedCodec Codec
	if len(codec) > 0 {
		selectedCodec = codec[0]
	}
	return &CacheManager{
		caches: make(map[BucketName]CacheBucket),
		codec:  defaultCodec(selectedCodec),
	}
}

// AddBucket 添加一个缓存桶
func (c *CacheManager) AddBucket(bucketName BucketName, bucket CacheBucket) *CacheManager {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, flag := c.caches[bucketName]; !flag {
		c.caches[bucketName] = bucket
	} else {
		logger.Logrus().Errorln("caching: duplicate bucketName", bucketName)
	}
	return c
}

// GetBucket 获取缓存桶
func (c *CacheManager) GetBucket(bucketName BucketName) CacheBucket {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.caches[bucketName]
}

// Get 获取缓存
func (c *CacheManager) Get(bucketName BucketName, key CacheKey, result any, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	bs, err := bucket.GetBytes(key, keyAppend...)
	if err != nil {
		return err
	}
	return c.codec.Decode(bs, result)
}

// GetBytes 获取缓存字节数组
func (c *CacheManager) GetBytes(bucketName BucketName, key CacheKey, keyAppend ...interface{}) ([]byte, error) {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return nil, toolkitError.ErrBadBucketName
	}
	return bucket.GetBytes(key, keyAppend...)
}

// Put 缓存数据
func (c *CacheManager) Put(bucketName BucketName, key CacheKey, data any, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	bs, err := c.codec.Encode(data)
	if err != nil {
		return err
	}
	return bucket.PutBytes(key, bs, keyAppend...)
}

// PutBytes 缓存字节数组
func (c *CacheManager) PutBytes(bucketName BucketName, key CacheKey, data []byte, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	return bucket.PutBytes(key, data, keyAppend...)
}

// Evict 清除缓存
func (c *CacheManager) Evict(bucketName BucketName, key CacheKey, keyAppend ...interface{}) error {
	bucket := c.GetBucket(bucketName)
	if bucket == nil {
		logger.Logrus().Errorln("caching: bad bucketName", bucketName)
		return toolkitError.ErrBadBucketName
	}
	return bucket.Evict(key, keyAppend...)
}
