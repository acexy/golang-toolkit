package asymmetric

// KeyPair 公私钥信息
type KeyPair interface {
	PrivateKey() interface{}
	PublicKey() interface{}
}

type CryptAsymmetric interface {

	// Create 创建原始公私钥
	create() KeyPair

	// Encrypt 加密数据
	encrypt(content []byte) ([]byte, error)
}
