package slice

// Contains 检查指定的原属是否存在切片中
func Contains[T comparable](slice []T, elem T) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}
