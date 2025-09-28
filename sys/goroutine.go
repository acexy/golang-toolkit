package sys

import (
	"strings"

	"github.com/google/uuid"
	"github.com/timandy/routine"
)

// ThreadLocal ThreadLocalStorage 线程本地存储
type ThreadLocal[T any] struct {
	thread routine.ThreadLocal[T]
}

// Supplier 创建初始化值的提供者
type Supplier[T any] func() T

// Set 设置值
func (l *ThreadLocal[T]) Set(value T) {
	l.thread.Set(value)
}

// Get 获取值
func (l *ThreadLocal[T]) Get() T {
	return l.thread.Get()
}

// Delete 删除值
func (l *ThreadLocal[T]) Delete() {
	l.thread.Remove()
}

// GetGoroutineId 获取当前协程id
func GetGoroutineId() uint64 {
	return routine.Goid()
}

// NewThreadLocal 创建线程本地存储
func NewThreadLocal[T any](supplier ...Supplier[T]) *ThreadLocal[T] {
	if len(supplier) > 0 {
		var f routine.Supplier[T] = func() T {
			return supplier[0]()
		}
		return &ThreadLocal[T]{
			thread: routine.NewThreadLocalWithInitial[T](f),
		}
	}
	return &ThreadLocal[T]{
		thread: routine.NewThreadLocal[T](),
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
