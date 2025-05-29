package symmetric

// CryptEncrypt 对称加密通用接口
type CryptEncrypt interface {

	// Encrypt 加密
	Encrypt(raw []byte) ([]byte, error)
	// EncryptBase64 加密
	EncryptBase64(base64Raw string) (string, error)
	// Decrypt 解密
	Decrypt(cipherText []byte) ([]byte, error)
	// DecryptBase64 解密
	DecryptBase64(base64Cipher string) (string, error)
}
