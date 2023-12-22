package asymmetric

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
)

type PaddingType uint8

const (
	// PaddingTypeNone 无填充
	PaddingTypeNone PaddingType = 0
	// PaddingTypePKCS1 PKCS1.5
	PaddingTypePKCS1 PaddingType = 1
	// PaddingTypeOAEP OAEP模式
	PaddingTypeOAEP PaddingType = 2
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

type CreateSetting struct {
	Length int
}

type RsaKeyManager struct {
	CreateSetting CreateSetting
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
		return nil, errors.New("empty public key")
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
		return nil, errors.New("empty privateKey key")
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
	paddingType     PaddingType
	hashForOAEP     hash.Hash
	labelForOAEP    []byte
	hashForSign     hash.Hash
	hashTypeForSign crypto.Hash
}

func (r *rsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	if r.hashForSign == nil {
		return nil, errors.New("nil hash function")
	}
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, errors.New("empty privateKey key")
	}
	hased := r.hashForSign.Sum(raw)
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), r.hashTypeForSign, hased)
	case PaddingTypeOAEP:
	default:

	}
	return nil, errors.New("not supported paddingType")

}

// Verify 签名验证
func (r *rsaSign) Verify(keyPair KeyPair, sign []byte) ([]byte, error) {
	return nil, nil
}

// NewRsaEncryptWithPaddingPKCS1 创建一个标准的RSA加密 Padding-PKCS1模式实例
func NewRsaEncryptWithPaddingPKCS1() *rsaEncrypt {
	var encrypt rsaEncrypt
	encrypt.paddingType = PaddingTypePKCS1
	return &encrypt
}

// NewRsaEncryptWithPaddingOAEP 创建一个标准的RSA加密 Padding-OAEP模式实例
func NewRsaEncryptWithPaddingOAEP(hash hash.Hash, label []byte) (*rsaEncrypt, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var encrypt rsaEncrypt
	encrypt.paddingType = PaddingTypeOAEP
	encrypt.hashForOAEP = hash
	encrypt.labelForOAEP = label
	return &encrypt, nil
}

func NewRsaSignWithPaddingPKCS1AndSHA256() *rsaSign {
	var sign rsaSign
	sign.paddingType = PaddingTypePKCS1
	sign.hashForSign = sha256.New()
	sign.hashTypeForSign = crypto.SHA256
	return &sign
}
