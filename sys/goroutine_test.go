package sys

import (
	"fmt"
	"testing"

	"github.com/acexy/golang-toolkit/math/random"
)

func TestLocal(t *testing.T) {
	l := NewThreadLocal[string](func() string {
		return random.UUID()
	})
	// 如果在获取线程本地变量之前，没有设置过，则返回默认值
	fmt.Println(l.Get())
	l.Set("123")
	fmt.Println(l.Get())
	f1(l)
}
func f1(local *ThreadLocal[string]) {
	fmt.Println(local.Get())
}

func f2(local *ThreadLocal[string]) {
	fmt.Println(local.Get())
}
