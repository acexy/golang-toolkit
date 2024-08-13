package asymmetric

import (
	"crypto"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/acexy/golang-toolkit/math/conversion"
	"testing"
)

func TestRsaKey(t *testing.T) {
	var manager = RsaKeyManager{
		CreateSetting: CreateRsaSetting{Length: 512},
	}
	key, err := manager.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 打印PEM编码的公共密钥
	fmt.Printf("PEM格式的公钥:\n%s\n", key.ToPublicPKCS1Pem())
	fmt.Printf("PEM格式的私钥:\n%s\n", key.ToPrivatePKCS1Pem())
}

func TestRsaEncrypt(t *testing.T) {
	//var manager = RsaKeyManager{
	//	CreateSetting: CreateRsaSetting{Length: 512},
	//}
	//keyPair, err := manager.Create()
	var manager = RsaKeyManager{}
	keyPair, err := manager.Load(`-----BEGIN RSA PUBLIC KEY-----
MEgCQQDdfMR8nyWFdEyg2aWkM4QjcTCvR7gqjGBo5rEASOOoCP52VfmRH686Koen
nTnq46LL3TjvJf0Q52tTKBj3X16BAgMBAAE=
-----END RSA PUBLIC KEY-----`, `-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBAN18xHyfJYV0TKDZpaQzhCNxMK9HuCqMYGjmsQBI46gI/nZV+ZEf
rzoqh6edOerjosvdOO8l/RDna1MoGPdfXoECAwEAAQJBAM6R5iOgvnroS+uc8irh
zTTNBa4EgtRUFjrgJWbxlDoLZPUZq+ckSpivuVdHZWfjsIJ7M0kzYWs4BBpgzbay
xH0CIQDmRdVdf9xM/Sa0+oUfaFWrHD9gv0jNLlJ3GlKOQ6KVOwIhAPY7qEXJGTEw
hRxZEN0EOjVtyVYlzzITxXYVAPtyya9zAiBR8eP+A/RHyYauvMAG70AdRk4fhbLI
oYVjMRDT46nF5QIhALFquMtXo6w6rp6HWkw1wI9AxKIq6gjGEDAN4EBNLB8bAiEA
xniLzimI7cQ8x+phTalMOazWfXdtnVcWxfWzQAVTaFE=
-----END RSA PRIVATE KEY-----`)
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
	decodeString, _ := base64.StdEncoding.DecodeString(`Kp2ws5r7Ia56s5dpypk0BIAUrXt7FtFP4sDP8ppGRF+kBcluiB/4gezea7ar4H6be8cN/gUIQMPLJncgeRJ8dg==`)
	fmt.Println(decodeString)
	fmt.Println(conversion.ParesBytesFromHex("2a9db0b39afb21ae7ab39769ca9934048014ad7b7b16d14fe2c0cff29a46445fa405c96e881ff881ecde6bb6abe07e9b7bc70dfe050840c3cb26772079127c76"))
	result, err = encrypt.Decrypt(keyPair, conversion.ParesBytesFromHex("305c300d06092a864886f70d0101010500034b003048024100dd7cc47c9f2585744ca0d9a5a43384237130af47b82a8c6068e6b10048e3a808fe7655f9911faf3a2a87a79d39eae3a2cbdd38ef25fd10e76b532818f75f5e810203010001"))
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

func TestEcdsaKey(t *testing.T) {
	manager := EcdsaKeyManager{CreateSetting: CreateEcdsaSetting{Curve: elliptic.P256()}}
	key, _ := manager.Create()
	// 打印PEM编码的公共密钥
	fmt.Printf("PEM格式的公钥:\n%s\n", key.ToPublicPKCS1Pem())
	fmt.Printf("PEM格式的私钥:\n%s\n", key.ToPrivatePKCS1Pem())
}

func TestEcdsaSignRS(t *testing.T) {
	//manager := EcdsaKeyManager{CreateSetting: CreateEcdsaSetting{Curve: elliptic.P256()}}
	//keyPair, _ := manager.Create()
	manager := EcdsaKeyManager{}
	keyPair, err := manager.Load(`-----BEGIN EC PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEaEgWXspJ7qHq3ZAC811gnXGECvgz
a7yDsRYEOjeUjKzRe2VqQew39pYpkLGtUo4HY63NIWs5vvrDutsqTOwMFw==
-----END EC PUBLIC KEY-----`, `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKThaZ1D5ut5opXP5SDvTXi5TnKcoSnTnAfj/xPyVM0moAoGCCqGSM49
AwEHoUQDQgAEaEgWXspJ7qHq3ZAC811gnXGECvgza7yDsRYEOjeUjKzRe2VqQew3
9pYpkLGtUo4HY63NIWs5vvrDutsqTOwMFw==
-----END EC PRIVATE KEY-----`)
	if err != nil {
		println(err)
		return
	}
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
