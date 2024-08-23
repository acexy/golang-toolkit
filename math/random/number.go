package random

import (
	"math/rand"
)

// RandInt 生成指定范围内随机数 [0,max]
func RandInt(max int) int {
	return rand.Intn(max + 1)
}

// RandRangeInt 生成指定范围内随机数 [min,max]
func RandRangeInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

type SeedRandom struct {
	rand *rand.Rand
}

// NewRandom 创建随机数生成器 指定随机算法的种子
func NewRandom(seed int64) *SeedRandom {
	return &SeedRandom{rand: rand.New(rand.NewSource(seed))}
}

// RandInt 生成指定范围内随机数 [0,max]
func (s *SeedRandom) RandInt(max int) int {
	return s.rand.Intn(max)
}

// RandRangeInt 生成指定范围内随机数 [min,max]
func (s *SeedRandom) RandRangeInt(min, max int) int {
	return s.rand.Intn(max-min+1) + min
}
