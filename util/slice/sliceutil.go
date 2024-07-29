package slice

// Contains 检查指定的元素是否存在切片中
// compare 比较函数，如果为空 则直接使用值比较
func Contains[T comparable](slice []T, elem T, compare ...func(*T, *T) bool) bool {
	var compareFn func(*T, *T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	for _, item := range slice {
		if compareFn != nil {
			if compareFn(&item, &elem) {
				return true
			}
		}
		if item == elem {
			return true
		}
	}
	return false
}

// Filter 筛选切片 通过函数筛选出符合要求的元素
func Filter[T comparable](slice []T, filter func(item *T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		flag := filter(&item)
		if flag {
			result = append(result, item)
		}
	}
	return result
}

// Intersection 求两个切片的交集 两个集合中共同的元素所组成的集合
func Intersection[T comparable](sliceA, sliceB []T, compare ...func(*T, *T) bool) []T {
	var compareFn func(*T, *T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)
	cache := make(map[T]struct{})

	if compareFn == nil {
		// 如果没有自定义的比较函数，使用默认的比较方式
		setA := make(map[T]struct{})
		for _, v := range sliceA {
			setA[v] = struct{}{}
		}
		for _, v := range sliceB {
			if _, found := setA[v]; found {
				if _, added := cache[v]; !added {
					cache[v] = struct{}{}
					result = append(result, v)
				}
			}
		}
	} else {
		// 使用自定义的比较函数
		for _, v1 := range sliceA {
			for _, v2 := range sliceB {
				if compareFn(&v1, &v2) {
					if _, found := cache[v1]; !found {
						cache[v1] = struct{}{}
						result = append(result, v1)
					}
				}
			}
		}
	}

	return result
}

// Union 求两个切片的并集 两个集合中所有元素（不重复）所组成的集合。
func Union[T comparable](sliceA, sliceB []T, compare ...func(*T, *T) bool) []T {
	var compareFn func(*T, *T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)

	// 处理 sliceA 中的元素
	for _, v := range sliceA {
		if !Contains(result, v, compareFn) {
			result = append(result, v)
		}
	}

	// 处理 sliceB 中的元素
	for _, v := range sliceB {
		if !Contains(result, v, compareFn) {
			result = append(result, v)
		}
	}
	return result
}

// Complement 求两个切片的补集 全集中(sliceAll)不属于某个集合(slicePart)的元素所组成的集合
func Complement[T comparable](sliceAll, slicePart []T, compare ...func(*T, *T) bool) []T {
	var compareFn func(*T, *T) bool
	if len(compare) > 0 {
		compareFn = compare[0]
	}
	result := make([]T, 0)

	cache := make(map[T]struct{})
	if compareFn == nil {
		for _, v := range slicePart {
			cache[v] = struct{}{}
		}
		for _, v := range sliceAll {
			if _, found := cache[v]; !found {
				result = append(result, v)
			}
		}
	} else {
		for _, v := range sliceAll {
			if !Contains(slicePart, v, compareFn) {
				result = append(result, v)
			}
		}
	}
	return result
}

// ToMap 将切片按照指定的过滤处理形成map
func ToMap[T any, K comparable, V any](slice []T, filter func(*T) (*K, *V, bool)) map[K]V {
	if len(slice) == 0 {
		return nil
	}
	result := make(map[K]V, len(slice))
	for _, item := range slice {
		key, value, ok := filter(&item)
		if ok {
			result[*key] = *value
		}
	}
	return result
}
