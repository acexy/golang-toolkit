package cache

import (
	"fmt"
	"github.com/acexy/golang-toolkit/log"
	"testing"
	"time"
)

var manager *CachingManager

func TestBigCache(t *testing.T) {

	manager = NewCacheBucketManager("b1", NewSimpleBigCache(time.Second*10))

	manager.AddBucket("b2", NewSimpleBigCache(time.Second*3))

	err := manager.Put("b1", "key1", "123")
	if err != nil {
		log.Logrus().Errorln(err)
		return
	}
	err = manager.Put("b2", "key1", "321")
	if err != nil {
		log.Logrus().Errorln(err)
		return
	}

	go func() {
		for {
			var result string
			err := manager.Get("b1", "key1", &result)
			if err != nil {
				return
			}
			fmt.Println("b1", "key1", result)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	go func() {
		for {
			var result string
			err := manager.Get("b2", "key1", &result)
			if err != nil {
				return
			}
			fmt.Println("b2", "key1", result)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	time.Sleep(time.Minute)
}
