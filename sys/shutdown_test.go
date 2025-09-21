package sys

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestShutdownHolding(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("get kill sig")
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()
	ShutdownHolding()
	fmt.Println("shutdown")
}

func TestShutdownCallback(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("get kill sig")
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()
	ShutdownCallback(func() {
		time.Sleep(10 * time.Second)
		fmt.Println("callback finish")
	})
	fmt.Println("shutdown")
}

func TestShutdownCallbackDeadline(t *testing.T) {
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("get kill sig")
		_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()
	ShutdownCallbackDeadline(func() {
		time.Sleep(5 * time.Second)
		fmt.Println("callback finish")
	}, time.Second)
	fmt.Println("shutdown")
}

func TestShutdownSignal(t *testing.T) {
	for {
		select {
		case <-ShutdownSignal():
			fmt.Println("get kill sig")
			return
		case <-time.After(time.Second):
			fmt.Println(1)
		}
	}
}
