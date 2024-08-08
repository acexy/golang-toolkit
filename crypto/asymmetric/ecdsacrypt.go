package asymmetric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/acexy/golang-toolkit/math/conversion"
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

func (e *EcdsaKeyManager) Load(pubPem, priPem string) (KeyPair, error) {
	block, _ := pem.Decode(conversion.ParseBytes(pubPem))
	if block == nil {
		return nil, errors.New("bak public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecdsaPubKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("bak public key")
	}
	block, _ = pem.Decode(conversion.ParseBytes(priPem))
	if block == nil {
		return nil, errors.New("bak private key")
	}
	pri, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key := &ecdsaKey{
		publicKey:  ecdsaPubKey,
		privateKey: pri,
	}
	return key, nil
}

type ecdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (e *ecdsaKey) PrivateKey() interface{} {
	return e.privateKey
}

func (e *ecdsaKey) PublicKey() interface{} {
	return e.publicKey
}

func (e *ecdsaKey) ToPublicPKCS1Pem() string {
	publicKey := e.PublicKey().(*ecdsa.PublicKey)
	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return ""
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: der,
	}))
}

func (e *ecdsaKey) ToPrivatePKCS1Pem() string {
	privateKey := e.PrivateKey().(*ecdsa.PrivateKey)
	der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return ""
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}))
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
