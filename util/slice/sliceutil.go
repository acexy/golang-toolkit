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

// ContainsWithFn ContainsWhitFn 检查指定的元素是否存在切片中
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
