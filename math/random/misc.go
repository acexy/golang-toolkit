package random

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	// 使用 sync.Pool 提供每个 Goroutine 独立的随机数生成器
	randProbabilityPool = sync.Pool{
		New: func() interface{} {
			seed := time.Now().UnixNano()
			return rand.New(rand.NewSource(seed))
		},
	}
)

// decimalPlaces 计算浮点数的有效小数位数
func decimalPlaces(num float64) int {
	str := fmt.Sprintf("%f", num)
	decimalPart := str[strings.Index(str, ".")+1:]
	return len(strings.TrimRight(decimalPart, "0"))
}

// ProbabilityTrue 支持任意小数精度的概率判断
// 0 <= percentage <= 100 超出范围将永远返回 false。
func ProbabilityTrue(percentage float64) bool {
	if percentage < 0 || percentage > 100 {
		return false
	}
	scale := math.Pow(10, float64(decimalPlaces(percentage)))
	maxValue := int(100 * scale)
	threshold := int(percentage * scale)
	r := randProbabilityPool.Get().(*rand.Rand)
	defer randProbabilityPool.Put(r)
	randomNumber := r.Intn(maxValue)
	return randomNumber < threshold
}
