package sys

import (
	"github.com/acexy/golang-toolkit/math/random"
	"github.com/timandy/routine"
	"sync"
)

// 一个traceId的默认共享策略
var traceIdLocal *Local[string]
var traceIdLocalOnce sync.Once

// Local ThreadLocalStorage 线程本地存储
type Local[T any] struct {
	thread routine.ThreadLocal[T]
}

// Supplier 创建初始化值的提供者
type Supplier[T any] func() T

// Set 设置值
func (l *Local[T]) Set(value T) {
	l.thread.Set(value)
}

// Get 获取值
func (l *Local[T]) Get() T {
	return l.thread.Get()
}

// Delete 删除值
func (l *Local[T]) Delete() {
	l.thread.Remove()
}

// GetGoroutineId 获取当前协程id
func GetGoroutineId() int64 {
	return routine.Goid()
}

// NewThreadLocal 创建线程本地存储
func NewThreadLocal[T any](supplier ...Supplier[T]) *Local[T] {
	if len(supplier) > 0 {
		var f routine.Supplier[T] = func() T {
			return supplier[0]()
		}
		return &Local[T]{
			thread: routine.NewThreadLocalWithInitial[T](f),
		}
	}
	return &Local[T]{
		thread: routine.NewThreadLocal[T](),
	}
}

// ======= 默认Local 场景处理器

// EnableLocalTraceId 激活TraceIdLocal处理器 如果不指定supplier则使用默认策略 uuid
func EnableLocalTraceId(supplier Supplier[string]) {
	traceIdLocalOnce.Do(func() {
		if supplier == nil {
			supplier = func() string {
				return random.UUID()
			}
		}
		traceIdLocal = NewThreadLocal(supplier)
	})
}

// IsEnabledLocalTraceId 判断是否启用了TraceIdLocal处理器
func IsEnabledLocalTraceId() bool {
	return traceIdLocal != nil
}

// GetLocalTraceId 获取当前线程的TraceId
func GetLocalTraceId() string {
	if traceIdLocal == nil {
		return ""
	}
	return traceIdLocal.Get()
}
