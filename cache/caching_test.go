package cache

import (
	"fmt"
	"github.com/acexy/golang-toolkit/logger"
	"testing"
	"time"
)

var manager *CachingManager

type User struct {
	Name string
	Sex  uint8
}

func TestBigCache(t *testing.T) {

	manager = NewCacheBucketManager("b1", NewSimpleBigCache(time.Second*10))
	manager.AddBucket("b2", NewSimpleBigCache(time.Second*3))

	err := manager.Put("b1", "key1", "123")
	if err != nil {
		logger.Logrus().Errorln(err)
		return
	}
	err = manager.Put("b2", "key1", User{Name: "Q", Sex: 1})
	if err != nil {
		logger.Logrus().Errorln(err)
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
			var result User
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
