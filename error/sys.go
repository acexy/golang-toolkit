package error

import "errors"

var (
	// ErrShutdownSignalsInitialized 表示底层信号订阅已初始化，不能再修改关闭信号。
	ErrShutdownSignalsInitialized = errors.New("shutdown signals already initialized")

	// ErrInvalidShutdownSignals 表示关闭信号为空或包含 nil。
	ErrInvalidShutdownSignals = errors.New("shutdown signals cannot be empty or contain nil")
)
