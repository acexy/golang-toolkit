package slice

import (
	"fmt"
	"testing"
)

type people struct {
	name string
	age  int
}

func TestContains(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"apple", "banana", "cherry"}

	fmt.Println(Contains(intSlice, 3))           // 输出: true
	fmt.Println(Contains(intSlice, 6))           // 输出: false
	fmt.Println(Contains(stringSlice, "banana")) // 输出: true
	fmt.Println(Contains(stringSlice, "grape"))  // 输出: false

	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	fmt.Println(Contains(peoples, people{name: "张三", age: 29}, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestFilter(t *testing.T) {
	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}

	fmt.Println(Filter(peoples, func(item *people) bool {
		return item.age == 20
	}))
}
func TestIntersection(t *testing.T) {
	intSlice1 := []int{3, 4, 4, 5, 6, 7}
	intSlice2 := []int{1, 2, 3, 4, 3, 5}
	fmt.Println(Intersection(intSlice1, intSlice2))

	intSlice3 := []int{1, 2, 3, 4, 5}
	intSlice4 := []int{6, 6, 7, 7, 8, 9, 0}
	fmt.Println(Intersection(intSlice3, intSlice4))
	peoples1 := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	peoples2 := []people{
		{name: "李四", age: 20},
		{name: "赵六", age: 20},
	}
	fmt.Println(Intersection(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestUnion(t *testing.T) {
	intSlice1 := []int{1, 2, 3, 4, 3, 5}
	intSlice2 := []int{3, 4, 4, 5, 6, 7}
	fmt.Println(Union(intSlice1, intSlice2))

	peoples1 := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
		{name: "赵六", age: 21},
	}
	peoples2 := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "赵六", age: 20},
		{name: "赵六", age: 21},
		{name: "赵六", age: 21},
	}
	fmt.Println(Union(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestComplement(t *testing.T) {
	intSlice1 := []int{1, 2, 3, 4, 3, 5}
	intSlice2 := []int{3, 4, 4, 5}
	fmt.Println(Complement(intSlice1, intSlice2))
	peoples1 := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
		{name: "赵六", age: 21},
	}
	peoples2 := []people{
		{name: "张三", age: 28},
		{name: "赵六", age: 20},
		{name: "赵六", age: 21},
		{name: "赵六", age: 21},
	}
	fmt.Println(Complement(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}
