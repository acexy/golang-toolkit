package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
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

type RsaKey struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *RsaKey) PrivateKey() interface{} {
	return r.privateKey
}

func (r *RsaKey) PublicKey() interface{} {
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
	return &RsaKey{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
	}, nil
}

func (r *RsaKeyManager) Load() (KeyPair, error) {
	return nil, errors.New("not support now")
}

type RsaAsymmetric struct {
	paddingType  PaddingType
	hashForOAEP  hash.Hash
	labelForOAEP []byte
}

func (a *RsaAsymmetric) Encrypt(keyPair KeyPair, raw []byte) ([]byte, error) {
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return nil, errors.New("empty public key")
	}
	switch a.paddingType {
	case PaddingTypePKCS1:
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), raw)
	case PaddingTypeOAEP:
		return rsa.EncryptOAEP(a.hashForOAEP, rand.Reader, publicKey.(*rsa.PublicKey), raw, a.labelForOAEP)
	default:

	}
	return nil, errors.New("not supported paddingType")
}

func (a *RsaAsymmetric) EncryptBase64(keyPair KeyPair, base64Raw string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	result, err := a.Encrypt(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func (a *RsaAsymmetric) Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error) {
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, errors.New("empty privateKey key")
	}
	switch a.paddingType {
	case PaddingTypePKCS1:
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), cipher)
	case PaddingTypeOAEP:
		return rsa.DecryptOAEP(a.hashForOAEP, rand.Reader, privateKey.(*rsa.PrivateKey), cipher, a.labelForOAEP)
	default:

	}
	return nil, errors.New("not supported paddingType")
}

func (a *RsaAsymmetric) DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return "", err
	}
	result, err := a.Decrypt(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

// NewRsaWithPaddingPKCS1 创建一个标准的RSA Padding-PKCS1模式实例
func NewRsaWithPaddingPKCS1() *RsaAsymmetric {
	var rsaAsymmetric RsaAsymmetric
	rsaAsymmetric.paddingType = PaddingTypePKCS1
	return &rsaAsymmetric
}

func NewRsaWithPaddingOAEP(hash hash.Hash, label []byte) (*RsaAsymmetric, error) {
	if hash == nil {
		return nil, errors.New("nil hash function")
	}
	var rsaAsymmetric RsaAsymmetric
	rsaAsymmetric.paddingType = PaddingTypeOAEP
	rsaAsymmetric.hashForOAEP = hash
	return &rsaAsymmetric, nil
}
