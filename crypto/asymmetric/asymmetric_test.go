package asymmetric

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRSAPKCS1(t *testing.T) {

	var manager = RsaKeyManager{
		CreateSetting: CreateSetting{Length: 512},
	}

	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rsaPKCS1 := NewRsaEncryptWithPaddingPKCS1()

	raw := []byte("hello rsa")
	result, err := rsaPKCS1.Encrypt(keyPair, raw)
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
	bae64Result, _ := rsaPKCS1.EncryptBase64(keyPair, base64.StdEncoding.EncodeToString(raw))
	bae64Result, _ = rsaPKCS1.DecryptBase64(keyPair, bae64Result)
	fmt.Println(string(result))

}

func TestRSAOAEP(t *testing.T) {
	var manager = RsaKeyManager{
		CreateSetting: CreateSetting{Length: 2048},
	}

	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	label := []byte("label")
	rsaOAEP, err := NewRsaEncryptWithPaddingOAEP(sha1.New(), label)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	raw := []byte("hello rsa")
	result, err := rsaOAEP.Encrypt(keyPair, raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err = rsaOAEP.Decrypt(keyPair, result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(result))
}
