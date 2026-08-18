package random

import (
	"errors"
	"math"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func TestProbabilityTrue(t *testing.T) {
	if ProbabilityTrue(-1) {
		t.Fatal("negative probability should be false")
	}
	if ProbabilityTrue(101) {
		t.Fatal("probability greater than 100 should be false")
	}
	if ProbabilityTrue(0) {
		t.Fatal("zero probability should be false")
	}
	if !ProbabilityTrue(100) {
		t.Fatal("100 probability should be true")
	}
	if ProbabilityTrue(math.NaN()) || ProbabilityTrue(math.Inf(1)) {
		t.Fatal("non-finite probability should be false")
	}
}

func TestProbabilityResult(t *testing.T) {
	result, err := ProbabilityResult(map[any]float64{
		"A": 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "A" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestProbabilityResultError(t *testing.T) {
	if _, err := ProbabilityResult(nil); !errors.Is(err, toolkitError.ErrEmptyProbability) {
		t.Fatalf("expected ErrEmptyProbability, got %v", err)
	}
	if _, err := ProbabilityResult(map[any]float64{"A": 50}); !errors.Is(err, toolkitError.ErrInvalidProbabilityTotal) {
		t.Fatalf("expected ErrInvalidProbabilityTotal, got %v", err)
	}
	invalidValues := []map[any]float64{
		{"A": -10, "B": 110},
		{"A": math.NaN(), "B": 100},
		{"A": math.Inf(1), "B": 100},
	}
	for _, values := range invalidValues {
		if _, err := ProbabilityResult(values); !errors.Is(err, toolkitError.ErrInvalidProbabilityTotal) {
			t.Fatalf("expected ErrInvalidProbabilityTotal, got %v", err)
		}
	}
}
