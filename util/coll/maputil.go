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
	}
	return key, value
}
