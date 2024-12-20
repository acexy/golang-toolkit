package coll

import (
	"fmt"
	"github.com/acexy/golang-toolkit/util/json"
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

type User struct {
	Name string
	Age  int
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
	users := map[User]int{
		User{"one", 1}: 0,
		User{"two", 1}: 0,
		User{"two", 2}: 0,
		User{"two", 1}: 0,
	}
	userSeys := MapKeyToSlice(users)
	fmt.Printf("取出的元素: %v\n", userSeys)
}

type U struct {
	K string
	V int
}

func TestMapToSlice(t *testing.T) {
	exampleMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
	result := MapFilterToSlice(exampleMap, func(key string, value int) (U, bool) {
		if value > 3 {
			return U{}, false
		}
		return U{key, value}, true
	})
	fmt.Printf("取出的元素: %v\n", json.ToJson(result))
}
