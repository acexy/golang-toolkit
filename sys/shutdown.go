package sys

import (
	"context"
	"fmt"
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

// ShutdownSignal 监听指定的信号，若不传递则使用默认信号
// 方法会一直阻塞直到触发所监听的信号为止，并返回一个通道
// ** 如果使用 for select ，强烈建议不要在使用 case: <-ShutdownSignal() 特别是高频循环中，这样会一直创建新的简单信号管道
// 应当定义一个全局的通道 变了sig := ShutdownSignal()，然后使用 case: <-sig
func ShutdownSignal(sig ...os.Signal) <-chan struct{} {
	chn := make(chan struct{})
	go func() {
		holding(sig...)
		fmt.Println("shutdown")
		chn <- struct{}{}
	}()
	return chn
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
	case <-done:
	}
}
