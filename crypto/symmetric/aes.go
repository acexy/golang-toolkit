package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"sync"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

const (
	// AESModeCBC selects CBC mode with padding.
	AESModeCBC AESMode = iota
	// AESModeGCM selects authenticated GCM mode.
	AESModeGCM
)

// AESMode identifies the AES encryption mode.
type AESMode int

// AESEncrypt encrypts and decrypts with a copied AES key and configured strategies.
type AESEncrypt struct {
	key            []byte         // Copied key material.
	mode           AESMode        // Encryption mode.
	ivCreator      IVCreator      // IV/nonce strategy.
	resultCreator  ResultCreator  // Ciphertext framing strategy.
	paddingCreator PaddingCreator // CBC padding strategy.
	block          cipher.Block
	blockErr       error
	once           sync.Once
}

// AESOption configures AES mode, IV/nonce handling, output framing, and CBC padding.
type AESOption struct {
	// Mode selects CBC or GCM mode.
	Mode AESMode
	// IVCreator controls IV/nonce creation and extraction; nil selects the mode default.
	IVCreator IVCreator
	// ResultCreator controls ciphertext framing; nil prefixes the IV/nonce.
	ResultCreator ResultCreator
	// PaddingCreator controls CBC padding; nil selects PKCS7 and GCM ignores padding.
	PaddingCreator PaddingCreator
}

// NewAES creates a CBC instance with random IVs and PKCS7 padding; keys must be 16, 24, or 32 bytes.
func NewAES(key []byte) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, toolkitError.ErrInvalidAESKeySize
	}
	// Copy the key to isolate it from caller mutation.
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	return &AESEncrypt{
		key:            keyCopy,
		mode:           AESModeCBC,
		ivCreator:      &RandomIvCreator{},
		resultCreator:  &AppendResultCreator{},
		paddingCreator: &Pkcs7PaddingCreator{},
	}, nil
}

// NewAESWithOption creates an instance with copied key material and optional strategies.
func NewAESWithOption(key []byte, option AESOption) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, toolkitError.ErrInvalidAESKeySize
	}

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	aesInstance := &AESEncrypt{
		key:  keyCopy,
		mode: AESModeCBC,
	}

	// Select the encryption mode.
	switch option.Mode {
	case AESModeCBC:
		aesInstance.mode = AESModeCBC
	case AESModeGCM:
		aesInstance.mode = AESModeGCM
	default:
		return nil, toolkitError.ErrUnsupportedAESMode
	}

	// Select a nonce strategy for GCM and an IV strategy for CBC.
	if option.IVCreator != nil {
		aesInstance.ivCreator = option.IVCreator
	} else {
		if aesInstance.mode == AESModeGCM {
			aesInstance.ivCreator = &RandomGCMNonceCreator{}
		} else {
			aesInstance.ivCreator = &RandomIvCreator{}
		}
	}

	// Select the ciphertext framing strategy.
	if option.ResultCreator != nil {
		aesInstance.resultCreator = option.ResultCreator
	} else {
		aesInstance.resultCreator = &AppendResultCreator{}
	}

	// Only CBC requires a padding strategy.
	if aesInstance.mode == AESModeCBC {
		if option.PaddingCreator != nil {
			aesInstance.paddingCreator = option.PaddingCreator
		} else {
			aesInstance.paddingCreator = &Pkcs7PaddingCreator{}
		}
	}

	return aesInstance, nil
}

// IVCreator defines IV/nonce creation and ciphertext extraction.
type IVCreator interface {
	// CreateForEncrypt creates the IV or nonce used for encryption.
	CreateForEncrypt(key, rawData []byte) ([]byte, error)
	// ExtractForDecrypt extracts the IV or nonce required for decryption.
	ExtractForDecrypt(key, cipherData []byte) ([]byte, error)
}

// ResultCreator defines how the IV/nonce and ciphertext are framed.
type ResultCreator interface {
	// CombineResult combines the IV/nonce and ciphertext.
	CombineResult(iv, cipherData []byte) []byte
	// SeparateResult extracts ciphertext from the combined result.
	SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error)
}

// PaddingCreator defines block padding and validation.
type PaddingCreator interface {
	// Pad pads plaintext to the block size.
	Pad(rawData []byte, blockSize int) ([]byte, error)
	// UnPad validates and removes padding.
	UnPad(paddedData []byte) ([]byte, error)
}

// RandomIvCreator generates CBC IVs and reads them from a ciphertext prefix.
type RandomIvCreator struct{}

// CreateForEncrypt creates a random CBC IV.
func (r *RandomIvCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, toolkitError.ErrCreateIVFailed
	}
	return iv, nil
}

// ExtractForDecrypt reads a CBC IV from the ciphertext prefix.
func (r *RandomIvCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < aes.BlockSize {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return cipherData[:aes.BlockSize], nil
}

// RandomGCMNonceCreator generates GCM nonces and reads them from a ciphertext prefix.
type RandomGCMNonceCreator struct{}

// CreateForEncrypt creates a random 12-byte GCM nonce.
func (r *RandomGCMNonceCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	// Use the default 12-byte GCM nonce size.
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, toolkitError.ErrCreateNonceFailed
	}
	return nonce, nil
}

// ExtractForDecrypt reads a GCM nonce from the ciphertext prefix.
func (r *RandomGCMNonceCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < 12 {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return cipherData[:12], nil
}

// AppendResultCreator prefixes the ciphertext with its IV or nonce.
type AppendResultCreator struct{}

// CombineResult prefixes ciphertext with the IV/nonce.
func (a *AppendResultCreator) CombineResult(iv, cipherData []byte) []byte {
	result := make([]byte, len(iv)+len(cipherData))
	copy(result, iv)
	copy(result[len(iv):], cipherData)
	return result
}

// SeparateResult removes the IV/nonce prefix from ciphertext.
func (a *AppendResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	if len(combinedData) < ivSize {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return combinedData[ivSize:], nil
}

// PureResultCreator returns ciphertext without including its IV or nonce.
type PureResultCreator struct{}

// CombineResult returns ciphertext without the IV/nonce.
func (p *PureResultCreator) CombineResult(iv, cipherData []byte) []byte {
	return cipherData
}

// SeparateResult returns the input ciphertext unchanged.
func (p *PureResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	// Pure ciphertext has no IV prefix; extraction needs an external IV strategy.
	return combinedData, nil
}

// Pkcs7PaddingCreator implements PKCS7 padding and validation.
type Pkcs7PaddingCreator struct{}

// Pad pads data to the requested block size.
func (p *Pkcs7PaddingCreator) Pad(rawData []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, toolkitError.ErrInvalidBlockSize
	}

	padding := blockSize - len(rawData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	result := make([]byte, len(rawData)+len(padText))
	copy(result, rawData)
	copy(result[len(rawData):], padText)

	return result, nil
}

// UnPad validates and removes PKCS7 padding.
func (p *Pkcs7PaddingCreator) UnPad(paddedData []byte) ([]byte, error) {
	if len(paddedData) == 0 {
		return nil, toolkitError.ErrInvalidPadding
	}

	padding := int(paddedData[len(paddedData)-1])
	if padding == 0 || padding > len(paddedData) {
		return nil, toolkitError.ErrInvalidPadding
	}

	// Compare padding bytes without early exit; CBC still requires separate authentication.
	for i := 0; i < padding; i++ {
		if subtle.ConstantTimeByteEq(paddedData[len(paddedData)-1-i], byte(padding)) != 1 {
			return nil, toolkitError.ErrInvalidPadding
		}
	}

	return paddedData[:len(paddedData)-padding], nil
}

// cipherBlock initializes the AES block once.
func (a *AESEncrypt) cipherBlock() (cipher.Block, error) {
	a.once.Do(func() {
		a.block, a.blockErr = aes.NewCipher(a.key)
	})
	return a.block, a.blockErr
}

// Encrypt encrypts nonempty plaintext bytes.
func (a *AESEncrypt) Encrypt(rawData []byte) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, toolkitError.ErrEmptyEncryptData
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, toolkitError.ErrCreateCipherBlockFailed
	}

	switch a.mode {
	case AESModeCBC:
		return a.encryptCBC(block, rawData)
	case AESModeGCM:
		return a.encryptGCM(block, rawData)
	default:
		return nil, toolkitError.ErrUnsupportedAESMode
	}
}

// encryptCBC applies padding, a fresh IV, and ciphertext framing.
func (a *AESEncrypt) encryptCBC(block cipher.Block, rawData []byte) ([]byte, error) {
	// Pad the plaintext.
	paddedData, err := a.paddingCreator.Pad(rawData, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	// Create the IV.
	iv, err := a.ivCreator.CreateForEncrypt(a.key, paddedData)
	if err != nil {
		return nil, toolkitError.ErrCreateIVFailed
	}
	if len(iv) != aes.BlockSize {
		return nil, toolkitError.ErrInvalidIVSize
	}

	// Encrypt the padded blocks.
	cipherData := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, paddedData)

	// Frame the ciphertext with its IV or nonce.
	return a.resultCreator.CombineResult(iv, cipherData), nil
}

// encryptGCM encrypts and authenticates using a fresh nonce.
func (a *AESEncrypt) encryptGCM(block cipher.Block, rawData []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, toolkitError.ErrCreateGCMFailed
	}

	// Create the nonce.
	nonce, err := a.ivCreator.CreateForEncrypt(a.key, rawData)
	if err != nil {
		return nil, toolkitError.ErrCreateNonceFailed
	}

	// Validate the nonce size before sealing.
	if len(nonce) != gcm.NonceSize() {
		return nil, toolkitError.ErrInvalidNonceSize
	}

	// Seal the plaintext, including the authentication tag.
	cipherData := gcm.Seal(nil, nonce, rawData, nil)

	// Frame the ciphertext with its IV or nonce.
	return a.resultCreator.CombineResult(nonce, cipherData), nil
}

// EncryptBase64 encrypts plaintext bytes and returns Base64 ciphertext.
func (a *AESEncrypt) EncryptBase64(rawData []byte) (string, error) {
	cipherData, err := a.Encrypt(rawData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// Decrypt decrypts nonempty ciphertext bytes.
func (a *AESEncrypt) Decrypt(cipherData []byte) ([]byte, error) {
	if len(cipherData) == 0 {
		return nil, toolkitError.ErrEmptyCipherData
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, toolkitError.ErrCreateCipherBlockFailed
	}

	switch a.mode {
	case AESModeCBC:
		return a.decryptCBC(block, cipherData)
	case AESModeGCM:
		return a.decryptGCM(block, cipherData)
	default:
		return nil, toolkitError.ErrUnsupportedAESMode
	}
}

// decryptCBC extracts the IV, decrypts blocks, and validates padding.
func (a *AESEncrypt) decryptCBC(block cipher.Block, cipherData []byte) ([]byte, error) {
	// Extract the IV using the configured strategy.
	iv, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, toolkitError.ErrExtractIVFailed
	}
	if len(iv) != aes.BlockSize {
		return nil, toolkitError.ErrInvalidIVSize
	}

	// Separate the ciphertext from its framing.
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(iv))
	if err != nil {
		return nil, toolkitError.ErrSeparateCipherDataFailed
	}

	if len(actualCipherData)%aes.BlockSize != 0 {
		return nil, toolkitError.ErrInvalidBlockSize
	}

	// Decrypt the blocks.
	rawData := make([]byte, len(actualCipherData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(rawData, actualCipherData)

	// Validate and remove padding.
	return a.paddingCreator.UnPad(rawData)
}

// decryptGCM extracts the nonce and authenticates before returning plaintext.
func (a *AESEncrypt) decryptGCM(block cipher.Block, cipherData []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, toolkitError.ErrCreateGCMFailed
	}

	// Extract the nonce using the configured strategy.
	nonce, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, toolkitError.ErrExtractNonceFailed
	}

	// Separate the ciphertext from its framing.
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(nonce))
	if err != nil {
		return nil, toolkitError.ErrSeparateCipherDataFailed
	}

	// Validate the nonce size before opening.
	if len(nonce) != gcm.NonceSize() {
		return nil, toolkitError.ErrInvalidNonceSize
	}

	// Open the ciphertext and verify its authentication tag.
	return gcm.Open(nil, nonce, actualCipherData, nil)
}

// DecryptBase64 decodes Base64 ciphertext and returns the plaintext string.
func (a *AESEncrypt) DecryptBase64(base64CipherData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64CipherData)
	if err != nil {
		return "", toolkitError.ErrInvalidBase64Data
	}

	plain, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
