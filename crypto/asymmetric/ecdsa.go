package asymmetric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"math/big"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/math/conversion"
)

type CreateEcdsaSetting struct {
	Curve elliptic.Curve
}

func NewEcdsaKeyManager(setting CreateEcdsaSetting) *EcdsaKeyManager {
	return &EcdsaKeyManager{
		CreateSetting: setting,
	}
}

func NewEmptyEcdsaKeyManager() *EcdsaKeyManager {
	return &EcdsaKeyManager{}
}

type EcdsaKeyManager struct {
	CreateSetting CreateEcdsaSetting
}

func (e *EcdsaKeyManager) Create() (EcdsaKeyPair, error) {
	if e.CreateSetting.Curve == nil {
		return nil, toolkitError.ErrNilCurve
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

func (e *EcdsaKeyManager) LoadPublicKey(pubPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, "")
}

func (e *EcdsaKeyManager) LoadPrivateKey(priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair("", priPem)
}

func (e *EcdsaKeyManager) LoadKeyPair(pubPem, priPem string) (EcdsaKeyPair, error) {
	if pubPem == "" && priPem == "" {
		return nil, toolkitError.ErrBadKey
	}

	var pub *ecdsa.PublicKey
	var pri *ecdsa.PrivateKey
	var err error

	// 解析公钥
	if pubPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(pubPem))
		if block == nil {
			return nil, toolkitError.ErrBadPublicKey
		}

		switch block.Type {
		case "PUBLIC KEY", "EC PUBLIC KEY":
			var iface any
			iface, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pub, ok = iface.(*ecdsa.PublicKey)
			if !ok {
				return nil, toolkitError.ErrNotEcdsaPublicKey
			}
		default:
			return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
		}
	}

	// 解析私钥
	if priPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(priPem))
		if block == nil {
			return nil, toolkitError.ErrBadPrivateKey
		}

		switch block.Type {
		case "EC PRIVATE KEY":
			pri, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PRIVATE KEY":
			var iface any
			iface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pri, ok = iface.(*ecdsa.PrivateKey)
			if !ok {
				return nil, toolkitError.ErrNotEcdsaPrivateKey
			}
		default:
			return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
		}
	}
	if pub == nil && pri != nil {
		pub = &pri.PublicKey
	}

	key := &ecdsaKey{
		publicKey:  pub,
		privateKey: pri,
	}
	return key, nil
}

func (e *EcdsaKeyManager) Load(pubPem, priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, priPem)
}

type ecdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (e *ecdsaKey) PrivateKey() any {
	return e.privateKey
}

func (e *ecdsaKey) PublicKey() any {
	return e.publicKey
}

func (e *ecdsaKey) ToPublicPem() (string, error) {
	return e.ToPKIXPublicPem()
}

func (e *ecdsaKey) ToPrivatePem() (string, error) {
	return e.ToPKCS8PrivatePem()
}

func (e *ecdsaKey) ToPKIXPublicPem() (string, error) {
	if e.publicKey == nil {
		return "", toolkitError.ErrNilPublicKey
	}
	der, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})), nil
}

func (e *ecdsaKey) ToPKCS8PrivatePem() (string, error) {
	if e.privateKey == nil {
		return "", toolkitError.ErrNilPrivateKey
	}
	der, err := x509.MarshalPKCS8PrivateKey(e.privateKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})), nil
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

func loadEcdsaPublicKey(keyPair KeyPair) (*ecdsa.PublicKey, error) {
	if keyPair == nil {
		return nil, toolkitError.ErrNilKeyPair
	}
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return nil, toolkitError.ErrNilPublicKey
	}
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, toolkitError.ErrNotEcdsaPublicKey
	}
	return ecdsaPublicKey, nil
}

func loadEcdsaPrivateKey(keyPair KeyPair) (*ecdsa.PrivateKey, error) {
	if keyPair == nil {
		return nil, toolkitError.ErrNilKeyPair
	}
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, toolkitError.ErrNilPrivateKey
	}
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, toolkitError.ErrNotEcdsaPrivateKey
	}
	return ecdsaPrivateKey, nil
}

func (e *EcdsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	privateKey, err := loadEcdsaPrivateKey(keyPair)
	if err != nil {
		return nil, err
	}
	if e.hash == nil {
		return ecdsa.SignASN1(rand.Reader, privateKey, raw)
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.SignASN1(rand.Reader, privateKey, raw)
}

func (e *EcdsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	publicKey, err := loadEcdsaPublicKey(keyPair)
	if err != nil {
		return err
	}
	if e.hash == nil {
		flag := ecdsa.VerifyASN1(publicKey, raw, sign)
		if flag {
			return nil
		}
		return toolkitError.ErrVerifyFailed
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	flag := ecdsa.VerifyASN1(publicKey, raw, sign)
	if flag {
		return nil
	}
	return toolkitError.ErrVerifyFailed
}

func (e *EcdsaSign) SignRS(keyPair KeyPair, raw []byte) (*big.Int, *big.Int, error) {
	privateKey, err := loadEcdsaPrivateKey(keyPair)
	if err != nil {
		return nil, nil, err
	}
	if e.hash == nil {
		return ecdsa.Sign(rand.Reader, privateKey, raw)
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.Sign(rand.Reader, privateKey, raw)
}

func (e *EcdsaSign) VerifyRS(keyPair KeyPair, raw []byte, r, s *big.Int) (bool, error) {
	publicKey, err := loadEcdsaPublicKey(keyPair)
	if err != nil {
		return false, err
	}
	if e.hash == nil {
		return ecdsa.Verify(publicKey, raw, r, s), nil
	}
	e.Lock()
	e.hash.Reset()
	e.hash.Write(raw)
	raw = e.hash.Sum(nil)
	e.Unlock()
	return ecdsa.Verify(publicKey, raw, r, s), nil
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
