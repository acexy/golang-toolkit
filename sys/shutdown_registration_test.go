package sys

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func awaitShutdownRegistration(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("等待关闭通知超时")
	}
}

func TestShutdownRegistrationScenarios(t *testing.T) {
	const marker = "shutdown-registration-scenario:"
	if scenario, child := strings.CutPrefix(os.Args[len(os.Args)-1], marker); child {
		// Isolate one-time subscription and test signal delivery before the waiter runs.
		runtime.GOMAXPROCS(1)
		selectedSignal := syscall.SIGUSR1
		if scenario == "default-signals" {
			selectedSignal = syscall.SIGTERM
		} else if err := SetShutdownSignals(selectedSignal); err != nil {
			t.Fatal(err)
		}
		trigger := func() {
			t.Helper()
			if err := syscall.Kill(os.Getpid(), selectedSignal); err != nil {
				t.Fatal(err)
			}
		}
		switch scenario {
		case "immediate-signal", "default-signals":
			done := ShutdownSignal()
			trigger()
			awaitShutdownRegistration(t, done)
		case "context-ready":
			ctx, cancel := ShutdownContext(nil)
			defer cancel()
			if err := SetShutdownSignals(syscall.SIGUSR2); !errors.Is(err, toolkitError.ErrShutdownSignalsInitialized) {
				t.Fatalf("ShutdownContext 返回后应已完成初始化: %v", err)
			}
			trigger()
			awaitShutdownRegistration(t, ctx.Done())
		case "broadcast":
			const count = 16
			channels := make(chan (<-chan struct{}), count)
			var ready sync.WaitGroup
			for i := 0; i < count; i++ {
				ready.Add(1)
				go func() {
					defer ready.Done()
					channels <- ShutdownSignal()
				}()
			}
			ready.Wait()
			shared := ShutdownSignal()
			trigger()
			for i := 0; i < count; i++ {
				ch := <-channels
				if ch != shared {
					t.Fatal("所有监听者应复用同一信道")
				}
				awaitShutdownRegistration(t, ch)
			}
			awaitShutdownRegistration(t, ShutdownSignal())
			ShutdownHolding()
		case "context-cancellation":
			parent, parentCancel := context.WithCancel(context.Background())
			ctx, cancel := ShutdownContext(parent)
			defer cancel()
			parentCancel()
			awaitShutdownRegistration(t, ctx.Done())
			other, otherCancel := ShutdownContext(nil)
			otherCancel()
			awaitShutdownRegistration(t, other.Done())
			select {
			case <-ShutdownSignal():
				t.Fatal("context 取消不应触发进程关闭广播")
			default:
			}
		case "callback":
			done := ShutdownSignal()
			trigger()
			awaitShutdownRegistration(t, done)
			started, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
			go func() {
				ShutdownCallback(func() { close(started); <-release })
				close(finished)
			}()
			awaitShutdownRegistration(t, started)
			select {
			case <-finished:
				t.Error("同步回调完成前不应返回")
			default:
			}
			close(release)
			awaitShutdownRegistration(t, finished)
			ShutdownCallback(nil)
		case "deadline":
			done := ShutdownSignal()
			trigger()
			awaitShutdownRegistration(t, done)
			started, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
			ShutdownCallbackDeadline(func() {
				close(started)
				<-release
				close(finished)
			}, 20*time.Millisecond)
			awaitShutdownRegistration(t, started)
			select {
			case <-finished:
				t.Error("超时不应终止仍在执行的回调")
			default:
			}
			close(release)
			awaitShutdownRegistration(t, finished)
			ShutdownCallbackDeadline(nil, time.Second)
		case "configuration":
			if err := SetShutdownSignals(); !errors.Is(err, toolkitError.ErrInvalidShutdownSignals) {
				t.Fatalf("空信号配置 = %v", err)
			}
			if err := SetShutdownSignals(nil); !errors.Is(err, toolkitError.ErrInvalidShutdownSignals) {
				t.Fatalf("nil 信号配置 = %v", err)
			}
			ShutdownSignal()
			if err := SetShutdownSignals(syscall.SIGUSR2); !errors.Is(err, toolkitError.ErrShutdownSignalsInitialized) {
				t.Fatalf("初始化后修改配置 = %v", err)
			}
		default:
			t.Fatal("未知场景")
		}
		return
	}
	for _, scenario := range []string{"immediate-signal", "default-signals", "context-ready", "broadcast", "context-cancellation", "callback", "deadline", "configuration"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestShutdownRegistrationScenarios$", "--", marker+scenario)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("隔离关闭测试失败: %v, timeout=%v\n%s", err, ctx.Err(), output)
			}
		})
	}
}
