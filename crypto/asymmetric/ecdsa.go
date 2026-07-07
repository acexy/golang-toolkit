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

// NewEcdsaKeyManager 创建带指定曲线配置的 ECDSA KeyManager
func NewEcdsaKeyManager(curve elliptic.Curve) *EcdsaKeyManager {
	return &EcdsaKeyManager{
		Curve: curve,
	}
}

// NewEmptyEcdsaKeyManager 创建空 ECDSA KeyManager，仅适合加载已有 PEM 密钥
func NewEmptyEcdsaKeyManager() *EcdsaKeyManager {
	return &EcdsaKeyManager{}
}

// EcdsaKeyManager 管理 ECDSA 密钥创建和 PEM 加载
type EcdsaKeyManager struct {
	// Curve 表示创建 ECDSA 密钥时使用的椭圆曲线
	Curve elliptic.Curve
}

// Create 创建新的 ECDSA 公私钥对
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

// LoadPublicKey 从 PEM 字符串加载 ECDSA 公钥
func (e *EcdsaKeyManager) LoadPublicKey(pubPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, "")
}

// LoadPrivateKey 从 PEM 字符串加载 ECDSA 私钥，并派生对应公钥
func (e *EcdsaKeyManager) LoadPrivateKey(priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair("", priPem)
}

// LoadKeyPair 从 PEM 字符串加载 ECDSA 公私钥，并校验二者是否匹配
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
			return nil, toolkitError.NewUnsupportedPublicKeyTypeError(block.Type)
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
			return nil, toolkitError.NewUnsupportedPrivateKeyTypeError(block.Type)
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
	if pub1 == nil || pub2 == nil || pub1.Curve == nil || pub2.Curve == nil || pub1.X == nil || pub1.Y == nil || pub2.X == nil || pub2.Y == nil {
		return false
	}
	return pub1.Curve == pub2.Curve && pub1.X.Cmp(pub2.X) == 0 && pub1.Y.Cmp(pub2.Y) == 0
}

// Load 从 PEM 字符串加载 ECDSA 公私钥
func (e *EcdsaKeyManager) Load(pubPem, priPem string) (EcdsaKeyPair, error) {
	return e.LoadKeyPair(pubPem, priPem)
}

type ecdsaKey struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// PrivateKey 返回原始 ECDSA 私钥
func (e *ecdsaKey) PrivateKey() any {
	if e.privateKey == nil {
		return nil
	}
	return e.privateKey
}

// PublicKey 返回原始 ECDSA 公钥
func (e *ecdsaKey) PublicKey() any {
	if e.publicKey == nil {
		return nil
	}
	return e.publicKey
}

// ToPublicPem 将 ECDSA 公钥导出为默认 PKIX PEM 格式
func (e *ecdsaKey) ToPublicPem() (string, error) {
	return e.ToPKIXPublicPem()
}

// ToPrivatePem 将 ECDSA 私钥导出为默认 PKCS8 PEM 格式
func (e *ecdsaKey) ToPrivatePem() (string, error) {
	return e.ToPKCS8PrivatePem()
}

// ToECPrivatePem 将 ECDSA 私钥导出为 SEC1 EC PRIVATE KEY PEM 格式
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

// ToPKIXPublicPem 将 ECDSA 公钥导出为 PKIX PEM 格式
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

// ToPKCS8PrivatePem 将 ECDSA 私钥导出为 PKCS8 PEM 格式
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

// EcdsaSign 提供 ECDSA 签名和验签能力
type EcdsaSign struct {
	sync.Mutex
	hashFunc hash.Hash
}

// NewEcdsaSign 创建 ECDSA 签名验签实例，hashFunc 为空时调用方必须传入摘要数据
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

// Sign 使用 ECDSA 私钥对数据进行 ASN.1 DER 签名，未配置 hashFunc 时 raw 必须是摘要
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

// Verify 使用 ECDSA 公钥验证 ASN.1 DER 签名，未配置 hashFunc 时 raw 必须是摘要
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

// SignRS 使用 ECDSA 私钥签名并返回 r、s 两段签名值，未配置 hashFunc 时 raw 必须是摘要
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

// VerifyRS 使用 ECDSA 公钥验证 r、s 两段签名值，未配置 hashFunc 时 raw 必须是摘要
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

// SignBase64 解码 Base64 原文后签名，并返回 Base64 签名
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

// VerifyBase64 解码 Base64 原文和签名后执行验签
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
