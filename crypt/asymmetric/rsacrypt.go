package asymmetric

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/acexy/golang-toolkit/logger"
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
)

type RsaKey struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (k *RsaKey) PrivateKey() interface{} {
	return k.privateKey
}

func (k *RsaKey) PublicKey() interface{} {
	return k.publicKey
}

type RsaAsymmetric struct {
	length      int
	paddingType PaddingType
	keyPair     KeyPair
	sync.Once
}

func NewRsaWithPaddingPKCS1(length int) *RsaAsymmetric {
	var rsaAsymmetric RsaAsymmetric
	rsaAsymmetric.length = length
	rsaAsymmetric.paddingType = PaddingTypePKCS1
	return &rsaAsymmetric
}

func (a *RsaAsymmetric) create() KeyPair {
	a.Once.Do(func() {
		privateKey, err := rsa.GenerateKey(rand.Reader, a.length)
		if err != nil {
			logger.Logrus().Error("generate rsa key error", err)
		}
		a.keyPair = &RsaKey{
			publicKey:  &privateKey.PublicKey,
			privateKey: privateKey,
		}
	})
	pair := a.keyPair
	return pair
}

func (a *RsaAsymmetric) encrypt(content []byte) ([]byte, error) {
	a.create()
	switch a.paddingType {
	case PaddingTypePKCS1:
	case PaddingTypeOAEP:
		//return rsa.EncryptOAEP(rand.Reader, a.rsaKey.publicKey, content)
	default:

	}
	return nil, nil
}
