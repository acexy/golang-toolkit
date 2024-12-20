package coll

// MapFirst 从map中抽取第一个元素
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

// MapValueToSlice 将map的value转换为slice
func MapValueToSlice[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}
	result := make([]V, 0)
	for i := range m {
		result = append(result, m[i])
	}
	return result
}

// MapCollect 将map转换成新的map
func MapCollect[K, RK comparable, V, RV any](m map[K]V, mapFn func(K, V) (RK, RV)) map[RK]RV {
	if len(m) == 0 {
		return nil
	}
	result := make(map[RK]RV)
	for k, v := range m {
		key, value := mapFn(k, v)
		result[key] = value
	}
	return result
}

// MapFilterCollect 将map过滤转换成新的map
func MapFilterCollect[K, RK comparable, V, RV any](m map[K]V, mapFn func(K, V) (RK, RV, bool)) map[RK]RV {
	if len(m) == 0 {
		return nil
	}
	result := make(map[RK]RV)
	for k, v := range m {
		key, value, ok := mapFn(k, v)
		if ok {
			result[key] = value
		}
	}
	return result
}

// MapFilterToSlice 将map过滤转换成切片
func MapFilterToSlice[K comparable, V, R any](m map[K]V, mapFn func(K, V) (R, bool)) []R {
	if len(m) == 0 {
		return nil
	}
	result := make([]R, 0)
	for k, v := range m {
		n, f := mapFn(k, v)
		if f {
			result = append(result, n)
		}
	}
	return result
}
