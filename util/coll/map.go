package coll

import "math/rand"

// MapAny 从map中抽取任意一个元素
func MapAny[K comparable, V any](m map[K]V) (K, V) {
	var key K
	var value V
	if len(m) == 0 {
		return key, value
	}
	for k, v := range m {
		return k, v
	}
	return key, value
}

// MapKeys 将map的key转换为slice
func MapKeys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// MapValues 将map的value转换为slice
func MapValues[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// MapCollect 将map转换成新的map
func MapCollect[K, RK comparable, V, RV any](m map[K]V, mapFn func(K, V) (RK, RV)) map[RK]RV {
	if len(m) == 0 {
		return nil
	}
	result := make(map[RK]RV, len(m))
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
	result := make(map[RK]RV, len(m))
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
	result := make([]R, 0, len(m))
	for k, v := range m {
		n, f := mapFn(k, v)
		if f {
			result = append(result, n)
		}
	}
	return result
}

// MapForEach 遍历map return false时停止迭代
func MapForEach[K comparable, V any](m map[K]V, fn func(k K, v V) bool) {
	for k, v := range m {
		if !fn(k, v) {
			return
		}
	}
}

// MapForEachAll 遍历map
func MapForEachAll[K comparable, V any](m map[K]V, fn func(k K, v V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// MapMerge 将source追加到target中
func MapMerge[K comparable, V any](target map[K]V, source map[K]V) map[K]V {
	if len(source) == 0 {
		return target
	}
	if target == nil {
		target = make(map[K]V, len(source))
	}
	for k, v := range source {
		target[k] = v
	}
	return target
}

// MapRandom 随机获取一个map的元素
func MapRandom[K comparable, V any](m map[K]V) (K, V) {
	if len(m) == 0 {
		var zk K
		var zv V
		return zk, zv
	}
	r := rand.Intn(len(m))
	i := 0
	for k, v := range m {
		if i == r {
			return k, v
		}
		i++
	}
	var zk K
	var zv V
	return zk, zv
}
