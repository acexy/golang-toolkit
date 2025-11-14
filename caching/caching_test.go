package caching

import (
	"fmt"
	"testing"
	"time"

	lgr "github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/sys"
)

var manager *CacheManager

type User struct {
	Name string
	Sex  uint8
}

func init() {
	manager = NewCacheBucketManager("b1", NewSimpleBigCache(time.Second*10))
	manager.AddBucket("b2", NewSimpleBigCache(time.Second*4))
}

func TestBigCache(t *testing.T) {
	err := manager.Put("b1", NewNemCacheKey("key1"), "123")
	if err != nil {
		lgr.Logrus().Errorln(err)
		return
	}
	err = manager.Put("b2", NewNemCacheKey("key1"), User{Name: "Q", Sex: 1})
	if err != nil {
		lgr.Logrus().Errorln(err)
		return
	}

	go func() {
		for {
			var result string
			err := manager.Get("b1", NewNemCacheKey("key1"), &result)
			if err != nil {
				return
			}
			fmt.Println("b1", "key1", result)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	go func() {
		for {
			var result User
			err := manager.Get("b2", NewNemCacheKey("key1"), &result)
			if err != nil {
				return
			}
			fmt.Println("b2", "key1", result)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	sys.ShutdownHolding()
}

func TestKeyFormat(t *testing.T) {
	key := NewNemCacheKey("key:%s")
	_ = manager.Put("b1", key, 3, "1")
	var result int
	_ = manager.Get("b1", key, &result, "1")
	fmt.Println(result)
}
