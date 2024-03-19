package asymmetric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash"
	"math/big"
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

type EcdsaSign struct {
	sync.Mutex
	hash hash.Hash
}

func NewEcdsaSign(hash hash.Hash) *EcdsaSign {
	return &EcdsaSign{
		hash: hash,
	}
}

func (e *EcdsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, errors.New("nil privateKey key")
	}
	if e.hash == nil {
		return ecdsa.SignASN1(rand.Reader, privateKey.(*ecdsa.PrivateKey), raw)
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.SignASN1(rand.Reader, privateKey.(*ecdsa.PrivateKey), raw)
}

func (e *EcdsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return errors.New("nil publicKey key")
	}
	if e.hash == nil {
		flag := ecdsa.VerifyASN1(publicKey.(*ecdsa.PublicKey), raw, sign)
		if flag {
			return nil
		}
		return errors.New("verify failed")
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	flag := ecdsa.VerifyASN1(publicKey.(*ecdsa.PublicKey), raw, sign)
	if flag {
		return nil
	}
	return errors.New("verify failed")
}

func (e *EcdsaSign) SignRS(keyPair KeyPair, raw []byte) (*big.Int, *big.Int, error) {
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, nil, errors.New("nil privateKey key")
	}
	if e.hash == nil {
		return ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), raw)
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), raw)
}

func (e *EcdsaSign) VerifyRS(keyPair KeyPair, raw []byte, r, s *big.Int) (bool, error) {
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return false, errors.New("nil publicKey key")
	}
	if e.hash == nil {
		return ecdsa.Verify(keyPair.PublicKey().(*ecdsa.PublicKey), raw, r, s), nil
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.Verify(keyPair.PublicKey().(*ecdsa.PublicKey), raw, r, s), nil
}

func (e *EcdsaSign) SignBase64(keyPair KeyPair, base64Raw string) (string, error) {
	content, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	result, err := e.Sign(keyPair, content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func (e *EcdsaSign) VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error {
	rawContent, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return err
	}
	signContent, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}
	return e.Verify(keyPair, rawContent, signContent)
}
