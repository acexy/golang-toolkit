package sys

import (
	"fmt"
	"testing"
)

func TestSetGOMAXPROCS(t *testing.T) {
	fmt.Println(DetectCPULimit())
	SetGoMaxProc(1)
}
