package random

import (
	"math/rand"
	"time"
)

func loadSeed() {
	rand.Seed(time.Now().UnixNano())
}

// RandInt 生成指定范围内随机数 [0,max]
func RandInt(max int) int {
	loadSeed()
	return rand.Intn(max + 1)
}

// RandRangeInt 生成指定范围内随机数 [min,max]
func RandRangeInt(min, max int) int {
	loadSeed()
	return rand.Intn(max-min+1) + min
}
