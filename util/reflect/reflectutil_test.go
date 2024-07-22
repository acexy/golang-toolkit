package reflect

import (
	"fmt"
	"testing"
)

func TestNonZeroField(t *testing.T) {
	i := 1
	testStruct := struct {
		A string
		B *int
		C bool
		D int
	}{
		A: "a",
		B: &i,
		C: true,
	}
	fields, err := NonZeroField(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields)
}
