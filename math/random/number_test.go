package random

import "testing"

func TestRandomInt(t *testing.T) {

	for range 20 {
		t.Log(RandInt(0))
	}
}
