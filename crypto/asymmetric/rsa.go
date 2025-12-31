package asymmetric

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"sync"

	"github.com/acexy/golang-toolkit/math/conversion"
)

type PaddingType uint8

const (
	// PaddingTypeNone 无填充
	PaddingTypeNone PaddingType = 0
	// PaddingTypePKCS1 PKCS1.5
	PaddingTypePKCS1 PaddingType = 1
	// PaddingTypeOAEP OAEP模式
	PaddingTypeOAEP PaddingType = 2
	// PaddingTypePSS PSS模式
	PaddingTypePSS PaddingType = 3
)

type rsaKey struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *rsaKey) PrivateKey() interface{} {
	return r.privateKey
}

func (r *rsaKey) PublicKey() interface{} {
	return r.publicKey
}

func (r *rsaKey) ToPublicPKCS1Pem() string {
	if r.publicKey == nil {
		return ""
	}
	publicKey := r.PublicKey().(*rsa.PublicKey)
	der := x509.MarshalPKCS1PublicKey(publicKey)
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	}))
}

func (r *rsaKey) ToPrivatePKCS1Pem() string {
	if r.privateKey == nil {
		return ""
	}
	privateKey := r.PrivateKey().(*rsa.PrivateKey)
	der := x509.MarshalPKCS1PrivateKey(privateKey)
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	}))
}

func (r *rsaKey) ToPublicPKCS8Pem() (string, error) {
	if r.publicKey == nil {
		return "", errors.New("nil public key")
	}
	publicKey := r.PublicKey().(*rsa.PublicKey)
	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})), nil
}

func (r *rsaKey) ToPrivatePKCS8Pem() string {
	if r.privateKey == nil {
		return ""
	}
	privateKey := r.PrivateKey().(*rsa.PrivateKey)
	der := x509.MarshalPKCS1PrivateKey(privateKey)
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}))
}

type CreateRsaSetting struct {
	Length int
}

func NewRsaKeyManager(keyLength int) *RsaKeyManager {
	return &RsaKeyManager{
		CreateSetting: CreateRsaSetting{
			Length: keyLength,
		},
	}
}

func NewEmptyRasKeyManager() *RsaKeyManager {
	return &RsaKeyManager{}
}

type RsaKeyManager struct {
	CreateSetting CreateRsaSetting
}

func (r *RsaKeyManager) Create() (KeyPair, error) {
	if r.CreateSetting.Length == 0 {
		return nil, errors.New("bad key length")
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, r.CreateSetting.Length)
	if err != nil {
		return nil, err
	}
	return &rsaKey{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
	}, nil
}

func (r *RsaKeyManager) Load(pubPem, priPem string) (KeyPair, error) {
	if pubPem == "" && priPem == "" {
		return nil, errors.New("bad key")
	}

	var pub *rsa.PublicKey
	var pri *rsa.PrivateKey
	var err error

	// 解析公钥
	if pubPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(pubPem))
		if block == nil {
			return nil, errors.New("bad public key")
		}

		switch block.Type {
		case "RSA PUBLIC KEY":
			// PKCS#1 公钥
			pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PUBLIC KEY":
			// PKIX (PKCS#8) 公钥
			var iface any
			iface, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pub, ok = iface.(*rsa.PublicKey)
			if !ok {
				return nil, errors.New("not an RSA public key")
			}
		default:
			return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
		}
	}

	// 解析私钥
	if priPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(priPem))
		if block == nil {
			return nil, errors.New("bad private key")
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			// PKCS#1 私钥
			pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PRIVATE KEY":
			// PKCS#8 私钥
			var iface any
			iface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pri, ok = iface.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("not an RSA private key")
			}
		default:
			return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
		}
	}

	return &rsaKey{
		publicKey:  pub,
		privateKey: pri,
	}, nil
}

type RsaEncrypt struct {
	paddingType  PaddingType
	hashForOAEP  hash.Hash
	labelForOAEP []byte
}

func (r *RsaEncrypt) Encrypt(keyPair KeyPair, raw []byte) ([]byte, error) {
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return nil, errors.New("nil public key")
	}
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), raw)
	case PaddingTypeOAEP:
		return rsa.EncryptOAEP(r.hashForOAEP, rand.Reader, publicKey.(*rsa.PublicKey), raw, r.labelForOAEP)
	default:

	}
	return nil, errors.New("not supported paddingType")
}

func (r *RsaEncrypt) EncryptBase64(keyPair KeyPair, base64Raw string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	result, err := r.Encrypt(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func (r *RsaEncrypt) Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error) {
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, errors.New("nil privateKey key")
	}
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), cipher)
	case PaddingTypeOAEP:
		return rsa.DecryptOAEP(r.hashForOAEP, rand.Reader, privateKey.(*rsa.PrivateKey), cipher, r.labelForOAEP)
	default:

	}
	return nil, errors.New("not supported paddingType")
}

func (r *RsaEncrypt) DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return "", err
	}
	result, err := r.Decrypt(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

type RsaSign struct {
	sync.Mutex
	paddingType     PaddingType
	hashForOAEP     hash.Hash
	labelForOAEP    []byte
	hashForSign     hash.Hash
	hashTypeForSign crypto.Hash
	options         *rsa.PSSOptions
}

func (r *RsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	defer r.Unlock()
	r.Lock()
	if r.hashForSign == nil {
		return nil, errors.New("nil hash function")
	}
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, errors.New("nil privateKey key")
	}
	r.hashForSign.Reset()
	r.hashForSign.Write(raw)
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), r.hashTypeForSign, r.hashForSign.Sum(nil))
	case PaddingTypePSS:
		return rsa.SignPSS(rand.Reader, privateKey.(*rsa.PrivateKey), r.hashTypeForSign, r.hashForSign.Sum(nil), r.options)
	default:

	}
	return nil, errors.New("not supported paddingType")
}

// Verify 签名验证
func (r *RsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	defer r.Unlock()
	r.Lock()
	if r.hashForSign == nil {
		return errors.New("nil hash function")
	}
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return errors.New("nil public key")
	}
	r.hashForSign.Reset()
	r.hashForSign.Write(raw)
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), r.hashTypeForSign, r.hashForSign.Sum(nil), sign)
	case PaddingTypePSS:
		return rsa.VerifyPSS(publicKey.(*rsa.PublicKey), r.hashTypeForSign, r.hashForSign.Sum(nil), sign, r.options)
	default:

	}
	return errors.New("not supported paddingType")
}

func (r *RsaSign) SignBase64(keyPair KeyPair, base64Raw string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	result, err := r.Sign(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func (r *RsaSign) VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error {
	rawContent, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return err
	}
	signContent, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}
	return r.Verify(keyPair, rawContent, signContent)
}

// NewRsaEncryptWithPKCS1 创建一个标准的RSA加密 Padding-PKCS1模式实例
func NewRsaEncryptWithPKCS1() *RsaEncrypt {
	var encrypt RsaEncrypt
	encrypt.paddingType = PaddingTypePKCS1
	return &encrypt
}

// NewRsaEncryptWithOAEP 创建一个标准的RSA加密 Padding-OAEP模式实例
func NewRsaEncryptWithOAEP(hash hash.Hash, label []byte) (*RsaEncrypt, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var encrypt RsaEncrypt
	encrypt.paddingType = PaddingTypeOAEP
	encrypt.hashForOAEP = hash
	encrypt.labelForOAEP = label
	return &encrypt, nil
}

// NewRsaSignWithPKCS1AndSHA256 创建一个标准RSA签名 Padding-PKCS1 hash函数为sha256
func NewRsaSignWithPKCS1AndSHA256() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha256.New(), crypto.SHA256)
	return sign
}

// NewRsaSignWithPKCS1AndSHA512 创建一个标准RSA签名 Padding-PKCS1 hash函数为sha512
func NewRsaSignWithPKCS1AndSHA512() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha512.New(), crypto.SHA512)
	return sign
}

// NewRsaSignWithPKCS1 创建一个标准RSA签名 Padding-PKCS1 自定义hash函数
func NewRsaSignWithPKCS1(hash hash.Hash, hashType crypto.Hash) (*RsaSign, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var sign RsaSign
	sign.paddingType = PaddingTypePKCS1
	sign.hashForSign = hash
	sign.hashTypeForSign = hashType
	return &sign, nil
}

// NewRsaSignWithPSSAndSHA256 创建一个标准RSA签名 Padding-PSS hash函数为sha256
func NewRsaSignWithPSSAndSHA256(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOps(sha256.New(), crypto.SHA256, length)
	return sign
}

// NewRsaSignWithPSSAndSHA512 创建一个标准RSA签名 Padding-PSS hash函数为sha512
func NewRsaSignWithPSSAndSHA512(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOps(sha512.New(), crypto.SHA512, length)
	return sign
}

// NewRsaSignWithPSS 创建一个标准RSA签名 Padding-PSS
func NewRsaSignWithPSS(hash hash.Hash, hashType crypto.Hash) (*RsaSign, error) {
	return NewRsaSignWithPSSAndOps(hash, hashType, -1)
}

// NewRsaSignWithPSSAndOps NewRsaSignWithPSS 创建一个标准RSA签名 Padding-PSS 并指定saltLength
func NewRsaSignWithPSSAndOps(hash hash.Hash, hashType crypto.Hash, saltLength int) (*RsaSign, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var sign RsaSign
	sign.paddingType = PaddingTypePSS
	sign.hashForSign = hash
	sign.hashTypeForSign = hashType
	if saltLength >= 0 {
		sign.options = &rsa.PSSOptions{
			Hash:       hashType,
			SaltLength: saltLength,
		}
	}
	return &sign, nil
}
