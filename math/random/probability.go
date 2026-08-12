package random

import (
	"math"
	"math/rand/v2"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/shopspring/decimal"
)

// ProbabilityTrue 随机执行指定概率(0-100%)返回true的计算
// 0 <= percentage <= 100 超出范围将永远返回 false。
func ProbabilityTrue(percentage float64) bool {
	if math.IsNaN(percentage) || math.IsInf(percentage, 0) || percentage < 0 || percentage > 100 {
		return false
	}
	return rand.Float64() < percentage/100
}

// ProbabilityResult 按照设置的各种key(概率结果)执行随机返回对应的key(发生的结果)
// 所有key的概率之和必须为100%
func ProbabilityResult(percentage map[any]float64) (any, error) {
	if len(percentage) == 0 {
		return nil, toolkitError.ErrEmptyProbability
	}
	total := decimal.Zero
	for _, v := range percentage {
		if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > 100 {
			return nil, toolkitError.ErrInvalidProbabilityTotal
		}
		total = total.Add(decimal.NewFromFloat(v))
	}
	if !total.Equal(decimal.NewFromInt(100)) {
		return nil, toolkitError.ErrInvalidProbabilityTotal
	}

	randomValue := decimal.NewFromFloat(rand.Float64() * 100)
	cumulative := decimal.Zero
	for key, value := range percentage {
		cumulative = cumulative.Add(decimal.NewFromFloat(value))
		if randomValue.LessThan(cumulative) {
			return key, nil
		}
	}
	return nil, toolkitError.ErrProbabilityResultNotFound
}
