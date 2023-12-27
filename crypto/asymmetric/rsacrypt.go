package asymmetric

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash"
	"sync"
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

type CreateRsaSetting struct {
	Length int
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

func (r *RsaKeyManager) Load() (KeyPair, error) {
	return nil, errors.New("not support now")
}

type rsaEncrypt struct {
	paddingType  PaddingType
	hashForOAEP  hash.Hash
	labelForOAEP []byte
}

func (r *rsaEncrypt) Encrypt(keyPair KeyPair, raw []byte) ([]byte, error) {
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

func (r *rsaEncrypt) EncryptBase64(keyPair KeyPair, base64Raw string) (string, error) {
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

func (r *rsaEncrypt) Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error) {
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

func (r *rsaEncrypt) DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error) {
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

type rsaSign struct {
	sync.Mutex
	paddingType     PaddingType
	hashForOAEP     hash.Hash
	labelForOAEP    []byte
	hashForSign     hash.Hash
	hashTypeForSign crypto.Hash
	options         *rsa.PSSOptions
}

func (r *rsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
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
func (r *rsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
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

func (r *rsaSign) SignBase64(keyPair KeyPair, base64Raw string) (string, error) {
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

func (r *rsaSign) VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error {
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
func NewRsaEncryptWithPKCS1() *rsaEncrypt {
	var encrypt rsaEncrypt
	encrypt.paddingType = PaddingTypePKCS1
	return &encrypt
}

// NewRsaEncryptWithOAEP 创建一个标准的RSA加密 Padding-OAEP模式实例
func NewRsaEncryptWithOAEP(hash hash.Hash, label []byte) (*rsaEncrypt, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var encrypt rsaEncrypt
	encrypt.paddingType = PaddingTypeOAEP
	encrypt.hashForOAEP = hash
	encrypt.labelForOAEP = label
	return &encrypt, nil
}

// NewRsaSignWithPKCS1AndSHA256 创建一个标准RSA签名 Padding-PKCS1 hash函数为sha256
func NewRsaSignWithPKCS1AndSHA256() *rsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha256.New(), crypto.SHA256)
	return sign
}

// NewRsaSignWithPKCS1AndSHA512 创建一个标准RSA签名 Padding-PKCS1 hash函数为sha512
func NewRsaSignWithPKCS1AndSHA512() *rsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha512.New(), crypto.SHA512)
	return sign
}

// NewRsaSignWithPKCS1 创建一个标准RSA签名 Padding-PKCS1 自定义hash函数
func NewRsaSignWithPKCS1(hash hash.Hash, hashType crypto.Hash) (*rsaSign, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var sign rsaSign
	sign.paddingType = PaddingTypePKCS1
	sign.hashForSign = hash
	sign.hashTypeForSign = hashType
	return &sign, nil
}

// NewRsaSignWithPSSAndSHA256 创建一个标准RSA签名 Padding-PSS hash函数为sha256
func NewRsaSignWithPSSAndSHA256(saltLength ...int) *rsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOps(sha256.New(), crypto.SHA256, length)
	return sign
}

// NewRsaSignWithPSSAndSHA512 创建一个标准RSA签名 Padding-PSS hash函数为sha512
func NewRsaSignWithPSSAndSHA512(saltLength ...int) *rsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOps(sha512.New(), crypto.SHA512, length)
	return sign
}

// NewRsaSignWithPSS 创建一个标准RSA签名 Padding-PSS
func NewRsaSignWithPSS(hash hash.Hash, hashType crypto.Hash) (*rsaSign, error) {
	return NewRsaSignWithPSSAndOps(hash, hashType, -1)
}

// NewRsaSignWithPSSAndOps NewRsaSignWithPSS 创建一个标准RSA签名 Padding-PSS 并指定saltLength
func NewRsaSignWithPSSAndOps(hash hash.Hash, hashType crypto.Hash, saltLength int) (*rsaSign, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var sign rsaSign
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
