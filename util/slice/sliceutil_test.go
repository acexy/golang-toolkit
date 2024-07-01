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

type people struct {
	name string
	age  int
}

func TestFilterWithFn(t *testing.T) {
	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}

	fmt.Println(FilterWithFn(&peoples, func(item *people) bool {
		return item.age == 20
	}))
}
