package routine

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/acexy/golang-toolkit/math/random"
)

var l = NewThreadLocal[string](func() string {
	return random.UUID()
})

func TestLocal(t *testing.T) {
	// An unset local uses its default supplier on first access.
	fmt.Println(l.Get())
	l.Set(func() string {
		fmt.Println("set")
		return "123"
	}())
	fmt.Println(l.Get())
	f1()
	f2()

}
func f1() {
	fmt.Println("f1", l.Get())
}

func f2() {
	go func() {
		fmt.Println("f2", "goroutine", l.Get())
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("f2", l.Get())
}

func TestThreadLocalSetLazy(t *testing.T) {
	local := NewThreadLocal[string]()
	var called atomic.Int32

	local.SetLazy(func() string {
		called.Add(1)
		return "lazy value"
	})

	if got := called.Load(); got != 0 {
		t.Fatalf("SetLazy 不应执行 supplier，实际执行次数: %d", got)
	}
	if got := local.Get(); got != "lazy value" {
		t.Fatalf("Get() = %q, want %q", got, "lazy value")
	}
	if got := local.Get(); got != "lazy value" {
		t.Fatalf("第二次 Get() = %q, want %q", got, "lazy value")
	}
	if got := called.Load(); got != 1 {
		t.Fatalf("supplier 应只执行一次，实际执行次数: %d", got)
	}
}

func TestThreadLocalSetLazyOverridesInitialSupplier(t *testing.T) {
	var defaultCalled atomic.Int32
	var lazyCalled atomic.Int32
	local := NewThreadLocal(func() string {
		defaultCalled.Add(1)
		return "default"
	})

	local.SetLazy(func() string {
		lazyCalled.Add(1)
		return "request"
	})

	if got := local.Get(); got != "request" {
		t.Fatalf("Get() = %q, want %q", got, "request")
	}
	if got := defaultCalled.Load(); got != 0 {
		t.Fatalf("默认 supplier 不应执行，实际执行次数: %d", got)
	}
	if got := lazyCalled.Load(); got != 1 {
		t.Fatalf("请求 supplier 应执行一次，实际执行次数: %d", got)
	}
}

func TestThreadLocalSetAndSetLazyOverrideEachOther(t *testing.T) {
	local := NewThreadLocal[string]()
	var called atomic.Int32

	local.SetLazy(func() string {
		called.Add(1)
		return "lazy value"
	})
	local.Set("direct value")
	if got := local.Get(); got != "direct value" {
		t.Fatalf("Set 覆盖 SetLazy 后 Get() = %q, want %q", got, "direct value")
	}
	if got := called.Load(); got != 0 {
		t.Fatalf("被 Set 覆盖的 supplier 不应执行，实际执行次数: %d", got)
	}

	local.SetLazy(func() string {
		called.Add(1)
		return "next lazy value"
	})
	if got := local.Get(); got != "next lazy value" {
		t.Fatalf("SetLazy 覆盖 Set 后 Get() = %q, want %q", got, "next lazy value")
	}
	if got := called.Load(); got != 1 {
		t.Fatalf("supplier 应执行一次，实际执行次数: %d", got)
	}
}
