package coll

import (
	"slices"
	"testing"
)

func TestMapAny(t *testing.T) {
	key, value := MapAny(map[string]int{"one": 1})
	if key != "one" || value != 1 {
		t.Fatalf("MapAny() = %s, %d", key, value)
	}

	emptyKey, emptyValue := MapAny(map[string]int{})
	if emptyKey != "" || emptyValue != 0 {
		t.Fatalf("MapAny(empty) = %s, %d", emptyKey, emptyValue)
	}
}

func TestMapKeysAndValues(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2}

	keys := MapKeys(m)
	slices.Sort(keys)
	if !slices.Equal(keys, []string{"one", "two"}) {
		t.Fatalf("unexpected keys: %v", keys)
	}

	values := MapValues(m)
	slices.Sort(values)
	if !slices.Equal(values, []int{1, 2}) {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestMapCollect(t *testing.T) {
	result := MapCollect(map[string]int{"one": 1}, func(k string, v int) (int, string) {
		return v, k
	})
	if result[1] != "one" {
		t.Fatalf("unexpected collect result: %v", result)
	}
}

func TestMapFilterCollect(t *testing.T) {
	result := MapFilterCollect(map[string]int{"one": 1, "two": 2}, func(k string, v int) (string, int, bool) {
		return k, v, v > 1
	})
	if len(result) != 1 || result["two"] != 2 {
		t.Fatalf("unexpected filter collect result: %v", result)
	}
}

func TestMapFilterToSlice(t *testing.T) {
	result := MapFilterToSlice(map[string]int{"one": 1, "two": 2}, func(k string, v int) (string, bool) {
		return k, v > 1
	})
	if len(result) != 1 || result[0] != "two" {
		t.Fatalf("unexpected filter result: %v", result)
	}
}

func TestMapForEach(t *testing.T) {
	count := 0
	MapForEach(map[string]int{"one": 1, "two": 2}, func(k string, v int) bool {
		count++
		return false
	})
	if count != 1 {
		t.Fatalf("MapForEach should stop after first item, got %d", count)
	}

	count = 0
	MapForEachAll(map[string]int{"one": 1, "two": 2}, func(k string, v int) {
		count++
	})
	if count != 2 {
		t.Fatalf("MapForEachAll count = %d", count)
	}
}

func TestMapMerge(t *testing.T) {
	target := map[string]int{}
	target = MapMerge(target, map[string]int{"one": 1})
	if target["one"] != 1 {
		t.Fatalf("unexpected merge result: %v", target)
	}

	var nilTarget map[string]int
	nilTarget = MapMerge(nilTarget, map[string]int{"one": 1})
	if nilTarget["one"] != 1 {
		t.Fatalf("unexpected nil target merge result: %v", nilTarget)
	}
}

func TestMapRandom(t *testing.T) {
	key, value := MapRandom(map[string]int{"one": 1})
	if key != "one" || value != 1 {
		t.Fatalf("MapRandom() = %s, %d", key, value)
	}
}
