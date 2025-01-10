package random

import (
	"testing"
)

func TestProbabilityTrue(t *testing.T) {
	var count int
	for i := 0; i < 1000; i++ {
		if ProbabilityTrue(50) {
			count++
		}
	}
	print(count)
}
