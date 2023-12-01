package conversion

import (
	"fmt"
	"testing"
)

func a() {
	binaryNumber := int64(0b0101010)
	fmt.Printf("Binary: %d\n", binaryNumber)
}

func Test1(t *testing.T) {
	a()
}
