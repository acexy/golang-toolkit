package random

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandomInt(t *testing.T) {
	for range 20 {
		t.Log(RandInt(0))
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(rand.Intn(2))
	}
}
