package asymmetric

// KeyPairManager KeyPair管理器
type KeyPairManager interface {
	// Create 生成新的公私钥对
	Create() (KeyPair, error)
	// Load 加载公私钥
	Load(pubPem, priPem string) (KeyPair, error)
}

// KeyPair 公私钥信息
type KeyPair interface {
	// PrivateKey 获取原始私钥
	PrivateKey() interface{}
	// PublicKey 获取原始公钥信息
	PublicKey() interface{}
	// ToPublicPKCS1Pem  将公钥转换为pkcs1 PEM格式
	ToPublicPKCS1Pem() string
	// ToPrivatePKCS1Pem 将私钥转换为pkcs1 PEM格式
	ToPrivatePKCS1Pem() string
	// ToPublicPKCS8Pem  将公钥转换为pkcs8 PEM格式
	ToPublicPKCS8Pem() (string, error)
	// ToPrivatePKCS8Pem 将私钥转换为pkcs8 PEM格式
	ToPrivatePKCS8Pem() string
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
