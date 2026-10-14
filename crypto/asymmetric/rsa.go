package asymmetric

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/math/conversion"
)

// PaddingType identifies the RSA signature padding scheme.
type PaddingType uint8

const (
	// PaddingTypePKCS1 selects PKCS1 v1.5 signature padding.
	PaddingTypePKCS1 PaddingType = 1
	// PaddingTypePSS selects RSA PSS signature padding.
	PaddingTypePSS PaddingType = 3
)

type rsaKey struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// PrivateKey returns the raw private key, or nil if it is absent.
func (r *rsaKey) PrivateKey() any {
	if r.privateKey == nil {
		return nil
	}
	return r.privateKey
}

// PublicKey returns the raw public key, or nil if it is absent.
func (r *rsaKey) PublicKey() any {
	if r.publicKey == nil {
		return nil
	}
	return r.publicKey
}

// ToPublicPem exports the public key in PKIX PEM format.
func (r *rsaKey) ToPublicPem() (string, error) {
	return r.ToPKIXPublicPem()
}

// ToPrivatePem exports the private key in PKCS8 PEM format.
func (r *rsaKey) ToPrivatePem() (string, error) {
	return r.ToPKCS8PrivatePem()
}

// ToPKCS1PublicPem exports an RSA PKCS1 public-key PEM.
func (r *rsaKey) ToPKCS1PublicPem() (string, error) {
	if r.publicKey == nil {
		return "", toolkitError.ErrNilPublicKey
	}
	der := x509.MarshalPKCS1PublicKey(r.publicKey)
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: der,
	})), nil
}

// ToPKCS1PrivatePem exports an RSA PKCS1 private-key PEM.
func (r *rsaKey) ToPKCS1PrivatePem() (string, error) {
	if r.privateKey == nil {
		return "", toolkitError.ErrNilPrivateKey
	}
	der := x509.MarshalPKCS1PrivateKey(r.privateKey)
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	})), nil
}

// ToPKIXPublicPem exports a PKIX public-key PEM.
func (r *rsaKey) ToPKIXPublicPem() (string, error) {
	if r.publicKey == nil {
		return "", toolkitError.ErrNilPublicKey
	}
	der, err := x509.MarshalPKIXPublicKey(r.publicKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})), nil
}

// ToPKCS8PrivatePem exports a PKCS8 private-key PEM.
func (r *rsaKey) ToPKCS8PrivatePem() (string, error) {
	if r.privateKey == nil {
		return "", toolkitError.ErrNilPrivateKey
	}
	der, err := x509.MarshalPKCS8PrivateKey(r.privateKey)
	if err != nil {
		return "", err
	}
	return conversion.FromBytes(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})), nil
}

// CreateRsaSetting configures RSA key generation.
type CreateRsaSetting struct {
	// Length is the RSA key length in bits.
	Length int
}

// NewRsaKeyManager creates a key manager with the supplied RSA key length.
func NewRsaKeyManager(keyLength int) *RsaKeyManager {
	return &RsaKeyManager{
		CreateSetting: CreateRsaSetting{
			Length: keyLength,
		},
	}
}

// NewEmptyRsaKeyManager creates a load-only RSA key manager without a generation length.
func NewEmptyRsaKeyManager() *RsaKeyManager {
	return &RsaKeyManager{}
}

// RsaKeyManager generates RSA keys and loads PEM key material.
type RsaKeyManager struct {
	// CreateSetting holds the default key generation settings.
	CreateSetting CreateRsaSetting
}

// Create generates a new public/private key pair using the manager configuration.
func (r *RsaKeyManager) Create() (RsaKeyPair, error) {
	if r.CreateSetting.Length == 0 {
		return nil, toolkitError.ErrBadKeyLength
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

// LoadPublicKey loads a public key from PEM.
func (r *RsaKeyManager) LoadPublicKey(pubPem string) (RsaKeyPair, error) {
	return r.LoadKeyPair(pubPem, "")
}

// LoadPrivateKey loads a private key from PEM and derives its public key.
func (r *RsaKeyManager) LoadPrivateKey(priPem string) (RsaKeyPair, error) {
	return r.LoadKeyPair("", priPem)
}

// LoadKeyPair loads PEM keys and validates their match when both are supplied.
func (r *RsaKeyManager) LoadKeyPair(pubPem, priPem string) (RsaKeyPair, error) {
	if pubPem == "" && priPem == "" {
		return nil, toolkitError.ErrBadKey
	}

	var pub *rsa.PublicKey
	var pri *rsa.PrivateKey
	var err error

	// Parse the public key.
	if pubPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(pubPem))
		if block == nil {
			return nil, toolkitError.ErrBadPublicKey
		}

		switch block.Type {
		case "RSA PUBLIC KEY":
			// PKCS1 public key.
			pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PUBLIC KEY":
			// PKIX public key.
			var iface any
			iface, err = x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pub, ok = iface.(*rsa.PublicKey)
			if !ok {
				return nil, toolkitError.ErrNotRsaPublicKey
			}
		default:
			return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
		}
	}

	// Parse the private key.
	if priPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(priPem))
		if block == nil {
			return nil, toolkitError.ErrBadPrivateKey
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			// PKCS1 private key.
			pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PRIVATE KEY":
			// PKCS8 private key.
			var iface any
			iface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			var ok bool
			pri, ok = iface.(*rsa.PrivateKey)
			if !ok {
				return nil, toolkitError.ErrNotRsaPrivateKey
			}
		default:
			return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
		}
	}
	if pub == nil && pri != nil {
		pub = &pri.PublicKey
	}
	if pub != nil && pri != nil && (pub.E != pri.PublicKey.E || pub.N.Cmp(pri.PublicKey.N) != 0) {
		return nil, toolkitError.ErrKeyPairMismatch
	}

	return &rsaKey{
		publicKey:  pub,
		privateKey: pri,
	}, nil
}

// Load loads PEM keys through LoadKeyPair.
func (r *RsaKeyManager) Load(pubPem, priPem string) (RsaKeyPair, error) {
	return r.LoadKeyPair(pubPem, priPem)
}

func loadRsaPublicKey(keyPair KeyPair) (*rsa.PublicKey, error) {
	if keyPair == nil {
		return nil, toolkitError.ErrNilKeyPair
	}
	publicKey := keyPair.PublicKey()
	if publicKey == nil {
		return nil, toolkitError.ErrNilPublicKey
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, toolkitError.ErrNotRsaPublicKey
	}
	return rsaPublicKey, nil
}

func loadRsaPrivateKey(keyPair KeyPair) (*rsa.PrivateKey, error) {
	if keyPair == nil {
		return nil, toolkitError.ErrNilKeyPair
	}
	privateKey := keyPair.PrivateKey()
	if privateKey == nil {
		return nil, toolkitError.ErrNilPrivateKey
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, toolkitError.ErrNotRsaPrivateKey
	}
	return rsaPrivateKey, nil
}

// RsaEncrypt encrypts and decrypts with RSA OAEP using the configured hash and label.
type RsaEncrypt struct {
	hashForOAEP  hash.Hash
	labelForOAEP []byte
}

// Encrypt encrypts plaintext bytes with the public key.
func (r *RsaEncrypt) Encrypt(keyPair KeyPair, raw []byte) ([]byte, error) {
	publicKey, err := loadRsaPublicKey(keyPair)
	if err != nil {
		return nil, err
	}
	if r.hashForOAEP == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	return rsa.EncryptOAEP(r.hashForOAEP, rand.Reader, publicKey, raw, r.labelForOAEP)
}

// EncryptBase64 decodes Base64 plaintext, encrypts it, and returns Base64 ciphertext.
func (r *RsaEncrypt) EncryptBase64(keyPair KeyPair, base64Raw string) (string, error) {
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

// Decrypt decrypts ciphertext bytes with the private key.
func (r *RsaEncrypt) Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error) {
	privateKey, err := loadRsaPrivateKey(keyPair)
	if err != nil {
		return nil, err
	}
	if r.hashForOAEP == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	return rsa.DecryptOAEP(r.hashForOAEP, rand.Reader, privateKey, cipher, r.labelForOAEP)
}

// DecryptBase64 decodes Base64 ciphertext, decrypts it, and returns Base64 plaintext.
func (r *RsaEncrypt) DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error) {
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

// RsaSign signs and verifies RSA data, synchronizing access to its configured hash.
type RsaSign struct {
	sync.Mutex
	paddingType     PaddingType
	hashForSign     hash.Hash
	hashTypeForSign crypto.Hash
	options         *rsa.PSSOptions
}

// Sign signs data with the private key.
func (r *RsaSign) Sign(keyPair KeyPair, raw []byte) ([]byte, error) {
	if r.hashForSign == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	privateKey, err := loadRsaPrivateKey(keyPair)
	if err != nil {
		return nil, err
	}
	r.Lock()
	r.hashForSign.Reset()
	r.hashForSign.Write(raw)
	hashSum := r.hashForSign.Sum(nil)
	r.Unlock()
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.SignPKCS1v15(rand.Reader, privateKey, r.hashTypeForSign, hashSum)
	case PaddingTypePSS:
		return rsa.SignPSS(rand.Reader, privateKey, r.hashTypeForSign, hashSum, r.options)
	default:

	}
	return nil, toolkitError.ErrUnsupportedPaddingType
}

// Verify verifies data with the public key, returning an error on failure.
func (r *RsaSign) Verify(keyPair KeyPair, raw, sign []byte) error {
	if r.hashForSign == nil {
		return toolkitError.ErrNilHashFunction
	}
	publicKey, err := loadRsaPublicKey(keyPair)
	if err != nil {
		return err
	}
	r.Lock()
	r.hashForSign.Reset()
	r.hashForSign.Write(raw)
	hashSum := r.hashForSign.Sum(nil)
	r.Unlock()
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.VerifyPKCS1v15(publicKey, r.hashTypeForSign, hashSum, sign)
	case PaddingTypePSS:
		return rsa.VerifyPSS(publicKey, r.hashTypeForSign, hashSum, sign, r.options)
	default:

	}
	return toolkitError.ErrUnsupportedPaddingType
}

// SignBase64 decodes Base64 data and returns its signature as Base64.
func (r *RsaSign) SignBase64(keyPair KeyPair, base64Raw string) (string, error) {
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

// VerifyBase64 decodes Base64 data and signature before verification.
func (r *RsaSign) VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error {
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

// NewRsaEncryptWithOAEP creates an OAEP instance; hash must be nonnil and the label must match on both sides.
func NewRsaEncryptWithOAEP(hashFunc hash.Hash, label []byte) (*RsaEncrypt, error) {
	if hashFunc == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	var encrypt RsaEncrypt
	encrypt.hashForOAEP = hashFunc
	encrypt.labelForOAEP = label
	return &encrypt, nil
}

// NewRsaSignWithPKCS1AndSHA256 creates a PKCS1 v1.5 signer using SHA256.
func NewRsaSignWithPKCS1AndSHA256() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha256.New(), crypto.SHA256)
	return sign
}

// NewRsaSignWithPKCS1AndSHA512 creates a PKCS1 v1.5 signer using SHA512.
func NewRsaSignWithPKCS1AndSHA512() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha512.New(), crypto.SHA512)
	return sign
}

// NewRsaSignWithPKCS1 creates a PKCS1 v1.5 signer; hashType must match hashFunc.
func NewRsaSignWithPKCS1(hashFunc hash.Hash, hashType crypto.Hash) (*RsaSign, error) {
	if hashFunc == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	var sign RsaSign
	sign.paddingType = PaddingTypePKCS1
	sign.hashForSign = hashFunc
	sign.hashTypeForSign = hashType
	return &sign, nil
}

// NewRsaSignWithPSSAndSHA256 creates a SHA256 PSS signer with an optional salt length.
func NewRsaSignWithPSSAndSHA256(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOptions(sha256.New(), crypto.SHA256, length)
	return sign
}

// NewRsaSignWithPSSAndSHA512 creates a SHA512 PSS signer with an optional salt length.
func NewRsaSignWithPSSAndSHA512(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOptions(sha512.New(), crypto.SHA512, length)
	return sign
}

// NewRsaSignWithPSS creates a PSS signer using the automatic salt-length policy.
func NewRsaSignWithPSS(hashFunc hash.Hash, hashType crypto.Hash) (*RsaSign, error) {
	return NewRsaSignWithPSSAndOptions(hashFunc, hashType, -1)
}

// NewRsaSignWithPSSAndOptions creates a PSS signer with the supplied hash and salt-length policy.
func NewRsaSignWithPSSAndOptions(hashFunc hash.Hash, hashType crypto.Hash, saltLength int) (*RsaSign, error) {
	if hashFunc == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	var sign RsaSign
	sign.paddingType = PaddingTypePSS
	sign.hashForSign = hashFunc
	sign.hashTypeForSign = hashType
	if saltLength >= 0 {
		sign.options = &rsa.PSSOptions{
			Hash:       hashType,
			SaltLength: saltLength,
		}
	}
	return &sign, nil
}
