package routine_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/acexy/golang-toolkit/sys/routine"
)

func TestThreadLocalWithProfilerLabels(t *testing.T) {
	id := routine.GetGoroutineId()
	if id == 0 {
		t.Fatal("current goroutine ID is zero")
	}
	childID := make(chan uint64, 1)
	go func() { childID <- routine.GetGoroutineId() }()
	if got := <-childID; got == 0 || got == id {
		t.Fatalf("child goroutine ID = %d, parent = %d", got, id)
	}

	pprof.Do(context.Background(), pprof.Labels("toolkit", "thread-local"), func(context.Context) {
		local := routine.NewThreadLocal[string]()
		defer local.Delete()
		local.Set("value")
		stack := make([]byte, 4096)
		n := runtime.Stack(stack, false)
		if !bytes.HasPrefix(stack[:n], []byte("goroutine ")) {
			t.Fatalf("unexpected current goroutine stack: %q", stack[:n])
		}
		var profile bytes.Buffer
		if err := pprof.Lookup("goroutine").WriteTo(&profile, 1); err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(profile.Bytes(), []byte(`"toolkit":"thread-local"`)) {
			t.Fatal("goroutine profile omitted the label after using thread-local storage")
		}
		if got := local.Get(); got != "value" {
			t.Fatalf("thread-local value = %q after collecting diagnostics", got)
		}
	})
}

func TestThreadLocalPublicMethods(t *testing.T) {
	// Verify that private embedded storage still promotes public methods externally.
	type localValue interface {
		Get() string
		Set(string)
		SetLazy(routine.LazySupplier[string])
		Delete()
	}
	for _, tc := range []struct {
		name string
		new  func() localValue
	}{
		{"ordinary", func() localValue {
			return routine.NewThreadLocal(func() string { return "default" })
		}},
		{"inheritable", func() localValue {
			return routine.NewInheritableThreadLocal(func() string { return "default" })
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			local := tc.new()
			defer local.Delete()
			local.Set("direct")
			if got := local.Get(); got != "direct" {
				t.Fatalf("Set 后 = %q", got)
			}
			local.SetLazy(func() string { return "lazy" })
			if got := local.Get(); got != "lazy" {
				t.Fatalf("SetLazy 后 = %q", got)
			}
			local.Delete()
			if got := local.Get(); got != "default" {
				t.Fatalf("Delete 后 = %q", got)
			}
		})
	}
}

func TestTraceIDPublicAPI(t *testing.T) {
	trace := routine.NewTraceIdThreadLocal(nil)
	defer trace.Delete()
	id := trace.Get()
	decoded, err := hex.DecodeString(id)
	if err != nil || len(id) != 32 || len(decoded) != 16 {
		t.Fatalf("TraceId 应为 32 位十六进制 UUID: %q, err=%v", id, err)
	}
	if decoded[6]>>4 != 4 || decoded[8]>>6 != 2 {
		t.Fatalf("TraceId 应保持 UUID v4 格式: %q", id)
	}
	for _, run := range []func(func()){trace.Run, trace.RunAll} {
		result := make(chan string, 1)
		run(func() { result <- trace.Get() })
		select {
		case got := <-result:
			if got != id {
				t.Fatalf("继承的 TraceId = %q, want %q", got, id)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("等待 TraceId 继承测试超时")
		}
	}
}
