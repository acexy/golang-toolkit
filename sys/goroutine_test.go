package sys

import (
	"fmt"
	"testing"

	"github.com/acexy/golang-toolkit/math/random"
)

var l = NewThreadLocal[string](func() string {
	return random.UUID()
})

func TestLocal(t *testing.T) {
	// 如果在获取线程本地变量之前，没有设置过，则返回默认值
	fmt.Println(l.Get())
	l.Set("123")
	fmt.Println(l.Get())
	f1()
	f2()
}
func f1() {
	fmt.Println(l.Get())
}

func f2() {
	fmt.Println(l.Get())
}
