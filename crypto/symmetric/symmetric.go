package symmetric

// CryptEncrypt 对称加密通用接口
type CryptEncrypt interface {

	// Encrypt 加密
	Encrypt(rawData []byte) ([]byte, error)
	// EncryptBase64 加密
	EncryptBase64(rawData []byte) (string, error)
	// Decrypt 解密
	Decrypt(cipherData []byte) ([]byte, error)
	// DecryptBase64 解密
	DecryptBase64(base64CipherData string) (string, error)
}
