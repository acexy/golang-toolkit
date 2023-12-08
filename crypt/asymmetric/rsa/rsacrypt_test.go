package rsa

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	piKey, pbKey := Generate(512)
	fmt.Println(piKey, pbKey)
}
