package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
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
	paddingType PaddingType
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
		//return rsa.EncryptOAEP(rand.Reader, a.rsaKey.publicKey, raw)
	default:

	}
	return nil, nil
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
		//return rsa.EncryptOAEP(rand.Reader, a.rsaKey.publicKey, raw)
	default:

	}
	return nil, nil
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

func NewRsaWithPaddingPKCS1() *RsaAsymmetric {
	var rsaAsymmetric RsaAsymmetric
	rsaAsymmetric.paddingType = PaddingTypePKCS1
	return &rsaAsymmetric
}
