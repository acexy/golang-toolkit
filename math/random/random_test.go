package random

import (
	"fmt"
	"testing"
)

func TestRandRangeInt(t *testing.T) {
	fmt.Println(RandRangeInt(1, 2))
}

func TestRandString(t *testing.T) {
	fmt.Println(RandString(5))
}
