package coll

import (
	"fmt"
	"testing"

	"github.com/acexy/golang-toolkit/math/conversion"
	"github.com/acexy/golang-toolkit/util/str"
)

type people struct {
	name string
	age  int
}

func TestSliceContains(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"apple", "banana", "cherry"}

	fmt.Println(SliceContains(intSlice, 3)) // 输出: true
	fmt.Println(SliceIndexOf(intSlice, 3))
	fmt.Println(SliceContains(intSlice, 6)) // 输出: false
	fmt.Println(SliceIndexOf(intSlice, 6))
	fmt.Println(SliceContains(stringSlice, "banana")) // 输出: true
	fmt.Println(SliceContains(stringSlice, "grape"))  // 输出: false

	peoples := []*people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	fmt.Println(SliceContains(peoples, &people{name: "张三", age: 28}, func(a, b *people) bool {
		return a.name == b.name && a.age == b.age
	}))
	list := []string{"US", ""}
	fmt.Println(SliceContains(list, ""))
}

func TestSliceAnyContains(t *testing.T) {
	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	fmt.Println(SliceAnyContains(peoples, func(a people) bool {
		return a.name == "张三"
	}))
	fmt.Println(SliceAnyIndexOf(peoples, func(a people) bool {
		return a.name == "张三"
	}))
}

func TestSliceFilter(t *testing.T) {
	peoples := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}

	peoplesPt := []*people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
	}
	fmt.Println(SliceFilter(peoples, func(item people) bool {
		return item.age == 20
	}))
	fmt.Println(SliceFilterFirstOne(peoplesPt, func(item *people) bool {
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
	fmt.Println(SliceIntersection(peoples1, peoples2, func(a, b people) bool {
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
	fmt.Println(SliceUnion(peoples1, peoples2, func(a, b people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

func TestSliceSort(t *testing.T) {
	peoples1 := []people{
		{name: "张三", age: 28},
		{name: "李四", age: 20},
		{name: "王五", age: 22},
		{name: "赵六", age: 20},
		{name: "赵六", age: 21},
	}
	fmt.Println(peoples1)
	SliceSort(peoples1, func(e people) int {
		return e.age
	}, true)
	fmt.Println(peoples1)

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
	fmt.Println(SliceComplement(peoples1, peoples2, func(a, b people) bool {
		return a.name == b.name && a.age == b.age
	}))
}

type Person struct {
	Name string
	Age  int
}

func TestSliceFilterToMap(t *testing.T) {
	// 示例1: 整数切片
	ints := []int{1, 2, 3, 4, 5, 6}
	intMap := SliceFilterToMap(ints, func(t int) (string, int, bool) {
		if t > 3 {
			fromInt := conversion.FromInt(t)
			return fromInt, t, true
		}
		return "", 0, false
	})
	fmt.Println(intMap)

	// 示例2: 字符串切片
	strings := []string{"apple", "banana", "cherry", "date"}
	stringMap := SliceFilterToMap(strings, func(t string) (string, string, bool) {
		if str.CharLength(t) > 4 {
			return t, t, true
		}
		return "", "", false
	})
	fmt.Println(stringMap)

	// 定义一个结构体切片
	people := []*Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Dave", Age: 40},
	}

	// 调用通用方法
	personMap := SliceFilterToMap(people, func(t *Person) (string, int, bool) {
		if t.Age > 30 {
			return t.Name, t.Age, true
		}
		return "", 0, false
	})

	// 打印结果
	fmt.Println(personMap)
}

func TestSliceCollect(t *testing.T) {
	// 输入切片
	input := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}

	// 定义映射函数，用于从 Person 提取 Age
	collect := func(p Person) int {
		return p.Age
	}

	// 调用 SliceCollect
	output := SliceCollect(input, collect)
	fmt.Println(output)

}

func TestSliceDistinct(t *testing.T) {
	input := []Person{
		{Name: "Alice", Age: 35},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Charlie", Age: 35},
	}
	fmt.Println(SliceDistinct(input))

	fmt.Println(SliceDistinctAny(input, func(t Person) int {
		return t.Age
	}))
}
func TestSliceGroupBy(t *testing.T) {
	input := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
		{Name: "Charlie", Age: 36},
	}
	fmt.Println(SliceGroupBy(input, func(t Person) string {
		return t.Name
	}))
	fmt.Println(SliceAnyGroupBy(input, func(t Person) (string, int) {
		return t.Name, t.Age
	}))
}

func TestSliceGroupBySingle(t *testing.T) {
	input := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie1", Age: 35},
		{Name: "Charlie", Age: 36},
	}
	fmt.Println(SliceGroupBySingle(input, func(t Person) string {
		return t.Name
	}))
	fmt.Println(SliceAnyGroupBySingle(input, func(t Person) (string, int) {
		return t.Name, t.Age
	}))
}

func TestSliceIsSubset(t *testing.T) {
	set := []string{"a", "b", "c", "d"}
	subset1 := []string{"a", "c"}
	subset2 := []string{"a", "e"}

	fmt.Println(SliceIsSubset(subset1, set)) // true
	fmt.Println(SliceIsSubset(subset2, set)) // false
}

func TestSliceDiff(t *testing.T) {
	old := []int{1, 2, 3, 4}
	new := []int{2, 3, 5, 6}
	added, removed := SliceDiff(old, new)
	fmt.Printf("Added: %v, Removed: %v\n", added, removed)
	// 输出: Added: [5 6], Removed: [1 4]

	// 示例2：结构体类型（假设CronConfig实现了comparable）
	type Person struct {
		ID   int
		Name string
	}

	oldPeople := []Person{{1, "Alice"}, {2, "Bob"}}
	newPeople := []Person{{2, "Bob"}, {3, "Charlie"}}
	addedPeople, removedPeople := SliceDiff(oldPeople, newPeople)
	fmt.Printf("Added: %v, Removed: %v\n", addedPeople, removedPeople)

	// 示例3：使用自定义比较函数（比较ID）
	addedByID, removedByID := SliceDiff(oldPeople, newPeople, func(a, b Person) bool {
		return a.ID == b.ID
	})
	fmt.Printf("Added by ID: %v, Removed by ID: %v\n", addedByID, removedByID)
}
