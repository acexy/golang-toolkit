package asymmetric

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRSA(t *testing.T) {

	var manager = RsaKeyManager{
		CreateSetting: CreateSetting{Length: 168},
	}

	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rsaPKCS1 := NewRsaWithPaddingPKCS1()

	context := []byte("hello rsa")
	result, err := rsaPKCS1.Encrypt(keyPair, context)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err = rsaPKCS1.Decrypt(keyPair, result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(result))
	bae64Result, _ := rsaPKCS1.EncryptBase64(keyPair, base64.StdEncoding.EncodeToString(context))
	bae64Result, _ = rsaPKCS1.DecryptBase64(keyPair, bae64Result)
	fmt.Println(string(result))

}
