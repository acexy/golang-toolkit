package coll

import (
	"fmt"
	"github.com/acexy/golang-toolkit/math/conversion"
	"github.com/acexy/golang-toolkit/util/str"
	"testing"
)

type people struct {
	name string
	age  int
}

func TestSliceContains(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"apple", "banana", "cherry"}

	fmt.Println(SliceContains(intSlice, 3))           // 输出: true
	fmt.Println(SliceContains(intSlice, 6))           // 输出: false
	fmt.Println(SliceContains(stringSlice, "banana")) // 输出: true
	fmt.Println(SliceContains(stringSlice, "grape"))  // 输出: false

	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	fmt.Println(SliceContains(peoples, people{name: "张三", age: 29}, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
	list := []string{"US"}
	fmt.Println(SliceContains(list, ""))
}

func TestSliceFilter(t *testing.T) {
	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}

	fmt.Println(SliceFilter(peoples, func(item *people) bool {
		return item.age == 20
	}))
}
func TestSliceIntersection(t *testing.T) {
	intSlice1 := []int{3, 4, 4, 5, 6, 7}
	intSlice2 := []int{1, 2, 3, 4, 3, 5}
	fmt.Println(SliceIntersection(intSlice1, intSlice2))

	intSlice3 := []int{1, 2, 3, 4, 5}
	intSlice4 := []int{6, 6, 7, 7, 8, 9, 0}
	fmt.Println(SliceIntersection(intSlice3, intSlice4))
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
	fmt.Println(SliceIntersection(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestSliceUnion(t *testing.T) {
	intSlice1 := []int{1, 2, 3, 4, 3, 5}
	intSlice2 := []int{3, 4, 4, 5, 6, 7}
	fmt.Println(SliceUnion(intSlice1, intSlice2))

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
	fmt.Println(SliceUnion(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestSliceComplement(t *testing.T) {
	intSlice1 := []int{1, 2, 3, 4, 3, 5}
	intSlice2 := []int{3, 4, 4, 5}
	fmt.Println(SliceComplement(intSlice1, intSlice2))
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
	fmt.Println(SliceComplement(peoples1, peoples2, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

type Person struct {
	Name string
	Age  int
}

func TestSliceToMap(t *testing.T) {
	// 示例1: 整数切片
	ints := []int{1, 2, 3, 4, 5, 6}
	intMap := SliceToMap(ints, func(t *int) (*string, *int, bool) {
		if *t > 3 {
			str := conversion.FromInt(*t)
			return &str, t, true
		}
		return nil, nil, false
	})
	fmt.Println(intMap)

	// 示例2: 字符串切片
	strings := []string{"apple", "banana", "cherry", "date"}
	stringMap := SliceToMap(strings, func(t *string) (*string, *string, bool) {
		if str.CharLength(*t) > 4 {
			return t, t, true
		}
		return nil, nil, false
	})
	fmt.Println(stringMap)

	// 定义一个结构体切片
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Dave", Age: 40},
	}

	// 调用通用方法
	personMap := SliceToMap(people, func(t *Person) (*string, *int, bool) {
		if t.Age > 30 {
			return &t.Name, &t.Age, true
		}
		return nil, nil, false
	})

	// 打印结果
	fmt.Println(personMap)
}
