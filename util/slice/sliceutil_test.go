package slice

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"apple", "banana", "cherry"}

	fmt.Println(Contains(intSlice, 3))           // 输出: true
	fmt.Println(Contains(intSlice, 6))           // 输出: false
	fmt.Println(Contains(stringSlice, "banana")) // 输出: true
	fmt.Println(Contains(stringSlice, "grape"))  // 输出: false
}