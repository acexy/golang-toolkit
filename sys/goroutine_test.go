package sys

import (
	"fmt"
	"github.com/acexy/golang-toolkit/math/random"
	"testing"
)

func TestLocal(t *testing.T) {
	l := NewThreadLocal[string](func() string {
		return random.UUID()
	})
	// 如果在获取线程本地变量之前，没有设置过，则返回默认值
	//l.Set("123")
	fmt.Println(l.Get())
	f1(l)
}
func f1(local *Local[string]) {
	fmt.Println(local.Get())
}

func TestEnableTraceIdLocal(t *testing.T) {
	EnableLocalTraceId(nil)
	fmt.Println(GetLocalTraceId())
	f2()
	SetLocalTraceId("123")
	f2()
}

func f2() {
	fmt.Println(GetLocalTraceId())
}
