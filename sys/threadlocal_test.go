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
	fmt.Println(l.Get())
	f1(l)
}
func f1(local *Local[string]) {
	fmt.Println(local.Get())
}

func TestEnableTraceIdLocal(t *testing.T) {
	EnableTraceIdLocal(nil)
	fmt.Println(GetTraceId())
	f2()
}

func f2() {
	fmt.Println(GetTraceId())
}
