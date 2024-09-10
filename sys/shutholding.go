package sys

import (
	"context"
	"github.com/acexy/golang-toolkit/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var defaultSig = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGALRM}

func holding(sig ...os.Signal) os.Signal {
	ch := make(chan os.Signal, 1)
	defer signal.Stop(ch)
	if len(sig) == 0 {
		signal.Notify(ch, defaultSig...)
	} else {
		signal.Notify(ch, sig...)
	}
	return <-ch // 返回捕获的信号
}

// ShutdownHolding 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止
func ShutdownHolding(sig ...os.Signal) {
	holding(sig...)
}

// ShutdownCallback 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止，并执行回调
func ShutdownCallback(f func(), sig ...os.Signal) {
	holding(sig...)
	if f != nil {
		f()
	}
}

// ShutdownCallbackDeadline 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止 并执行回调 若在指定时间未完成回调执行，则放弃等待
func ShutdownCallbackDeadline(f func(), deadline time.Duration, sig ...os.Signal) {
	holding(sig...)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	done := make(chan struct{})
	go func() {
		if f != nil {
			f()
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		logger.Logrus().Warningln("shutdown callback timeout")
	case <-done:
	}
}
