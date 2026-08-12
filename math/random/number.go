package random

import (
	"math/rand/v2"
)

// RandInt 生成指定范围内随机数 [0,max]
func RandInt(max int) int {
	if max < 0 {
		return -1
	}
	return int(rand.Uint64N(uint64(max) + 1))
}

// RandRangeInt 生成指定范围内随机数 [min,max]
func RandRangeInt(min, max int) int {
	if min > max {
		return -1
	}

	// 使用无符号运算计算区间宽度，避免有符号整数在边界处溢出。
	width := uint64(max) - uint64(min) + 1
	if width == 0 {
		// 64 位平台完整 int 区间的宽度是 2^64，无法用 uint64 表示。
		return int(rand.Uint64())
	}
	return int(uint64(min) + rand.Uint64N(width))
}
