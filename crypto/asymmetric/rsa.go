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

type PaddingType uint8

const (
	// PaddingTypePKCS1 表示 RSA PKCS#1 v1.5 填充，可用于加解密和签名验签
	PaddingTypePKCS1 PaddingType = 1
	// PaddingTypeOAEP 表示 RSA OAEP 填充，仅用于加解密
	PaddingTypeOAEP PaddingType = 2
	// PaddingTypePSS 表示 RSA PSS 填充，仅用于签名验签
	PaddingTypePSS PaddingType = 3
)

type rsaKey struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *rsaKey) PrivateKey() any {
	return r.privateKey
}

func (r *rsaKey) PublicKey() any {
	return r.publicKey
}

func (r *rsaKey) ToPublicPem() (string, error) {
	return r.ToPKIXPublicPem()
}

func (r *rsaKey) ToPrivatePem() (string, error) {
	return r.ToPKCS8PrivatePem()
}

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

type CreateRsaSetting struct {
	Length int
}

// NewRsaKeyManager 创建带默认密钥长度的 RSA KeyManager
func NewRsaKeyManager(keyLength int) *RsaKeyManager {
	return &RsaKeyManager{
		CreateSetting: CreateRsaSetting{
			Length: keyLength,
		},
	}
}

// NewEmptyRsaKeyManager 创建空 RSA KeyManager，仅适合加载已有 PEM 密钥
func NewEmptyRsaKeyManager() *RsaKeyManager {
	return &RsaKeyManager{}
}

type RsaKeyManager struct {
	CreateSetting CreateRsaSetting
}

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

func (r *RsaKeyManager) LoadPublicKey(pubPem string) (RsaKeyPair, error) {
	return r.LoadKeyPair(pubPem, "")
}

func (r *RsaKeyManager) LoadPrivateKey(priPem string) (RsaKeyPair, error) {
	return r.LoadKeyPair("", priPem)
}

func (r *RsaKeyManager) LoadKeyPair(pubPem, priPem string) (RsaKeyPair, error) {
	if pubPem == "" && priPem == "" {
		return nil, toolkitError.ErrBadKey
	}

	var pub *rsa.PublicKey
	var pri *rsa.PrivateKey
	var err error

	// 解析公钥
	if pubPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(pubPem))
		if block == nil {
			return nil, toolkitError.ErrBadPublicKey
		}

		switch block.Type {
		case "RSA PUBLIC KEY":
			// PKCS#1 公钥
			pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PUBLIC KEY":
			// PKIX (PKCS#8) 公钥
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

	// 解析私钥
	if priPem != "" {
		block, _ := pem.Decode(conversion.ParseBytes(priPem))
		if block == nil {
			return nil, toolkitError.ErrBadPrivateKey
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			// PKCS#1 私钥
			pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
		case "PRIVATE KEY":
			// PKCS#8 私钥
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

type RsaEncrypt struct {
	paddingType  PaddingType
	hashForOAEP  hash.Hash
	labelForOAEP []byte
}

func (r *RsaEncrypt) Encrypt(keyPair KeyPair, raw []byte) ([]byte, error) {
	publicKey, err := loadRsaPublicKey(keyPair)
	if err != nil {
		return nil, err
	}
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.EncryptPKCS1v15(rand.Reader, publicKey, raw)
	case PaddingTypeOAEP:
		if r.hashForOAEP == nil {
			return nil, toolkitError.ErrNilHashFunction
		}
		return rsa.EncryptOAEP(r.hashForOAEP, rand.Reader, publicKey, raw, r.labelForOAEP)
	default:

	}
	return nil, toolkitError.ErrUnsupportedPaddingType
}

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

func (r *RsaEncrypt) Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error) {
	privateKey, err := loadRsaPrivateKey(keyPair)
	if err != nil {
		return nil, err
	}
	switch r.paddingType {
	case PaddingTypePKCS1:
		return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipher)
	case PaddingTypeOAEP:
		if r.hashForOAEP == nil {
			return nil, toolkitError.ErrNilHashFunction
		}
		return rsa.DecryptOAEP(r.hashForOAEP, rand.Reader, privateKey, cipher, r.labelForOAEP)
	default:

	}
	return nil, toolkitError.ErrUnsupportedPaddingType
}

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

type RsaSign struct {
	sync.Mutex
	paddingType     PaddingType
	hashForSign     hash.Hash
	hashTypeForSign crypto.Hash
	options         *rsa.PSSOptions
}

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

// Verify 签名验证
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

// NewRsaEncryptWithPKCS1 创建 RSA PKCS#1 v1.5 加解密实例
func NewRsaEncryptWithPKCS1() *RsaEncrypt {
	var encrypt RsaEncrypt
	encrypt.paddingType = PaddingTypePKCS1
	return &encrypt
}

// NewRsaEncryptWithOAEP 创建 RSA OAEP 加解密实例，label 在加密和解密时必须一致
func NewRsaEncryptWithOAEP(hashFunc hash.Hash, label []byte) (*RsaEncrypt, error) {
	if hashFunc == nil {
		return nil, toolkitError.ErrNilHashFunction
	}
	var encrypt RsaEncrypt
	encrypt.paddingType = PaddingTypeOAEP
	encrypt.hashForOAEP = hashFunc
	encrypt.labelForOAEP = label
	return &encrypt, nil
}

// NewRsaSignWithPKCS1AndSHA256 创建 RSA PKCS#1 v1.5 SHA256 签名验签实例
func NewRsaSignWithPKCS1AndSHA256() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha256.New(), crypto.SHA256)
	return sign
}

// NewRsaSignWithPKCS1AndSHA512 创建 RSA PKCS#1 v1.5 SHA512 签名验签实例
func NewRsaSignWithPKCS1AndSHA512() *RsaSign {
	sign, _ := NewRsaSignWithPKCS1(sha512.New(), crypto.SHA512)
	return sign
}

// NewRsaSignWithPKCS1 创建 RSA PKCS#1 v1.5 自定义 hash 签名验签实例，hashType 必须与 hashFunc 匹配
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

// NewRsaSignWithPSSAndSHA256 创建 RSA PSS SHA256 签名验签实例，可选指定 saltLength
func NewRsaSignWithPSSAndSHA256(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOptions(sha256.New(), crypto.SHA256, length)
	return sign
}

// NewRsaSignWithPSSAndSHA512 创建 RSA PSS SHA512 签名验签实例，可选指定 saltLength
func NewRsaSignWithPSSAndSHA512(saltLength ...int) *RsaSign {
	length := -1
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	sign, _ := NewRsaSignWithPSSAndOptions(sha512.New(), crypto.SHA512, length)
	return sign
}

// NewRsaSignWithPSS 创建 RSA PSS 自定义 hash 签名验签实例，默认使用 rsa.PSSSaltLengthAuto
func NewRsaSignWithPSS(hashFunc hash.Hash, hashType crypto.Hash) (*RsaSign, error) {
	return NewRsaSignWithPSSAndOptions(hashFunc, hashType, -1)
}

// NewRsaSignWithPSSAndOptions 创建 RSA PSS 自定义 hash 签名验签实例，并指定 saltLength
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
