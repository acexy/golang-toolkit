package sys

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

var defaultShutdownSig = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

var shutdownSignalMu sync.Mutex
var shutdownSignals = defaultShutdownSig
var shutdownSignalsInitialized bool
var exitSignChan = make(chan struct{})
var exitOnce sync.Once

func holding(sig ...os.Signal) os.Signal {
	ch := make(chan os.Signal, 1)
	defer signal.Stop(ch)
	if len(sig) == 0 {
		signal.Notify(ch, defaultShutdownSig...)
	} else {
		signal.Notify(ch, sig...)
	}
	return <-ch // 返回捕获的信号
}

// SetShutdownSignals 设置进程级关闭信号。
// 该方法必须在首次初始化底层信号订阅前调用，且信号不能为空或包含 nil。
func SetShutdownSignals(sig ...os.Signal) error {
	if len(sig) == 0 {
		return toolkitError.ErrInvalidShutdownSignals
	}
	for _, item := range sig {
		if item == nil {
			return toolkitError.ErrInvalidShutdownSignals
		}
	}

	shutdownSignalMu.Lock()
	defer shutdownSignalMu.Unlock()
	if shutdownSignalsInitialized {
		return toolkitError.ErrShutdownSignalsInitialized
	}
	shutdownSignals = append([]os.Signal(nil), sig...)
	return nil
}

// ShutdownHolding 阻塞等待进程级共享的关闭信号。
func ShutdownHolding() {
	<-ShutdownSignal()
}

// ShutdownCallback 阻塞等待进程级共享的关闭信号，收到信号后同步执行回调。
func ShutdownCallback(f func()) {
	<-ShutdownSignal()
	if f != nil {
		f()
	}
}

// ShutdownSignal 返回进程级共享的关闭信道。
// 底层仅订阅一次系统信号；未调用 SetShutdownSignals 时使用默认关闭信号。
// 收到任一监听信号后，共享信道会被关闭并通知所有监听者。
func ShutdownSignal() <-chan struct{} {
	exitOnce.Do(func() {
		shutdownSignalMu.Lock()
		shutdownSignalsInitialized = true
		sig := append([]os.Signal(nil), shutdownSignals...)
		shutdownSignalMu.Unlock()
		go func() {
			holding(sig...)
			close(exitSignChan)
		}()
	})
	return exitSignChan
}

// ShutdownContext 返回一个可取消的 context。
// 父 context、调用方主动取消或进程收到关闭信号时，返回的 context 会被取消。
// 该方法复用进程级共享的信号订阅。
func ShutdownContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctxNew, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-ShutdownSignal():
			cancel()
		case <-ctxNew.Done():
		}
	}()
	return ctxNew, cancel
}

// ShutdownCallbackDeadline 阻塞等待进程级共享的关闭信号，收到信号后执行回调。
// 超过指定等待时间后方法返回，但不会强制终止仍在执行的回调。
func ShutdownCallbackDeadline(f func(), deadline time.Duration) {
	<-ShutdownSignal()
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if f != nil {
			f()
		}
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}
