package coll

// MapFirst 从map中抽取第一个元素，由于map的无序性，所以这里返回的结果并不确定
func MapFirst[K comparable, V any](m map[K]V) (K, V) {
	var key K
	var value V
	if len(m) == 0 {

		return key, value
	}
	for k, v := range m {
		key = k
		value = v
		return key, value
	}
	return key, value
}

// MapKeyToSlice 将map的key转换为slice
func MapKeyToSlice[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}
	result := make([]K, 0)
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

// MapToSlice 将map转换为slice
func MapToSlice[K comparable, V, R any](m map[K]V, mapFn func(K, V) R) []R {
	if len(m) == 0 {
		return nil
	}
	result := make([]R, 0)
	for k, v := range m {
		result = append(result, mapFn(k, v))
	}
	return result
}
