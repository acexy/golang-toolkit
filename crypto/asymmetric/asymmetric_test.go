package asymmetric

import (
	"crypto"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsaEncrypt(t *testing.T) {
	var manager = RsaKeyManager{
		CreateSetting: CreateRsaSetting{Length: 512},
	}
	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	encrypt := NewRsaEncryptWithPKCS1()
	raw := []byte("hello rsa")
	result, err := encrypt.Encrypt(keyPair, raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err = encrypt.Decrypt(keyPair, result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(result))
	bae64Result, _ := encrypt.EncryptBase64(keyPair, base64.StdEncoding.EncodeToString(raw))
	bae64Result, _ = encrypt.DecryptBase64(keyPair, bae64Result)
	fmt.Println(string(result))

}

func TestEncryptRSAOAEP(t *testing.T) {
	var manager = RsaKeyManager{
		CreateSetting: CreateRsaSetting{Length: 2048},
	}
	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	label := []byte("label")
	encrypt, err := NewRsaEncryptWithOAEP(sha1.New(), label)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	raw := []byte("hello rsa")
	result, err := encrypt.Encrypt(keyPair, raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err = encrypt.Decrypt(keyPair, result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(result))
}

func TestRsaSign(t *testing.T) {
	var manager = RsaKeyManager{
		CreateSetting: CreateRsaSetting{Length: 2048},
	}
	keyPair, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sign := NewRsaSignWithPKCS1AndSHA256()
	raw := []byte("hello rsa")

	result, err := sign.Sign(keyPair, raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sign.Verify(keyPair, raw, result)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("NewRsaSignWithPKCS1AndSHA256 -> pass")

	sign = NewRsaSignWithPKCS1AndSHA512()
	base64Raw := base64.StdEncoding.EncodeToString(raw)
	s, err := sign.SignBase64(keyPair, base64Raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sign.VerifyBase64(keyPair, base64Raw, s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("NewRsaSignWithPKCS1AndSHA512 -> pass")

	sign, err = NewRsaSignWithPSSAndOps(md5.New(), crypto.MD5, rsa.PSSSaltLengthAuto)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	s, err = sign.SignBase64(keyPair, base64Raw)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sign.VerifyBase64(keyPair, base64Raw, s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("NewRsaSignWithPSSAndOps -> pass")
}

func TestEcdsaSignRS(t *testing.T) {
	manager := EcdsaKeyManager{CreateSetting: CreateEcdsaSetting{Curve: elliptic.P256()}}
	keyPair, _ := manager.Create()
	sign := NewEcdsaSign(crypto.SHA256.New())
	raw := []byte("你好")
	r, s, _ := sign.SignRS(keyPair, raw)
	fmt.Println(r, s)
	fmt.Println(sign.VerifyRS(keyPair, raw, r, s))

	sign = NewEcdsaSign(nil)
	r, s, _ = sign.SignRS(keyPair, raw)
	fmt.Println(r, s)
	fmt.Println(sign.VerifyRS(keyPair, raw, r, s))
}

func TestEcdsaSign(t *testing.T) {
	manager := EcdsaKeyManager{CreateSetting: CreateEcdsaSetting{Curve: elliptic.P256()}}
	keyPair, _ := manager.Create()
	sign := NewEcdsaSign(crypto.SHA256.New())
	raw := []byte("你好")
	result, _ := sign.Sign(keyPair, raw)
	fmt.Println(result)
	fmt.Println(sign.Verify(keyPair, raw, result))

	sign = NewEcdsaSign(nil)
	result, _ = sign.Sign(keyPair, raw)
	fmt.Println(result)
	fmt.Println(sign.Verify(keyPair, raw, result))
}
