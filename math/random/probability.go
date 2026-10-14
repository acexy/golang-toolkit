package random

import (
	"math"
	"math/rand/v2"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/shopspring/decimal"
)

// ProbabilityTrue samples a percentage in [0, 100]; invalid or non-finite values return false.
func ProbabilityTrue(percentage float64) bool {
	if math.IsNaN(percentage) || math.IsInf(percentage, 0) || percentage < 0 || percentage > 100 {
		return false
	}
	return rand.Float64() < percentage/100
}

// ProbabilityResult samples a map key from finite nonnegative percentages totaling exactly 100.
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
