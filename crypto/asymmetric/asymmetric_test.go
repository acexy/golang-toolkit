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
		CreateSetting: CreateRsaSetting{Length: 2048},
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

func TestLoadPem(t *testing.T) {
	var manager = RsaKeyManager{}
	load, err := manager.Load("", `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCt7ki7Zpu1sYxs
87Q2eZK2kJk7PPQANEQm9+km2FRUxi16pZvRjlR7AFhaY24JlwAKrubrOvEhSn4A
NCdeTpyB9CLAj18ug8KoR3dBO5kJWYXGfSENGQB1GlGh1pFQqtdzByut3cN+KaH0
Q2jKmzdVfRuasf7Y0/3vvQpJ5owz3Nu0ynrCXT7lRd5LIxRcVNPTyci9/t6gkigk
ysrtuMDX+JUdtsGpxEEO+DhwG0BUdLQ55ZV6yCkJHTPKo2Iksy5IX/EDV97h2Oen
mdj0zsP9GX5pSEe53vOpkH5P/ucJm2LZ1VtEPk2qrxyHxzBzmcO11YC5XcxUsz60
eZ785ojNAgMBAAECggEAft0BG/+Zr4tTe9y+I8XFJ3iM69JnvwZgd3P4EadUy4jJ
f13wK4d/Q0BbMYK4rjH/i8tLD2SsoLWu9gMUfAsmyPsDpiRtLoZYo1mNOH16AIHV
u5NlWiJUXRXXWQ4ibA5Qm8wtJeJNR2aihIHum2lfkPFyAwG1ev4ViJoAOSr6NJC9
pKEUEPuz/3iJucsL3irEeb0BX/uzYL3lW6AVP+GKTv3XqAvw9Mzwvteg2x0nv/W2
5MwI+7exqBHGRkwFZ8nCMjnVdH7qHHRVzLJIfHAA4IP4OyrYNgxOYu9NqtCgX7M1
4TPxdIWtLVV1MJ7zlrDKLAnm/W4GjcwuKE4463zQGQKBgQDXWjUFhh9gfsCCDVqI
9ModYFPgcE53Mt9bO5QmiTHv2DS6zO+xpOBBUEUxJbsEIEiamb23JCj2jigIh/dp
HWN6U3nf1aVO2yFPbN79paKAq3okSzTli+cKHA+85BouyDcyRGhJ5iGlTQCIkeG/
eT7i4klXbvgZ3mhAxC8QdHmvGwKBgQDOwpfPhEi5a5sQlr5Zgb91zu8a/hzHoYg5
ahmfLTulitV6uMg2xxbueDAUSh5pkajQChICQ42Mj0+0uRxaQz42VESXkNXMTLBS
FBssSkaCGmHQbyPvYvk/tERljcqs4Pa85MMi/sTEHMKVQZFRM2Yhdj6x5n3h3obo
nxcbIsJeNwKBgQCQEDacJWdkJCcgWVKRgECdek1iPX8gWpX08FxhkzIm4xoTRPms
e0HSL6K6CCWd0wL++Y5ir/v/RIYffuXSGejscl97j+7wW8Ni4NAkGuQk85HYKb+P
OBItPyyADpA9b48NP5oMpbkoXeTXd8/vSWr6WKr7pK6wj4pqmVNqzKCCtQKBgDFR
GQTWeMJBeGssiZqv6AshoMa/df3n+aY4OIRPNbr2spTjHl5yfnXDFTTBuR9VLv7w
Z/tCQbEVPd3NiYW+JEixMOs2EMj6QdRSF2kpDkhaIpqk591hrjITvXy7pWw2/KNx
rVnbivN8KK2RRigoKmQw4CNu5vfJLESwLcK7N8FhAoGBAMbq1Zr/d8B8hEoJcWzf
922krJNd9uZJQJWMc4HBHrgcm2GlODCjcpmyL2oM9udy1adsrojySM7dyCD8my+B
bUhXNlXx4OW3sWEi9ebaNlV+Qz5IyFnC4J3DvGGtm6BBTisU5rMDY+Uv7iW8VBDH
4TlS4PKZixp/cFv2jW1F+tSX
-----END PRIVATE KEY-----`)
	if err != nil {
		println(err)
		return
	}
	fmt.Println(load.ToPrivatePKCS1Pem())
}
