package symmetric

// CryptEncrypt 对称加密通用接口
type CryptEncrypt interface {
	// Encrypt 加密原始字节
	Encrypt(rawData []byte) ([]byte, error)
	// EncryptBase64 加密原始字节并返回 Base64 密文
	EncryptBase64(rawData []byte) (string, error)
	// Decrypt 解密密文字节
	Decrypt(cipherData []byte) ([]byte, error)
	// DecryptBase64 解密 Base64 密文并返回明文字符串
	DecryptBase64(base64CipherData string) (string, error)
}
