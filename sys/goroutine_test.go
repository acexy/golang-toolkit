package sys

import (
	"fmt"
	"testing"
	"time"
)

func TestGoroutineID(t *testing.T) {
	fmt.Println(GoroutineID())
	for i := 0; i < 2; i++ {
		go func() {
			fmt.Println(GoroutineID())
			fun()
		}()
	}
	fmt.Println(GoroutineID())
	time.Sleep(time.Second)
}

func fun() {
	fmt.Println(GoroutineID())
}
