package random

import (
	"github.com/acexy/golang-toolkit/util/coll"
	"testing"
)

func TestProbabilityTrue(t *testing.T) {
	var count int
	for i := 0; i < 1000; i++ {
		if ProbabilityTrue(10) {
			count++
		}
	}
	print(count)
}

func TestProbabilityResult(t *testing.T) {
	result := map[any]int{
		"A": 0,
		"B": 0,
		"C": 0,
		"D": 0,
	}
	for i := 0; i < 10000; i++ {
		result[ProbabilityResult(map[any]float64{
			"A": 10.15,
			"B": 20.85,
			"C": 53.05,
			"D": 15.95,
		})]++
	}
	coll.MapForeachAll(result, func(k any, v int) {
		println(k.(string), v)
	})
}
