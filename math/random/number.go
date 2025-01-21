package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	// 使用 sync.Pool 提供每个 Goroutine 独立的随机数生成器
	randPool = sync.Pool{
		New: func() interface{} {
			seed := time.Now().UnixNano()
			return rand.New(rand.NewSource(seed))
		},
	}
)

// RandInt 生成指定范围内随机数 [0,max]
func RandInt(max int) int {
	if max < 0 {
		return -1
	}
	r := randPool.Get().(*rand.Rand)
	defer randPool.Put(r)
	return r.Intn(max + 1)
}

// RandRangeInt 生成指定范围内随机数 [min,max]
func RandRangeInt(min, max int) int {
	if min > max {
		return -1
	}
	r := randPool.Get().(*rand.Rand)
	defer randPool.Put(r)
	return r.Intn(max-min+1) + min
}
