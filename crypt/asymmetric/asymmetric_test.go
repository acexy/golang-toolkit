package asymmetric

import (
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestRSA(t *testing.T) {
	asymmetricRsa := NewRsaWithPaddingPKCS1(512)
	pubK, _ := asymmetricRsa.create().PublicKey().(rsa.PublicKey)
	fmt.Println(pubK)
}
