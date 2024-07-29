package slice

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

type Person struct {
	Name string
	Age  int
}

func TestToMap(t *testing.T) {
	// 示例1: 整数切片
	ints := []int{1, 2, 3, 4, 5, 6}
	intMap := ToMap(ints, func(t *int) (*string, *int, bool) {
		if *t > 3 {
			str := conversion.FromInt(*t)
			return &str, t, true
		}
		return nil, nil, false
	})
	fmt.Println(intMap)

	// 示例2: 字符串切片
	strings := []string{"apple", "banana", "cherry", "date"}
	stringMap := ToMap(strings, func(t *string) (*string, *string, bool) {
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
	personMap := ToMap(people, func(t *Person) (*string, *int, bool) {
		if t.Age > 30 {
			return &t.Name, &t.Age, true
		}
		return nil, nil, false
	})

	// 打印结果
	fmt.Println(personMap)
}
