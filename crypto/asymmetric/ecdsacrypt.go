package asymmetric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"hash"
	"sync"
)

type CreateEcdsaSetting struct {
	Curve elliptic.Curve
}

type EcdsaKeyManager struct {
	CreateSetting CreateEcdsaSetting
}

func (e *EcdsaKeyManager) Create() (KeyPair, error) {
	if e.CreateSetting.Curve == nil {
		return nil, errors.New("nil curve")
	}
	privateKey, err := ecdsa.GenerateKey(e.CreateSetting.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ecdsaKey{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
	}, nil
}

func (e *EcdsaKeyManager) Load() (KeyPair, error) {
	return nil, errors.New("not support now")
}

type ecdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (r *ecdsaKey) PrivateKey() interface{} {
	return r.privateKey
}

func (r *ecdsaKey) PublicKey() interface{} {
	return r.publicKey
}

type ecdsaSign struct {
	sync.Mutex
	hash hash.Hash
}

func (e *ecdsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	defer e.Unlock()
	e.Lock()
	privateKey := keyPair.PrivateKey()
	e.hash.Reset()
	e.hash.Write(raw)
	if privateKey == nil {
		return nil, errors.New("nil privateKey key")
	}
	return nil, nil
	//return ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), e.hash.Sum(nil))
}

func (e *ecdsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	return nil
}

func (e *ecdsaSign) SignBase64(keyPair KeyPair, base64Raw string) (string, error) {
	return "", nil
}

func (e *ecdsaSign) VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error { return nil }
