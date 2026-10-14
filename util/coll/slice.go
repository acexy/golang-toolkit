package coll

import (
	"cmp"
	"math/rand"
	"slices"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

// SliceRandom returns a pseudo-random element, or the zero value for an empty slice.
func SliceRandom[T any](slice []T) T {
	if len(slice) == 0 {
		var zeroValue T
		return zeroValue
	}
	return slice[rand.Intn(len(slice))]
}

// SliceContains reports membership using equality or the first nonnil comparator.
func SliceContains[T comparable](slice []T, target T, compare ...func(T, T) bool) bool {
	if len(compare) > 0 && compare[0] != nil {
		return slices.ContainsFunc(slice, func(item T) bool {
			return compare[0](item, target)
		})
	}
	return slices.Contains(slice, target)
}

// SliceIndexOf returns the first matching index, or -1 when absent.
func SliceIndexOf[T comparable](slice []T, target T, compare ...func(T, T) bool) int {
	if len(compare) > 0 && compare[0] != nil {
		return slices.IndexFunc(slice, func(item T) bool {
			return compare[0](item, target)
		})
	}
	return slices.Index(slice, target)
}

// SliceContainsBy reports whether any element matches the predicate.
func SliceContainsBy[T any](slice []T, compare func(ele T) bool) bool {
	return slices.ContainsFunc(slice, compare)
}

// SliceIndexBy returns the first predicate match, or -1 when absent.
func SliceIndexBy[T any](slice []T, compare func(ele T) bool) int {
	return slices.IndexFunc(slice, compare)
}

// SliceFind returns the first predicate match and whether one exists.
func SliceFind[T any](slice []T, filter func(item T) bool) (T, bool) {
	index := slices.IndexFunc(slice, filter)
	if index >= 0 {
		return slice[index], true
	}
	var zero T
	return zero, false
}

// SliceFindLast returns the last predicate match and whether one exists.
func SliceFindLast[T any](slice []T, filter func(item T) bool) (T, bool) {
	var t T
	if len(slice) == 0 {
		return t, false
	}
	var exist bool
	for i := len(slice) - 1; i >= 0; i-- {
		flag := filter(slice[i])
		if flag {
			t = slice[i]
			exist = true
			break
		}
	}
	return t, exist
}

// SliceFilter returns matching elements in input order.
func SliceFilter[T any](slice []T, filter func(item T) bool) []T {
	if len(slice) == 0 {
		return nil
	}
	result := make([]T, 0, len(slice))
	for i := range slice {
		flag := filter(slice[i])
		if flag {
			result = append(result, slice[i])
		}
	}
	return result
}

// SliceIntersection returns shared values, optionally using a custom comparator.
func SliceIntersection[T comparable](slicePart1, slicePart2 []T, compare ...func(part1, part2 T) bool) []T {
	var compareFn func(T, T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)
	cache := make(map[T]struct{})

	if compareFn == nil {
		// Use comparable-value membership when no comparator is supplied.
		setA := make(map[T]struct{})
		for i := range slicePart1 {
			setA[slicePart1[i]] = struct{}{}
		}
		for i := range slicePart2 {
			if _, found := setA[slicePart2[i]]; found {
				if _, added := cache[slicePart2[i]]; !added {
					cache[slicePart2[i]] = struct{}{}
					result = append(result, slicePart2[i])
				}
			}
		}
	} else {
		// Use the supplied comparator.
		for i := range slicePart1 {
			for j := range slicePart2 {
				if compareFn(slicePart1[i], slicePart2[j]) {
					if _, found := cache[slicePart1[i]]; !found {
						cache[slicePart1[i]] = struct{}{}
						result = append(result, slicePart1[i])
					}
				}
			}
		}
	}
	return result
}

// SliceUnion returns distinct values from both inputs using the selected equality rule.
func SliceUnion[T comparable](slicePart1, slicePart2 []T, compare ...func(part1 T, part2 T) bool) []T {
	var compareFn func(T, T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)

	// Preserve the first input order.
	for i := range slicePart1 {
		if !SliceContains(result, slicePart1[i], compareFn) {
			result = append(result, slicePart1[i])
		}
	}

	// Append unseen values from the second input.
	for i := range slicePart2 {
		if !SliceContains(result, slicePart2[i], compareFn) {
			result = append(result, slicePart2[i])
		}
	}
	return result
}

// SliceComplement returns values from the first input that are absent from the second.
func SliceComplement[T comparable](slicePart1, slicePart2 []T, compare ...func(part1, part2 T) bool) []T {
	var compareFn func(T, T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)
	cache := make(map[T]struct{})
	if compareFn == nil {
		for i := range slicePart2 {
			cache[slicePart2[i]] = struct{}{}
		}
		for i := range slicePart1 {
			if _, found := cache[slicePart1[i]]; !found {
				result = append(result, slicePart1[i])
			}
		}
	} else {
		for i := range slicePart1 {
			if !SliceContains(slicePart2, slicePart1[i], compareFn) {
				result = append(result, slicePart1[i])
			}
		}
	}
	return result
}

// SliceDiff returns values added in partNew and removed from partOld.
func SliceDiff[T comparable](partOld, partNew []T, compare ...func(T, T) bool) (added, removed []T) {
	if len(compare) > 0 {
		eq := compare[0]
		for _, va := range partOld {
			found := false
			for _, vb := range partNew {
				if eq(va, vb) {
					found = true
					break
				}
			}
			if !found {
				removed = append(removed, va)
			}
		}
		for _, vb := range partNew {
			found := false
			for _, va := range partOld {
				if eq(vb, va) {
					found = true
					break
				}
			}
			if !found {
				added = append(added, vb)
			}
		}
		return
	}
	amap := make(map[T]struct{}, len(partOld))
	bmap := make(map[T]struct{}, len(partNew))
	for _, v := range partOld {
		amap[v] = struct{}{}
	}
	for _, v := range partNew {
		bmap[v] = struct{}{}
	}
	for _, v := range partOld {
		if _, ok := bmap[v]; !ok {
			removed = append(removed, v)
		}
	}
	// Collect values present only in the new input.
	for _, v := range partNew {
		if _, ok := amap[v]; !ok {
			added = append(added, v)
		}
	}
	return
}

// SliceAnyDiff returns additions and removals using one-to-one predicate matches.
func SliceAnyDiff[T any](partOld, partNew []T, compare func(T, T) bool) (added, removed []T) {
	usedA := make([]bool, len(partOld))
	usedB := make([]bool, len(partNew))
	for i, va := range partOld {
		for j, vb := range partNew {
			if !usedB[j] && compare(va, vb) {
				usedA[i] = true
				usedB[j] = true
				break
			}
		}
	}
	for i, va := range partOld {
		if !usedA[i] {
			removed = append(removed, va)
		}
	}
	for j, vb := range partNew {
		if !usedB[j] {
			added = append(added, vb)
		}
	}
	return
}

// SliceIsSubset reports whether every value in the first slice belongs to the second.
func SliceIsSubset[T comparable](slicePart, sliceAll []T) bool {
	setMap := make(map[T]struct{}, len(sliceAll))
	for _, v := range sliceAll {
		setMap[v] = struct{}{}
	}

	for _, item := range slicePart {
		if _, ok := setMap[item]; !ok {
			return false
		}
	}
	return true
}

// SliceFilterToMap filters and maps elements into a map; later duplicate keys overwrite values.
func SliceFilterToMap[T any, K comparable, V any](slice []T, filter func(T) (K, V, bool)) map[K]V {
	if len(slice) == 0 {
		return nil
	}
	result := make(map[K]V, len(slice))
	for i := range slice {
		key, value, ok := filter(slice[i])
		if ok {
			result[key] = value
		}
	}
	return result
}

// SliceCollect maps each element to a new slice in input order.
func SliceCollect[T, R any](input []T, mapFn func(T) R) []R {
	if len(input) == 0 {
		return nil
	}
	output := make([]R, len(input))
	for i := range input {
		output[i] = mapFn(input[i])
	}
	return output
}

// SliceFilterCollect filters and maps elements into a new slice in input order.
func SliceFilterCollect[T, R any](input []T, mapFn func(T) (R, bool)) []R {
	if len(input) == 0 {
		return nil
	}
	output := make([]R, 0, len(input))
	for i := range input {
		v, f := mapFn(input[i])
		if f {
			output = append(output, v)
		}
	}
	return output
}

// SliceFlat flattens a two-dimensional slice, returning a nonnil empty slice for no elements.
func SliceFlat[T any](input [][]T) []T {
	if len(input) == 0 {
		return []T{}
	}
	output := slices.Concat(input...)
	if output == nil {
		return []T{}
	}
	return output
}

// SliceForEachAll visits every element in input order.
func SliceForEachAll[T any](slice []T, fn func(T)) {
	if len(slice) == 0 {
		return
	}
	for i := range slice {
		fn(slice[i])
	}
}

// SliceForEach visits elements until the callback returns false.
func SliceForEach[T any](slice []T, foreach func(T) bool) {
	if len(slice) == 0 {
		return
	}
	for i := range slice {
		if !foreach(slice[i]) {
			return
		}
	}
}

// SliceDistinct removes equal duplicates while keeping the first occurrence.
func SliceDistinct[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, ele := range slice {
		if _, ok := seen[ele]; ok {
			continue
		}
		seen[ele] = struct{}{}
		result = append(result, ele)
	}
	return result
}

// SliceDistinctBy removes duplicate keys while keeping the first matching element.
func SliceDistinctBy[T any, K comparable](slice []T, keyBuild func(ele T) K) []T {
	if len(slice) == 0 {
		return nil
	}
	seen := make(map[K]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, ele := range slice {
		key := keyBuild(ele)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, ele)
	}
	return result
}

// SliceSort sorts the input in place by integer weight, ascending unless desc is true.
func SliceSort[T any](slice []T, less func(e T) int, desc ...bool) {
	if len(slice) == 0 {
		return
	}
	descending := len(desc) > 0 && desc[0]
	slices.SortFunc(slice, func(a, b T) int {
		result := cmp.Compare(less(a), less(b))
		if descending {
			return -result
		}
		return result
	})
}

// SliceGroupBy groups elements into slices by key.
func SliceGroupBy[T any, K comparable](slice []T, groupFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for i := range slice {
		k := groupFn(slice[i])
		result[k] = append(result[k], slice[i])
	}
	return result
}

// SliceGroupBySingle groups elements by key, keeping the last element for each key.
func SliceGroupBySingle[T any, K comparable](slice []T, groupFn func(T) K) map[K]T {
	result := make(map[K]T)
	for i := range slice {
		k := groupFn(slice[i])
		result[k] = slice[i]
	}
	return result
}

// SliceGroupByValue groups projected values into slices by key.
func SliceGroupByValue[T, V any, K comparable](slice []T, groupFn func(T) (K, V)) map[K][]V {
	result := make(map[K][]V)
	for i := range slice {
		k, v := groupFn(slice[i])
		result[k] = append(result[k], v)
	}
	return result
}

// SliceGroupBySingleValue groups projected values by key, keeping the last value.
func SliceGroupBySingleValue[T, V any, K comparable](slice []T, groupFn func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for i := range slice {
		k, v := groupFn(slice[i])
		result[k] = v
	}
	return result
}

// SliceSplitChunk returns views of at most size elements; nonpositive size returns nil.
func SliceSplitChunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	return slices.Collect(slices.Chunk(slice, size))
}

// SliceRemove removes an indexed element, potentially modifying the backing array; invalid indexes return an error.
func SliceRemove[T any](index int, slice []T) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return slice, toolkitError.ErrSliceIndexOutOfRange
	}
	return slices.Delete(slice, index, index+1), nil
}

// SliceRemoveSafe removes an indexed element from a shallow copy of the slice.
func SliceRemoveSafe[T any](index int, slice []T) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return nil, toolkitError.ErrSliceIndexOutOfRange
	}
	return slices.Delete(slices.Clone(slice), index, index+1), nil
}

// SliceInsert inserts an element at a checked index, potentially modifying the backing array.
func SliceInsert[T any](index int, new T, slice []T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return slice, toolkitError.ErrSliceIndexOutOfRange
	}
	return slices.Insert(slice, index, new), nil
}

// SliceInsertSafe inserts an element into a shallow copy at a checked index.
func SliceInsertSafe[T any](index int, new T, slice []T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return nil, toolkitError.ErrSliceIndexOutOfRange
	}
	return slices.Insert(slices.Clone(slice), index, new), nil
}

// SliceInserts inserts multiple elements at a checked index, potentially modifying the backing array.
func SliceInserts[T any](index int, new []T, slice []T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return slice, toolkitError.ErrSliceIndexOutOfRange
	}
	return slices.Insert(slice, index, new...), nil
}

// SliceInsertsSafe inserts multiple elements into a shallow copy at a checked index.
func SliceInsertsSafe[T any](index int, new []T, slice []T) ([]T, error) {
	if index < 0 || index > len(slice) {
		return nil, toolkitError.ErrSliceIndexOutOfRange
	}
	cloned := slices.Clone(slice)
	if cloned == nil {
		cloned = []T{}
	}
	return slices.Insert(cloned, index, new...), nil
}
