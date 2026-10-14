package symmetric

// CryptEncrypt defines symmetric byte and Base64 encryption operations.
type CryptEncrypt interface {
	// Encrypt encrypts nonempty plaintext bytes.
	Encrypt(rawData []byte) ([]byte, error)
	// EncryptBase64 encrypts plaintext bytes and returns Base64 ciphertext.
	EncryptBase64(rawData []byte) (string, error)
	// Decrypt decrypts nonempty ciphertext bytes.
	Decrypt(cipherData []byte) ([]byte, error)
	// DecryptBase64 decodes Base64 ciphertext and returns the plaintext string.
	DecryptBase64(base64CipherData string) (string, error)
}
