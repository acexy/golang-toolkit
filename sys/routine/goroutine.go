package routine

import (
	"github.com/timandy/routine"
)

// ThreadLocal stores a typed value isolated to the current goroutine.
type ThreadLocal[T any] struct {
	threadLocalStorage[T]
}

// InheritableThreadLocal supports propagation through Run and RunAll; native go does not inherit values.
type InheritableThreadLocal[T any] struct {
	threadLocalStorage[T]
}

// threadLocalStorage shares accessors without exposing replaceable storage.
type threadLocalStorage[T any] struct {
	thread routine.ThreadLocal[*threadLocalValue[T]]
}

// Supplier provides a default value for the current goroutine.
type Supplier[T any] func() T

// LazySupplier provides a value on the first Get in the current goroutine.
// It must not reenter the same local through Get, Set, SetLazy, or Delete.
type LazySupplier[T any] func() T

type threadLocalValue[T any] struct {
	value       T
	supplier    LazySupplier[T]
	initialized bool
}

// Clone copies lazy initialization state; referenced value data is not deep-copied.
func (v *threadLocalValue[T]) Clone() any {
	return v.clone()
}

func (v *threadLocalValue[T]) clone() *threadLocalValue[T] {
	cloned := *v
	return &cloned
}

// Set replaces the current goroutine value without evaluating a lazy supplier.
func (l *threadLocalStorage[T]) Set(value T) {
	l.thread.Set(&threadLocalValue[T]{
		value:       value,
		initialized: true,
	})
}

// SetLazy sets a supplier to evaluate and cache on the next Get; nil supplies the zero value.
func (l *threadLocalStorage[T]) SetLazy(supplier LazySupplier[T]) {
	l.thread.Set(&threadLocalValue[T]{
		supplier:    supplier,
		initialized: supplier == nil,
	})
}

// Get returns the current value, evaluating and caching its lazy supplier if necessary.
func (l *threadLocalStorage[T]) Get() T {
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

// Delete removes the current value so the next Get uses the default supplier again.
func (l *threadLocalStorage[T]) Delete() {
	l.thread.Remove()
}

// GetGoroutineId returns the current goroutine identifier.
func GetGoroutineId() uint64 {
	return routine.Goid()
}

// NewThreadLocal creates isolated storage with an optional lazy default supplier. Reuse the instance across goroutines.
func NewThreadLocal[T any](supplier ...Supplier[T]) *ThreadLocal[T] {
	return &ThreadLocal[T]{
		threadLocalStorage: threadLocalStorage[T]{
			thread: routine.NewThreadLocalWithInitial(initialValueSupplier(supplier...)),
		},
	}
}

func initialValueSupplier[T any](supplier ...Supplier[T]) routine.Supplier[*threadLocalValue[T]] {
	var initialSupplier LazySupplier[T]
	if len(supplier) > 0 {
		initialSupplier = LazySupplier[T](supplier[0])
	}

	return func() *threadLocalValue[T] {
		return &threadLocalValue[T]{
			supplier:    initialSupplier,
			initialized: initialSupplier == nil,
		}
	}
}

// NewInheritableThreadLocal creates inheritable storage with an optional lazy default supplier.
// Uninitialized suppliers run independently in each goroutine; synchronize shared closure data.
// Reuse instances to avoid increasing storage capacity and RunAll copying costs.
func NewInheritableThreadLocal[T any](supplier ...Supplier[T]) *InheritableThreadLocal[T] {
	return &InheritableThreadLocal[T]{
		threadLocalStorage: threadLocalStorage[T]{
			thread: routine.NewInheritableThreadLocalWithInitial(initialValueSupplier(supplier...)),
		},
	}
}

// Run starts a goroutine with only this instance snapshot, without evaluating its supplier.
// Child Set, SetLazy, and Delete leave the parent wrapper unchanged; referenced data is shared.
// Task panics are caught and printed, as with RunAll.
func (l *InheritableThreadLocal[T]) Run(fn func()) {
	value := l.thread.Get().clone()
	go func() {
		// Wrap inside the child so other parent locals are not captured.
		routine.WrapTask(func() {
			l.thread.Set(value)
			fn()
		}).Run()
	}()
}

// RunAll starts a goroutine with all saved inheritable instances from the caller.
// Ordinary locals are excluded and unused instances are not initialized.
// Capture occurs at call time; referenced data is shared and task panics are caught and printed.
func (l *InheritableThreadLocal[T]) RunAll(fn func()) {
	routine.Go(fn)
}
