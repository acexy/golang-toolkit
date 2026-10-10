package routine

import (
	"strings"

	"github.com/google/uuid"
	"github.com/timandy/routine"
)

// ThreadLocal ThreadLocalStorage 线程本地存储
type ThreadLocal[T any] struct {
	thread routine.ThreadLocal[*threadLocalValue[T]]
}

// Supplier 创建初始化值的提供者
type Supplier[T any] func() T

// LazySupplier 创建延迟初始化值的提供者。
// 函数会在当前 goroutine 首次调用 Get 时执行一次，结果会被缓存。
// supplier 内不得调用同一 ThreadLocal 的 Get、Set、SetLazy 或 Delete 方法。
type LazySupplier[T any] func() T

type threadLocalValue[T any] struct {
	value       T
	supplier    LazySupplier[T]
	initialized bool
}

// Set 设置值
func (l *ThreadLocal[T]) Set(value T) {
	l.thread.Set(&threadLocalValue[T]{
		value:       value,
		initialized: true,
	})
}

// SetLazy 设置当前 goroutine 的延迟初始化器。
// supplier 不会在调用 SetLazy 时执行，而会在首次调用 Get 时执行一次。
func (l *ThreadLocal[T]) SetLazy(supplier LazySupplier[T]) {
	l.thread.Set(&threadLocalValue[T]{
		supplier:    supplier,
		initialized: supplier == nil,
	})
}

// Get 获取当前 goroutine 的值。
// 如果值由 LazySupplier 提供，则首次获取时初始化并缓存，后续获取直接返回缓存值。
func (l *ThreadLocal[T]) Get() T {
	localValue := l.thread.Get()
	if localValue.initialized {
		return localValue.value
	}

	if localValue.supplier != nil {
		localValue.value = localValue.supplier()
	}
	localValue.supplier = nil
	localValue.initialized = true
	return localValue.value
}

// Delete 删除当前 goroutine 的值。
// 后续调用 Get 会再次使用构造时的默认 supplier（若存在）。
func (l *ThreadLocal[T]) Delete() {
	l.thread.Remove()
}

// GetGoroutineId 获取当前协程id
func GetGoroutineId() uint64 {
	return routine.Goid()
}

// NewThreadLocal 创建线程本地存储。
// 可选 supplier 是当前 goroutine 未设置值时使用的默认延迟初始化器。
func NewThreadLocal[T any](supplier ...Supplier[T]) *ThreadLocal[T] {
	var initialSupplier LazySupplier[T]
	if len(supplier) > 0 {
		initialSupplier = LazySupplier[T](supplier[0])
	}

	return &ThreadLocal[T]{
		thread: routine.NewThreadLocalWithInitial(func() *threadLocalValue[T] {
			return &threadLocalValue[T]{
				supplier:    initialSupplier,
				initialized: initialSupplier == nil,
			}
		}),
	}
}

// ======= 默认Local 场景处理器

// NewTraceIdThreadLocal 基于GoRoutine的TraceId处理器 如果不指定supplier则使用默认策略 uuid
func NewTraceIdThreadLocal(supplier Supplier[string]) *ThreadLocal[string] {
	if supplier == nil {
		supplier = func() string {
			return strings.ReplaceAll(uuid.NewString(), "-", "")
		}
	}
	return NewThreadLocal(supplier)
}
