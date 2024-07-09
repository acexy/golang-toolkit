package conversion

import "unsafe"

// ParseBytes 通过go 1.22新函数 将字符串转换为字节数组 不进行内存分配高性能
func ParseBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
