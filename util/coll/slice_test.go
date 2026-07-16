package coll

import (
	"errors"
	"slices"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

type person struct {
	name string
	age  int
}

func TestSliceRandom(t *testing.T) {
	if got := SliceRandom([]int{10}); got != 10 {
		t.Fatalf("SliceRandom() = %d", got)
	}
	if got := SliceRandom([]int{}); got != 0 {
		t.Fatalf("SliceRandom(empty) = %d", got)
	}
}

func TestSliceContainsAndIndex(t *testing.T) {
	values := []int{1, 2, 3}
	if !SliceContains(values, 2) {
		t.Fatal("expected contains 2")
	}
	if SliceContains(values, 4) {
		t.Fatal("should not contain 4")
	}
	if index := SliceIndexOf(values, 3); index != 2 {
		t.Fatalf("SliceIndexOf() = %d", index)
	}

	people := []*person{{name: "a", age: 1}}
	if !SliceContains(people, &person{name: "a", age: 1}, func(a, b *person) bool {
		return a.name == b.name && a.age == b.age
	}) {
		t.Fatal("expected custom contains")
	}
}

func TestSliceContainsByAndIndexBy(t *testing.T) {
	people := []person{{name: "a", age: 1}, {name: "b", age: 2}}
	if !SliceContainsBy(people, func(p person) bool {
		return p.name == "b"
	}) {
		t.Fatal("expected contains by")
	}
	if index := SliceIndexBy(people, func(p person) bool {
		return p.name == "b"
	}); index != 1 {
		t.Fatalf("SliceIndexBy() = %d", index)
	}
}

func TestSliceFind(t *testing.T) {
	values := []int{1, 2, 3, 2}
	got, ok := SliceFind(values, func(v int) bool { return v == 2 })
	if !ok || got != 2 {
		t.Fatalf("SliceFind() = %d, %v", got, ok)
	}
	got, ok = SliceFindLast(values, func(v int) bool { return v == 2 })
	if !ok || got != 2 {
		t.Fatalf("SliceFindLast() = %d, %v", got, ok)
	}
}

func TestSliceFilter(t *testing.T) {
	got := SliceFilter([]int{1, 2, 3}, func(v int) bool { return v > 1 })
	if !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("unexpected filter result: %v", got)
	}
}

func TestSliceSetOperations(t *testing.T) {
	if got := SliceIntersection([]int{1, 2, 2, 3}, []int{2, 3, 4}); !slices.Equal(got, []int{2, 3}) {
		t.Fatalf("unexpected intersection: %v", got)
	}
	if got := SliceUnion([]int{1, 2}, []int{2, 3}); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("unexpected union: %v", got)
	}
	if got := SliceComplement([]int{1, 2, 3}, []int{2}); !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("unexpected complement: %v", got)
	}
}

func TestSliceDiff(t *testing.T) {
	added, removed := SliceDiff([]int{1, 2, 3}, []int{2, 3, 4})
	if !slices.Equal(added, []int{4}) || !slices.Equal(removed, []int{1}) {
		t.Fatalf("unexpected diff added=%v removed=%v", added, removed)
	}

	oldPeople := []person{{name: "a", age: 1}, {name: "b", age: 2}}
	newPeople := []person{{name: "b", age: 20}, {name: "c", age: 3}}
	addedPeople, removedPeople := SliceDiff(oldPeople, newPeople, func(a, b person) bool {
		return a.name == b.name
	})
	if len(addedPeople) != 1 || addedPeople[0].name != "c" || len(removedPeople) != 1 || removedPeople[0].name != "a" {
		t.Fatalf("unexpected custom diff added=%v removed=%v", addedPeople, removedPeople)
	}
}

func TestSliceAnyDiff(t *testing.T) {
	oldValues := [][]int{{1}, {2}}
	newValues := [][]int{{2}, {3}}
	added, removed := SliceAnyDiff(oldValues, newValues, func(a, b []int) bool {
		return slices.Equal(a, b)
	})
	if len(added) != 1 || !slices.Equal(added[0], []int{3}) || len(removed) != 1 || !slices.Equal(removed[0], []int{1}) {
		t.Fatalf("unexpected any diff added=%v removed=%v", added, removed)
	}
}

func TestSliceIsSubset(t *testing.T) {
	if !SliceIsSubset([]string{"a", "b"}, []string{"a", "b", "c"}) {
		t.Fatal("expected subset")
	}
	if SliceIsSubset([]string{"a", "d"}, []string{"a", "b", "c"}) {
		t.Fatal("should not be subset")
	}
}

func TestSliceFilterToMap(t *testing.T) {
	got := SliceFilterToMap([]int{1, 2, 3}, func(v int) (int, string, bool) {
		return v, "ok", v > 1
	})
	if len(got) != 2 || got[2] != "ok" || got[3] != "ok" {
		t.Fatalf("unexpected map: %v", got)
	}
}

func TestSliceCollect(t *testing.T) {
	got := SliceCollect([]int{1, 2}, func(v int) string {
		if v == 1 {
			return "one"
		}
		return "two"
	})
	if !slices.Equal(got, []string{"one", "two"}) {
		t.Fatalf("unexpected collect result: %v", got)
	}
}

func TestSliceFilterCollect(t *testing.T) {
	got := SliceFilterCollect([]int{1, 2, 3}, func(v int) (int, bool) {
		return v * 10, v > 1
	})
	if !slices.Equal(got, []int{20, 30}) {
		t.Fatalf("unexpected filter collect result: %v", got)
	}
}

func TestSliceFlat(t *testing.T) {
	got := SliceFlat([][]int{{1, 2}, {}, {3}})
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("unexpected flat result: %v", got)
	}
	got = SliceFlat[int](nil)
	if got == nil || len(got) != 0 {
		t.Fatalf("SliceFlat(nil) = %v, want empty slice", got)
	}
}

func TestSliceForEach(t *testing.T) {
	count := 0
	SliceForEach([]int{1, 2, 3}, func(v int) bool {
		count++
		return false
	})
	if count != 1 {
		t.Fatalf("SliceForEach should stop after first item, got %d", count)
	}

	count = 0
	SliceForEachAll([]int{1, 2, 3}, func(v int) {
		count++
	})
	if count != 3 {
		t.Fatalf("SliceForEachAll count = %d", count)
	}
}

func TestSliceDistinct(t *testing.T) {
	if got := SliceDistinct([]int{1, 2, 1, 3, 2}); !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("unexpected distinct result: %v", got)
	}

	people := []person{{name: "a", age: 1}, {name: "b", age: 2}, {name: "a", age: 3}}
	got := SliceDistinctBy(people, func(p person) string {
		return p.name
	})
	if len(got) != 2 || got[0].age != 1 || got[1].name != "b" {
		t.Fatalf("unexpected distinct by result: %v", got)
	}
}

func TestSliceSort(t *testing.T) {
	values := []int{3, 1, 2}
	SliceSort(values, func(v int) int { return v })
	if !slices.Equal(values, []int{1, 2, 3}) {
		t.Fatalf("unexpected ascending sort: %v", values)
	}
	SliceSort(values, func(v int) int { return v }, true)
	if !slices.Equal(values, []int{3, 2, 1}) {
		t.Fatalf("unexpected descending sort: %v", values)
	}
}

func TestSliceGroupBy(t *testing.T) {
	people := []person{{name: "a", age: 1}, {name: "a", age: 2}, {name: "b", age: 3}}
	group := SliceGroupBy(people, func(p person) string { return p.name })
	if len(group["a"]) != 2 || len(group["b"]) != 1 {
		t.Fatalf("unexpected group: %v", group)
	}

	single := SliceGroupBySingle(people, func(p person) string { return p.name })
	if single["a"].age != 2 {
		t.Fatalf("unexpected single group: %v", single)
	}

	valueGroup := SliceGroupByValue(people, func(p person) (string, int) {
		return p.name, p.age
	})
	if !slices.Equal(valueGroup["a"], []int{1, 2}) {
		t.Fatalf("unexpected value group: %v", valueGroup)
	}

	valueSingle := SliceGroupBySingleValue(people, func(p person) (string, int) {
		return p.name, p.age
	})
	if valueSingle["a"] != 2 {
		t.Fatalf("unexpected value single group: %v", valueSingle)
	}
}

func TestSliceSplitChunk(t *testing.T) {
	got := SliceSplitChunk([]int{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 || !slices.Equal(got[0], []int{1, 2}) || !slices.Equal(got[2], []int{5}) {
		t.Fatalf("unexpected chunks: %v", got)
	}
	if got := SliceSplitChunk([]int{1}, 0); got != nil {
		t.Fatalf("expected nil chunks, got %v", got)
	}
}

func TestSliceRemove(t *testing.T) {
	got, err := SliceRemove(1, []int{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 3}) {
		t.Fatalf("unexpected remove result: %v", got)
	}
	if _, err = SliceRemove(3, []int{1}); !errors.Is(err, toolkitError.ErrSliceIndexOutOfRange) {
		t.Fatalf("expected ErrSliceIndexOutOfRange, got %v", err)
	}
}

func TestSliceRemoveSafe(t *testing.T) {
	source := []int{1, 2, 3}
	got, err := SliceRemoveSafe(1, source)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 3}) || !slices.Equal(source, []int{1, 2, 3}) {
		t.Fatalf("unexpected safe remove result=%v source=%v", got, source)
	}
}

func TestSliceInsert(t *testing.T) {
	got, err := SliceInsert(1, 9, []int{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 9, 2}) {
		t.Fatalf("unexpected insert result: %v", got)
	}
	if _, err = SliceInsert(3, 9, []int{1}); !errors.Is(err, toolkitError.ErrSliceIndexOutOfRange) {
		t.Fatalf("expected ErrSliceIndexOutOfRange, got %v", err)
	}
}

func TestSliceInsertSafe(t *testing.T) {
	source := []int{1, 2}
	got, err := SliceInsertSafe(1, 9, source)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 9, 2}) || !slices.Equal(source, []int{1, 2}) {
		t.Fatalf("unexpected safe insert result=%v source=%v", got, source)
	}
}

func TestSliceInserts(t *testing.T) {
	got, err := SliceInserts(1, []int{8, 9}, []int{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 8, 9, 2}) {
		t.Fatalf("unexpected inserts result: %v", got)
	}
}

func TestSliceInsertsSafe(t *testing.T) {
	source := []int{1, 2}
	got, err := SliceInsertsSafe(1, []int{8, 9}, source)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 8, 9, 2}) || !slices.Equal(source, []int{1, 2}) {
		t.Fatalf("unexpected safe inserts result=%v source=%v", got, source)
	}
	got, err = SliceInsertsSafe(1, nil, source)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, source) || &got[0] == &source[0] {
		t.Fatalf("expected copied source when inserting empty slice")
	}
}
