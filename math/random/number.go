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

type NumberRandom struct {
	rand *rand.Rand
}

func NewNumberRandom(seed int64) *NumberRandom {
	return &NumberRandom{rand: rand.New(rand.NewSource(seed))}
}

func (r *NumberRandom) RandInt(max int) int {
	return r.rand.Intn(max)
}

func (r *NumberRandom) RandRangeInt(min, max int) int {
	return r.rand.Intn(max-min+1) + min
}
