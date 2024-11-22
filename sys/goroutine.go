package sys

import "github.com/petermattis/goid"

func GoroutineID() int64 {
	return goid.Get()
}
