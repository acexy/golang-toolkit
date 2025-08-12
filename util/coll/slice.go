package coll

import "sort"

// SliceContains 检查指定的元素是否存在切片中
func SliceContains[T comparable](slice []T, target T, compare ...func(T, T) bool) bool {
	for i := range slice {
		if len(compare) > 0 && compare[0] != nil {
			if compare[0](slice[i], target) {
				return true
			}
		} else {
			if slice[i] == target {
				return true
			}
		}
	}
	return false
}

// SliceIndexOf 获取指定元素在切片中的索引，如果元素不存在则返回-1
func SliceIndexOf[T comparable](slice []T, target T, compare ...func(T, T) bool) int {
	for i := range slice {
		if len(compare) > 0 && compare[0] != nil {
			if compare[0](slice[i], target) {
				return i
			}
		} else {
			if slice[i] == target {
				return i
			}
		}
	}
	return -1
}

// SliceAnyContains 检查指定的元素是否存在切片中，元素可以是任意类型
func SliceAnyContains[T comparable](slice []T, compare func(ele T) bool) bool {
	for i := range slice {
		if compare(slice[i]) {
			return true
		}
	}
	return false
}

// SliceAnyIndexOf 获取指定元素在切片中的索引，如果元素不存在则返回-1
func SliceAnyIndexOf[T comparable](slice []T, compare func(ele T) bool) int {
	for i := range slice {
		if compare(slice[i]) {
			return i
		}
	}
	return -1
}

// SliceFilterFirstOne 筛选切片 通过函数筛选出符合要求的第一个元素
func SliceFilterFirstOne[T any](slice []T, filter func(item T) bool) (T, bool) {
	var t T
	var exist bool
	for i := range slice {
		flag := filter(slice[i])
		if flag {
			t = slice[i]
			exist = true
			break
		}
	}
	return t, exist
}

// SliceFilter 筛选切片 通过函数筛选出符合要求的元素
func SliceFilter[T any](slice []T, filter func(item T) bool) []T {
	result := make([]T, 0)
	for i := range slice {
		flag := filter(slice[i])
		if flag {
			result = append(result, slice[i])
		}
	}
	return result
}

// SliceIntersection 求两个切片的交集 两个集合中共同的元素所组成的集合
func SliceIntersection[T comparable](slicePart1, slicePart2 []T, compare ...func(part1, part2 T) bool) []T {
	var compareFn func(T, T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)
	cache := make(map[T]struct{})

	if compareFn == nil {
		// 如果没有自定义的比较函数，使用默认的比较方式
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
		// 使用自定义的比较函数
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

// SliceUnion 求两个切片的并集 两个集合中所有元素（不重复）所组成的集合。
func SliceUnion[T comparable](slicePart1, slicePart2 []T, compare ...func(part1 T, part2 T) bool) []T {
	var compareFn func(T, T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)

	// 处理 slicePart1 中的元素
	for i := range slicePart1 {
		if !SliceContains(result, slicePart1[i], compareFn) {
			result = append(result, slicePart1[i])
		}
	}

	// 处理 slicePart2 中的元素
	for i := range slicePart2 {
		if !SliceContains(result, slicePart2[i], compareFn) {
			result = append(result, slicePart2[i])
		}
	}
	return result
}

// SliceComplement 求两个切片的补集 全集中(slicePart1)不属于某个集合(slicePart2)的元素所组成的集合
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

// SliceIsSubset 判断切片A是否是切片B的子集
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

// SliceFilterToMap 将切片按照指定的过滤处理形成map
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

// SliceCollect 将切片按照指定的采集映射方法处理为一个新的切片
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

// SliceFilterCollect 将切片按照指定的过滤条件采集映射方法处理为一个新的切片
func SliceFilterCollect[T, R any](input []T, mapFn func(T) (R, bool)) []R {
	if len(input) == 0 {
		return nil
	}
	output := make([]R, 0)
	for i := range input {
		v, f := mapFn(input[i])
		if f {
			output = append(output, v)
		}
	}
	return output
}

// SliceForeachAll 遍历所有切片并执行指定的函数
func SliceForeachAll[T any](slice []T, fn func(T)) {
	if len(slice) == 0 {
		return
	}
	for i := range slice {
		fn(slice[i])
	}
}

// SliceForeach 遍历切片并执行指定的函数，如果返回false则停止遍历
func SliceForeach[T any](slice []T, foreach func(T) bool) {
	if len(slice) == 0 {
		return
	}
	for i := range slice {
		if !foreach(slice[i]) {
			return
		}
	}
}

// SliceDistinct 去除切片的重复元素
func SliceDistinct[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}
	mapValue := SliceFilterToMap(slice, func(ele T) (T, any, bool) {
		return ele, struct{}{}, true
	})
	return MapKeyToSlice(mapValue)
}

// SliceDistinctAny 按照指定的切片去重值去重元素
func SliceDistinctAny[T any, K comparable](slice []T, keyBuild func(ele T) K) []T {
	if len(slice) == 0 {
		return nil
	}
	mapValue := SliceFilterToMap(slice, func(ele T) (K, T, bool) {
		return keyBuild(ele), ele, true
	})
	return MapValueToSlice(mapValue)
}

// SliceSort 对切片进行排序
// less: 元素排序的权重值
// desc: 是否降序 默认为升序
func SliceSort[T any](slice []T, less func(e T) int, desc ...bool) {
	if len(slice) == 0 {
		return
	}
	sort.Slice(slice, func(i, j int) bool {
		asc := true
		if len(desc) > 0 {
			asc = !desc[0]
		}
		if asc {
			return less(slice[i]) < less(slice[j])
		} else {
			return less(slice[i]) > less(slice[j])
		}
	})
}

// SliceGroupBy 将切片按照指定的分组函数进行分组
func SliceGroupBy[T any, K comparable](slice []T, groupFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for i := range slice {
		k := groupFn(slice[i])
		result[k] = append(result[k], slice[i])
	}
	return result
}

// SliceGroupBySingle 将切片按照指定的分组函数进行分组 适用于分组的值为单个元素
func SliceGroupBySingle[T any, K comparable](slice []T, groupFn func(T) K) map[K]T {
	result := make(map[K]T)
	for i := range slice {
		k := groupFn(slice[i])
		result[k] = slice[i]
	}
	return result
}

// SliceAnyGroupBy 将切片按照指定的分组函数进行分组
func SliceAnyGroupBy[T, V any, K comparable](slice []T, groupFn func(T) (K, V)) map[K][]V {
	result := make(map[K][]V)
	for i := range slice {
		k, v := groupFn(slice[i])
		result[k] = append(result[k], v)
	}
	return result
}

// SliceAnyGroupBySingle 将切片按照指定的分组函数进行分组 适用于分组的值为单个元素
func SliceAnyGroupBySingle[T, V any, K comparable](slice []T, groupFn func(T) (K, V)) map[K]V {
	result := make(map[K]V)
	for i := range slice {
		k, v := groupFn(slice[i])
		result[k] = v
	}
	return result
}
