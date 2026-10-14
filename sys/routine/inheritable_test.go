package routine

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func waitRoutine[T any](t *testing.T, result <-chan T) T {
	t.Helper()
	select {
	case value := <-result:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("等待子协程超时")
		var zero T
		return zero
	}
}

func runAndWait(t *testing.T, run func(func()), fn func()) {
	t.Helper()
	done := make(chan struct{})
	run(func() {
		defer close(done)
		fn()
	})
	waitRoutine(t, done)
}

func TestInheritableThreadLocalPropagationScope(t *testing.T) {
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("RunAll=%t", all), func(t *testing.T) {
			local := NewInheritableThreadLocal[string]()
			other := NewInheritableThreadLocal(func() int { return -1 })
			ordinary := NewThreadLocal(func() string { return "ordinary-default" })
			trace := NewTraceIdThreadLocal(func() string { return "trace-default" })
			local.Set("request-001")
			other.Set(1001)
			ordinary.Set("parent")
			trace.Set("parent-trace")
			defer local.Delete()
			defer other.Delete()
			defer ordinary.Delete()
			defer trace.Delete()

			run := local.Run
			if all {
				run = local.RunAll
			}
			runAndWait(t, run, func() {
				if got := local.Get(); got != "request-001" {
					t.Errorf("当前实例 = %q", got)
				}
				wantOther := -1
				if all {
					wantOther = 1001
				}
				if got := other.Get(); got != wantOther {
					t.Errorf("其他可继承实例 = %d, want %d", got, wantOther)
				}
				if got := ordinary.Get(); got != "ordinary-default" {
					t.Errorf("普通实例不应继承: %q", got)
				}
				wantTrace := "trace-default"
				if all {
					wantTrace = "parent-trace"
				}
				if got := trace.Get(); got != wantTrace {
					t.Errorf("TraceId = %q, want %q", got, wantTrace)
				}
			})
		})
	}
}

func TestInheritableThreadLocalNativeGoDoesNotInherit(t *testing.T) {
	local := NewInheritableThreadLocal(func() string { return "default" })
	local.Set("parent")
	defer local.Delete()
	result := make(chan string, 1)
	go func() { result <- local.Get() }()
	if got := waitRoutine(t, result); got != "default" {
		t.Fatalf("原生 go 不应继承: %q", got)
	}
}

func TestInheritableThreadLocalMixedNestedScope(t *testing.T) {
	local := NewInheritableThreadLocal[string]()
	other := NewInheritableThreadLocal(func() int { return -1 })
	local.Set("trace-001")
	other.Set(1001)
	defer local.Delete()
	defer other.Delete()
	runAndWait(t, local.RunAll, func() {
		if got := other.Get(); got != 1001 {
			t.Errorf("第一层 RunAll = %d", got)
		}
		runAndWait(t, local.Run, func() {
			if got := local.Get(); got != "trace-001" {
				t.Errorf("第二层 Run 当前实例 = %q", got)
			}
			if got := other.Get(); got != -1 {
				t.Errorf("第二层 Run 不应携带其他实例: %d", got)
			}
			runAndWait(t, local.RunAll, func() {
				if got := other.Get(); got != -1 {
					t.Errorf("RunAll 不应找回祖先上下文: %d", got)
				}
			})
			result := make(chan string, 1)
			go func() { result <- local.Get() }()
			if got := waitRoutine(t, result); got != "" {
				t.Errorf("嵌套中的原生 go 不应继承: %q", got)
			}
		})
		if got := other.Get(); got != 1001 {
			t.Errorf("下层初始化影响了上层: %d", got)
		}
	})
}

func TestInheritableThreadLocalCrossRunMatrix(t *testing.T) {
	entries := []string{"native", "A.Run", "B.Run", "Trace.Run", "A.RunAll", "B.RunAll", "Trace.RunAll"}
	for _, first := range entries {
		for _, second := range entries {
			t.Run(first+"/"+second, func(t *testing.T) {
				a := NewInheritableThreadLocal(func() string { return "A-default" })
				b := NewInheritableThreadLocal(func() string { return "B-default" })
				ordinary := NewThreadLocal(func() string { return "ordinary-default" })
				trace := NewTraceIdThreadLocal(func() string { return "trace-default" })
				a.Set("A-parent")
				b.Set("B-parent")
				ordinary.Set("ordinary-parent")
				trace.Set("trace-parent")
				defer a.Delete()
				defer b.Delete()
				defer ordinary.Delete()
				defer trace.Delete()

				launch := func(entry string) func(func()) {
					switch entry {
					case "A.Run":
						return a.Run
					case "B.Run":
						return b.Run
					case "Trace.Run":
						return trace.Run
					case "A.RunAll":
						return a.RunAll
					case "B.RunAll":
						return b.RunAll
					case "Trace.RunAll":
						return trace.RunAll
					default:
						return func(fn func()) { go fn() }
					}
				}
				inherited := func(entry, parentA, parentB, parentTrace string) (string, string, string) {
					switch entry {
					case "A.Run":
						return parentA, "B-default", "trace-default"
					case "B.Run":
						return "A-default", parentB, "trace-default"
					case "Trace.Run":
						return "A-default", "B-default", parentTrace
					case "A.RunAll", "B.RunAll", "Trace.RunAll":
						return parentA, parentB, parentTrace
					default:
						return "A-default", "B-default", "trace-default"
					}
				}
				check := func(wantA, wantB, wantTrace string) {
					t.Helper()
					if got := a.Get(); got != wantA {
						t.Errorf("A = %q, want %q", got, wantA)
					}
					if got := b.Get(); got != wantB {
						t.Errorf("B = %q, want %q", got, wantB)
					}
					if got := ordinary.Get(); got != "ordinary-default" {
						t.Errorf("普通变量不应传递: %q", got)
					}
					if got := trace.Get(); got != wantTrace {
						t.Errorf("TraceId = %q, want %q", got, wantTrace)
					}
				}

				wantA, wantB, wantTrace := inherited(first, "A-parent", "B-parent", "trace-parent")
				runAndWait(t, launch(first), func() {
					check(wantA, wantB, wantTrace)
					ordinary.Set("ordinary-child")
					trace.Set("trace-child")
					nestedA, nestedB, nestedTrace := inherited(second, wantA, wantB, "trace-child")
					runAndWait(t, launch(second), func() {
						check(nestedA, nestedB, nestedTrace)
						a.Set("A-grandchild")
						b.Delete()
						ordinary.Set("ordinary-grandchild")
						trace.Delete()
					})
					if a.Get() != wantA || b.Get() != wantB ||
						ordinary.Get() != "ordinary-child" || trace.Get() != "trace-child" {
						t.Error("孙协程修改影响了子协程")
					}
				})
				if a.Get() != "A-parent" || b.Get() != "B-parent" ||
					ordinary.Get() != "ordinary-parent" || trace.Get() != "trace-parent" {
					t.Fatal("下层修改影响了父协程")
				}
			})
		}
	}
}

func TestOrdinaryThreadLocalLazyIsNotInherited(t *testing.T) {
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("RunAll=%t", all), func(t *testing.T) {
			trigger := NewInheritableThreadLocal[string]()
			ordinary := NewThreadLocal(func() string { return "default" })
			var calls atomic.Int32
			ordinary.SetLazy(func() string {
				calls.Add(1)
				return "parent-lazy"
			})
			defer trigger.Delete()
			defer ordinary.Delete()
			run := trigger.Run
			if all {
				run = trigger.RunAll
			}
			runAndWait(t, run, func() {
				if got := ordinary.Get(); got != "default" {
					t.Errorf("普通实例的 SetLazy 不应继承: %q", got)
				}
			})
			if got := calls.Load(); got != 0 {
				t.Fatalf("子协程不应执行父协程的延迟初始化器: %d", got)
			}
			if got := ordinary.Get(); got != "parent-lazy" || calls.Load() != 1 {
				t.Fatalf("父协程延迟初始化状态被改变: %q", got)
			}
		})
	}
}

func TestInheritableThreadLocalNestedAndIsolated(t *testing.T) {
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("RunAll=%t", all), func(t *testing.T) {
			local := NewInheritableThreadLocal(func() string { return "default" })
			local.Set("parent")
			defer local.Delete()
			run := local.Run
			if all {
				run = local.RunAll
			}
			runAndWait(t, run, func() {
				if got := local.Get(); got != "parent" {
					t.Errorf("子协程 = %q", got)
				}
				local.Set("child")
				runAndWait(t, run, func() {
					if got := local.Get(); got != "child" {
						t.Errorf("孙协程 = %q", got)
					}
					local.SetLazy(func() string { return "grandchild" })
					if got := local.Get(); got != "grandchild" {
						t.Errorf("孙协程 SetLazy = %q", got)
					}
					local.Delete()
					if got := local.Get(); got != "default" {
						t.Errorf("Delete 后 = %q", got)
					}
				})
				if got := local.Get(); got != "child" {
					t.Errorf("孙协程修改影响了子协程: %q", got)
				}
			})
			if got := local.Get(); got != "parent" {
				t.Fatalf("子协程修改影响了父协程: %q", got)
			}
			runAndWait(t, run, func() {
				if got := local.Get(); got != "parent" {
					t.Errorf("兄弟协程 = %q", got)
				}
			})
		})
	}
}

func TestInheritableThreadLocalSnapshotAndAsync(t *testing.T) {
	for _, all := range []bool{false, true} {
		for _, lazy := range []bool{false, true} {
			t.Run(fmt.Sprintf("RunAll=%t/Lazy=%t", all, lazy), func(t *testing.T) {
				local := NewInheritableThreadLocal[string]()
				other := NewInheritableThreadLocal[int]()
				defer local.Delete()
				defer other.Delete()
				if lazy {
					local.SetLazy(func() string { return "captured" })
				} else {
					local.Set("captured")
				}
				other.Set(1001)
				run := local.Run
				if all {
					run = local.RunAll
				}
				gate := make(chan struct{})
				done := make(chan struct{})
				run(func() {
					defer close(done)
					<-gate
					if got := local.Get(); got != "captured" {
						t.Errorf("调用时快照 = %q", got)
					}
					want := 0
					if all {
						want = 1001
					}
					if got := other.Get(); got != want {
						t.Errorf("其他实例快照 = %d, want %d", got, want)
					}
				})
				local.Set("changed")
				other.Delete()
				close(gate)
				waitRoutine(t, done)
			})
		}
	}
}

func TestInheritableThreadLocalLazyIndependentInitialization(t *testing.T) {
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("RunAll=%t", all), func(t *testing.T) {
			local := NewInheritableThreadLocal[int32]()
			var calls atomic.Int32
			local.SetLazy(func() int32 { return calls.Add(1) })
			defer local.Delete()
			run := local.Run
			if all {
				run = local.RunAll
			}
			const children = 16
			gate := make(chan struct{})
			results := make(chan int32, children)
			for i := 0; i < children; i++ {
				run(func() {
					<-gate
					value := local.Get()
					if got := local.Get(); got != value {
						t.Errorf("同一协程重复 Get = %d, want %d", got, value)
					}
					results <- value
				})
			}
			if got := calls.Load(); got != 0 {
				t.Errorf("启动任务不应求值: %d", got)
			}
			close(gate)
			parent := local.Get()
			seen := map[int32]bool{parent: true}
			for i := 0; i < children; i++ {
				value := waitRoutine(t, results)
				if seen[value] {
					t.Errorf("父子或兄弟协程共享了初始化结果: %d", value)
				}
				seen[value] = true
			}
			if got := calls.Load(); got != children+1 {
				t.Fatalf("supplier 次数 = %d, want %d", got, children+1)
			}
		})
	}
}

func TestInheritableThreadLocalInitializedLazyIsInherited(t *testing.T) {
	for _, all := range []bool{false, true} {
		t.Run(fmt.Sprintf("RunAll=%t", all), func(t *testing.T) {
			local := NewInheritableThreadLocal[string]()
			var calls atomic.Int32
			local.SetLazy(func() string {
				calls.Add(1)
				return "cached"
			})
			defer local.Delete()
			local.Get()
			run := local.Run
			if all {
				run = local.RunAll
			}
			runAndWait(t, run, func() {
				if got := local.Get(); got != "cached" {
					t.Errorf("缓存值 = %q", got)
				}
			})
			if got := calls.Load(); got != 1 {
				t.Fatalf("已初始化的 supplier 不应再次执行: %d", got)
			}
		})
	}
}

func TestInheritableThreadLocalDefaultsAndOverrides(t *testing.T) {
	var calls atomic.Int32
	local := NewInheritableThreadLocal(func() int32 { return calls.Add(1) })
	defer local.Delete()
	if got := calls.Load(); got != 0 {
		t.Fatal("构造函数不应执行 supplier")
	}
	runAndWait(t, local.Run, func() {
		if got := local.Get(); got != 1 {
			t.Errorf("子协程默认初始化 = %d", got)
		}
	})
	if got := calls.Load(); got != 1 {
		t.Fatal("Run 不应在父协程提前执行 supplier")
	}
	if got := local.Get(); got != 2 {
		t.Fatalf("父协程独立初始化 = %d", got)
	}
	local.SetLazy(func() int32 { return 99 })
	local.Set(42)
	if got := local.Get(); got != 42 {
		t.Fatalf("Set 覆盖 SetLazy = %d", got)
	}
	local.SetLazy(func() int32 { return 99 })
	if got := local.Get(); got != 99 {
		t.Fatalf("SetLazy 覆盖 Set = %d", got)
	}
	local.Delete()
	if got := local.Get(); got != 3 {
		t.Fatalf("Delete 后重新初始化 = %d", got)
	}
	local.SetLazy(nil)
	if got := local.Get(); got != 0 {
		t.Fatalf("nil 延迟初始化器应返回零值: %d", got)
	}
	zero := NewInheritableThreadLocal[int](nil)
	defer zero.Delete()
	runAndWait(t, zero.Run, func() {
		if got := zero.Get(); got != 0 {
			t.Errorf("nil 默认初始化器 = %d", got)
		}
	})
}

func TestInheritableThreadLocalRunAllLazyOtherAndUnused(t *testing.T) {
	trigger := NewInheritableThreadLocal[string]()
	other := NewInheritableThreadLocal[int32]()
	var calls atomic.Int32
	other.SetLazy(func() int32 { return calls.Add(1) })
	defer other.Delete()
	unused := NewInheritableThreadLocal(func() string {
		t.Error("未使用的实例不应自动初始化")
		return "unused"
	})
	defer unused.Delete()
	runAndWait(t, trigger.RunAll, func() {
		if got := other.Get(); got != 1 {
			t.Errorf("其他实例的延迟初始化 = %d", got)
		}
	})
	if got := other.Get(); got != 2 {
		t.Fatalf("其他实例的包装状态被共享: %d", got)
	}
}

func TestInheritableTraceIDNested(t *testing.T) {
	var trace *InheritableThreadLocal[string] = NewTraceIdThreadLocal(nil)
	defer trace.Delete()
	id := trace.Get()
	if len(id) != 32 {
		t.Fatalf("TraceId 长度 = %d", len(id))
	}
	runAndWait(t, trace.Run, func() {
		if got := trace.Get(); got != id {
			t.Errorf("子协程 TraceId = %q, want %q", got, id)
		}
		runAndWait(t, trace.RunAll, func() {
			if got := trace.Get(); got != id {
				t.Errorf("孙协程 TraceId = %q, want %q", got, id)
			}
		})
	})
}

func TestInheritableThreadLocalReferenceValuesAreShallowCopied(t *testing.T) {
	local := NewInheritableThreadLocal[map[string]int]()
	local.Set(map[string]int{"value": 1})
	defer local.Delete()
	runAndWait(t, local.Run, func() { local.Get()["value"] = 2 })
	if got := local.Get()["value"]; got != 2 {
		t.Fatalf("引用数据应保持浅拷贝语义: %d", got)
	}
}

func TestInheritableThreadLocalPanicDoesNotEscape(t *testing.T) {
	local := NewInheritableThreadLocal[string]()
	local.Set("parent")
	defer local.Delete()
	for _, run := range []func(func()){local.Run, local.RunAll} {
		runAndWait(t, run, func() { panic("test child panic") })
		// A failed task leaves the parent context available for later inheritance.
		runAndWait(t, run, func() {
			if got := local.Get(); got != "parent" {
				t.Errorf("panic 后上下文 = %q", got)
			}
		})
	}
}
