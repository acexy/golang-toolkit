package coll

import (
	"fmt"
	"testing"
)

func TestMapFirst(t *testing.T) {
	exampleMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	key, value := MapFirst(exampleMap)
	fmt.Printf("取出的元素: key=%v, value=%v\n", key, value)
	fmt.Printf("取出的元素: key=%v, value=%v\n", key, value)
	fmt.Printf("取出的元素: key=%v, value=%v\n", key, value)
}

func TestMapKeyToSlice(t *testing.T) {
	exampleMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
	keys := MapKeyToSlice(exampleMap)
	fmt.Printf("取出的元素: %v\n", keys)
}
