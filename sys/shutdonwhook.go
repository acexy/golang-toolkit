package sys

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var defaultSig = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGALRM}

func hold(sig ...os.Signal) {
	ch := make(chan os.Signal)
	if len(sig) == 0 {
		signal.Notify(ch, defaultSig...)
	} else {
		signal.Notify(ch, sig...)
	}
	<-ch
}

// ShutdownHolding 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止
func ShutdownHolding(sig ...os.Signal) {
	hold(sig...)
}

// ShutdownCallback 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止 并执行回调
func ShutdownCallback(f func(), sig ...os.Signal) {
	hold(sig...)
	f()
}

// ShutdownCallbackDeadline 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止 并执行回调 若在指定时间未完成回调执行，则放弃等待
func ShutdownCallbackDeadline(f func(), deadline time.Duration, sig ...os.Signal) {
	hold(sig...)
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	done := make(chan bool)
	go func() {
		f()
		done <- true
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}
