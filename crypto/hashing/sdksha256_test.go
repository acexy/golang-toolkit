package hashing

import (
	"fmt"
	"github.com/acexy/golang-toolkit/math/conversion"
	"testing"
)

func TestSha256(t *testing.T) {
	fmt.Println(Sha256Hex("sha256"))
	fmt.Println(Sha256Base64("sha256"))
	fmt.Println(Sha256Bytes(conversion.ParseBytes("sha256")))
}
