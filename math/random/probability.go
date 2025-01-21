package random

import (
	"errors"
	"fmt"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/shopspring/decimal"
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

// ProbabilityTrue 随机执行指定概率(0-100%)返回true的计算
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

// ProbabilityResult 按照设置的各种key(概率结果)执行随机返回对应的key(发生的结果)
// 所有key的概率之和必须为100%
func ProbabilityResult(percentage map[any]float64) (any, error) {
	if len(percentage) == 0 {
		return nil, errors.New("empty param")
	}
	total := decimal.Zero
	for _, v := range percentage {
		total = total.Add(decimal.NewFromFloat(v))
	}
	if total.Compare(decimal.NewFromFloat(100)) != 0 {
		return nil, errors.New("percentage is not 100%")
	}
	var maxScale int
	coll.MapForeachAll(percentage, func(k any, v float64) {
		maxScale = int(math.Max(float64(maxScale), float64(decimalPlaces(v))))
	})
	calcPercentage := coll.MapCollect(percentage, func(k any, v float64) (any, decimal.Decimal) {
		return k, decimal.NewFromFloat(v).Mul(decimal.NewFromFloat(math.Pow(10, float64(maxScale))))
	})
	sum := decimal.Zero
	coll.MapForeachAll(calcPercentage, func(k any, v decimal.Decimal) {
		sum = sum.Add(v)
	})
	if sum.Compare(decimal.NewFromFloat(math.Pow(10, float64(maxScale))*100)) != 0 {
		return nil, errors.New("percentage is not 100%")
	}
	randomValue := RandInt(int(sum.IntPart() - 1)) // 随机数范围为 [0, total)
	decimalRandomValue := decimal.NewFromInt(int64(randomValue))
	cumulative := decimal.Zero
	for key, scaledProb := range calcPercentage {
		cumulative = cumulative.Add(scaledProb)
		if decimalRandomValue.Compare(cumulative) < 0 {
			return key, nil
		}
	}
	return nil, errors.New("")
}
