package random

import (
	"errors"
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
}
