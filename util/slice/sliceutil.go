package slice

// Contains 检查指定的元素是否存在切片中
func Contains[T comparable](slice []T, elem T) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}

// ContainsWithFn ContainsWhitFn 通过函数检查指定的元素是否存在切片中
func ContainsWithFn[T comparable](slice []T, fn func(item *T) bool) bool {
	for _, item := range slice {
		flag := fn(&item)
		if flag {
			return true
		}
	}
	return false
}

// FilterWithFn 检查原始切片是否满足符合条件的元素
func FilterWithFn[T comparable](slice []T, fn func(item *T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		flag := fn(&item)
		if flag {
			result = append(result, item)
		}
	}
	return result
}

// Intersection 求两个切片的交集 两个集合中共同的元素所组成的集合
func Intersection[T comparable](sliceA, sliceB []T) []T {
	set := make(map[T]struct{})
	result := make([]T, 0)

	for _, value := range sliceA {
		set[value] = struct{}{}
	}

	for _, value := range sliceB {
		if _, found := set[value]; found {
			result = append(result, value)
			delete(set, value) // 防止结果中出现重复元素
		}
	}
	return result
}

// Union 求两个切片的并集 两个集合中所有元素（不重复）所组成的集合。
func Union[T comparable](sliceA, sliceB []T) []T {
	set := make(map[T]struct{})
	result := make([]T, 0)

	for _, value := range sliceA {
		if _, found := set[value]; !found {
			set[value] = struct{}{}
			result = append(result, value)
		}
	}

	for _, value := range sliceB {
		if _, found := set[value]; !found {
			set[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

// Complement 求两个切片的补集 全集中(sliceAll)不属于某个集合(slicePart)的元素所组成的集合
func Complement[T comparable](sliceAll, slicePart []T) []T {
	set := make(map[T]struct{})
	result := make([]T, 0)

	for _, value := range slicePart {
		set[value] = struct{}{}
	}

	for _, value := range sliceAll {
		if _, found := set[value]; !found {
			result = append(result, value)
		}
	}
	return result
}
