package asymmetric

// KeyPairManager 通用 KeyPair 管理器
type KeyPairManager[T KeyPair] interface {
	// Create 生成新的公私钥对
	Create() (T, error)
	// LoadPublicKey 加载公钥
	LoadPublicKey(pubPem string) (T, error)
	// LoadPrivateKey 加载私钥
	LoadPrivateKey(priPem string) (T, error)
	// LoadKeyPair 加载公私钥
	LoadKeyPair(pubPem, priPem string) (T, error)
}

// KeyPair 公私钥信息
type KeyPair interface {
	// PrivateKey 获取原始私钥
	PrivateKey() any
	// PublicKey 获取原始公钥信息
	PublicKey() any
	// ToPublicPem 将公钥转换为默认 PEM 格式
	ToPublicPem() (string, error)
	// ToPrivatePem 将私钥转换为默认 PEM 格式
	ToPrivatePem() (string, error)
	// ToPKIXPublicPem 将公钥转换为 PKIX PEM 格式
	ToPKIXPublicPem() (string, error)
	// ToPKCS8PrivatePem 将私钥转换为 PKCS8 PEM 格式
	ToPKCS8PrivatePem() (string, error)
}

// RsaKeyPair RSA 专用公私钥信息
type RsaKeyPair interface {
	KeyPair
	// ToPKCS1PublicPem 将 RSA 公钥转换为 PKCS1 PEM 格式
	ToPKCS1PublicPem() (string, error)
	// ToPKCS1PrivatePem 将 RSA 私钥转换为 PKCS1 PEM 格式
	ToPKCS1PrivatePem() (string, error)
}

// RsaKeyPairManager RSA 专用 KeyPair 管理器
type RsaKeyPairManager interface {
	KeyPairManager[RsaKeyPair]
}

// EcdsaKeyPair ECDSA 专用公私钥信息
type EcdsaKeyPair interface {
	KeyPair
}

// EcdsaKeyPairManager ECDSA 专用 KeyPair 管理器
type EcdsaKeyPairManager interface {
	KeyPairManager[EcdsaKeyPair]
}

type CryptEncrypt interface {
	// Encrypt 加密
	Encrypt(keyPair KeyPair, raw []byte) ([]byte, error)
	// EncryptBase64 使用标准Base64传递的数据进行加密
	EncryptBase64(keyPair KeyPair, base64Raw string) (string, error)
	// Decrypt 解密
	Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error)
	// DecryptBase64 使用标准Base64传递的数据进行解密
	DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error)
}

type CryptSign interface {
	// Sign 数据签名
	Sign(keyPair KeyPair, raw []byte) ([]byte, error)
	// Verify 签名验证
	Verify(keyPair KeyPair, raw, sign []byte) error
	// SignBase64 使用标准Base64传递的数据进行加签
	SignBase64(keyPair KeyPair, base64Raw string) (string, error)
	// VerifyBase64 使用标准Base64传递的数据进行验签
	VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error
}
