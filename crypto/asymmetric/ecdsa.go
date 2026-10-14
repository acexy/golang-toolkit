package asymmetric

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"hash"
	"math/big"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/math/conversion"
)

// NewEcdsaKeyManager creates a key manager using the supplied elliptic curve.
func NewEcdsaKeyManager(curve elliptic.Curve) *EcdsaKeyManager {
	return &EcdsaKeyManager{
		Curve: curve,
	}
}

// NewEmptyEcdsaKeyManager creates a load-only ECDSA key manager without a generation curve.
func NewEmptyEcdsaKeyManager() *EcdsaKeyManager {
	return &EcdsaKeyManager{}
}

// EcdsaKeyManager generates ECDSA keys and loads PEM key material.
type EcdsaKeyManager struct {
	// Curve selects the elliptic curve used for key generation.
	Curve elliptic.Curve
}

// Create generates a new public/private key pair using the manager configuration.
func (e *EcdsaKeyManager) Create() (EcdsaKeyPair, error) {
	if e.Curve == nil {
		return nil, toolkitError.ErrNilCurve
	}
	privateKey, err := ecdsa.GenerateKey(e.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ecdsaKey{
		publicKey:  &privateKey.PublicKey,
		privateKey: privateKey,
	}, nil
}

// LoadPublicKey loads a public key from PEM.
func (e *EcdsaKeyManager) LoadPublicKey(pubPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, "")
}

// LoadPrivateKey loads a private key from PEM and derives its public key.
func (e *EcdsaKeyManager) LoadPrivateKey(priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair("", priPem)
}

// LoadKeyPair loads PEM keys and validates their match when both are supplied.
func (e *EcdsaKeyManager) LoadKeyPair(pubPem, priPem string) (EcdsaKeyPair, error) {
	if pubPem == "" && priPem == "" {
		return nil, toolkitError.ErrBadKey
	}

	var pub *ecdsa.PublicKey
	var pri *ecdsa.PrivateKey
	var err error

	// Parse the public key.
	if pubPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(pubPem))
		if block == nil {
			return nil, toolkitError.ErrBadPublicKey
		}

		switch block.Type {
		case "PUBLIC KEY":
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
			return nil, toolkitError.ErrUnsupportedPublicKeyType
		}
	}

	// Parse the private key.
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
			return nil, toolkitError.ErrUnsupportedPrivateKeyType
		}
	}
	if pub == nil && pri != nil {
		pub = &pri.PublicKey
	}
	if pub != nil && pri != nil && !sameEcdsaPublicKey(pub, &pri.PublicKey) {
		return nil, toolkitError.ErrKeyPairMismatch
	}

	key := &ecdsaKey{
		publicKey:  pub,
		privateKey: pri,
	}
	return key, nil
}

func sameEcdsaPublicKey(pub1, pub2 *ecdsa.PublicKey) bool {
	if pub1 == nil || pub2 == nil {
		return false
	}
	return pub1.Equal(pub2)
}

// Load loads PEM keys through LoadKeyPair.
func (e *EcdsaKeyManager) Load(pubPem, priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, priPem)
}

type ecdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// PrivateKey returns the raw private key, or nil if it is absent.
func (e *ecdsaKey) PrivateKey() any {
	if e.privateKey == nil {
		return nil
	}
	return e.privateKey
}

// PublicKey returns the raw public key, or nil if it is absent.
func (e *ecdsaKey) PublicKey() any {
	if e.publicKey == nil {
		return nil
	}
	return e.publicKey
}

// ToPublicPem exports the public key in PKIX PEM format.
func (e *ecdsaKey) ToPublicPem() (string, error) {
	return e.ToPKIXPublicPem()
}

// ToPrivatePem exports the private key in PKCS8 PEM format.
func (e *ecdsaKey) ToPrivatePem() (string, error) {
	return e.ToPKCS8PrivatePem()
}

// ToECPrivatePem exports a SEC1 EC private-key PEM.
func (e *ecdsaKey) ToECPrivatePem() (string, error) {
	if e.privateKey == nil {
		return "", toolkitError.ErrNilPrivateKey
	}
	der, err := x509.MarshalECPrivateKey(e.privateKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	})), nil
}

// ToPKIXPublicPem exports a PKIX public-key PEM.
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

// ToPKCS8PrivatePem exports a PKCS8 private-key PEM.
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

// EcdsaSign signs and verifies ECDSA data, synchronizing access to its configured hash.
type EcdsaSign struct {
	sync.Mutex
	hashFunc hash.Hash
}

// NewEcdsaSign creates an ECDSA signer; a nil hash requires precomputed digest input.
func NewEcdsaSign(hashFunc hash.Hash) *EcdsaSign {
	return &EcdsaSign{
		hashFunc: hashFunc,
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

// Sign returns an ASN.1 DER signature; input must be a digest when no hash is configured.
func (e *EcdsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	privateKey, err := loadEcdsaPrivateKey(keyPair)
	if err != nil {
		return nil, err
	}
	if e.hashFunc == nil {
		return ecdsa.SignASN1(rand.Reader, privateKey, raw)
	}
	e.Lock()
	e.hashFunc.Reset()
	e.hashFunc.Write(raw)
	raw = e.hashFunc.Sum(nil)
	e.Unlock()
	return ecdsa.SignASN1(rand.Reader, privateKey, raw)
}

// Verify verifies an ASN.1 DER signature; input must be a digest when no hash is configured.
func (e *EcdsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	publicKey, err := loadEcdsaPublicKey(keyPair)
	if err != nil {
		return err
	}
	if e.hashFunc == nil {
		flag := ecdsa.VerifyASN1(publicKey, raw, sign)
		if flag {
			return nil
		}
		return toolkitError.ErrVerifyFailed
	}
	e.Lock()
	e.hashFunc.Reset()
	e.hashFunc.Write(raw)
	raw = e.hashFunc.Sum(nil)
	e.Unlock()
	flag := ecdsa.VerifyASN1(publicKey, raw, sign)
	if flag {
		return nil
	}
	return toolkitError.ErrVerifyFailed
}

// SignRS returns ECDSA r and s values; input must be a digest when no hash is configured.
func (e *EcdsaSign) SignRS(keyPair KeyPair, raw []byte) (*big.Int, *big.Int, error) {
	privateKey, err := loadEcdsaPrivateKey(keyPair)
	if err != nil {
		return nil, nil, err
	}
	if e.hashFunc == nil {
		return ecdsa.Sign(rand.Reader, privateKey, raw)
	}
	e.Lock()
	e.hashFunc.Reset()
	e.hashFunc.Write(raw)
	raw = e.hashFunc.Sum(nil)
	e.Unlock()
	return ecdsa.Sign(rand.Reader, privateKey, raw)
}

// VerifyRS verifies ECDSA r and s; input must be a digest when no hash is configured.
func (e *EcdsaSign) VerifyRS(keyPair KeyPair, raw []byte, r, s *big.Int) (bool, error) {
	publicKey, err := loadEcdsaPublicKey(keyPair)
	if err != nil {
		return false, err
	}
	if e.hashFunc == nil {
		return ecdsa.Verify(publicKey, raw, r, s), nil
	}
	e.Lock()
	e.hashFunc.Reset()
	e.hashFunc.Write(raw)
	raw = e.hashFunc.Sum(nil)
	e.Unlock()
	return ecdsa.Verify(publicKey, raw, r, s), nil
}

// SignBase64 decodes Base64 data and returns its signature as Base64.
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

// VerifyBase64 decodes Base64 data and signature before verification.
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
