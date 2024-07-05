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

type seedRandom struct {
	rand *rand.Rand
}

func NewRandom(seed int64) *seedRandom {
	return &seedRandom{rand: rand.New(rand.NewSource(seed))}
}

func (s *seedRandom) RandInt(max int) int {
	return s.rand.Intn(max)
}

func (s *seedRandom) RandRangeInt(min, max int) int {
	return s.rand.Intn(max-min+1) + min
}
