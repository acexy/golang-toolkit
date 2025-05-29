package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
)

const (
	AESModelCBC AESModel = iota
)

type AESModel int

// AESEncrypt AES加密实现
type AESEncrypt struct {
	Key   []byte
	Model AESModel

	block cipher.Block
	once  sync.Once
}

func NewAES(key []byte) *AESEncrypt {
	return &AESEncrypt{
		Key: key,
	}
}

func encIV() []byte {
	iv := make([]byte, aes.BlockSize)
	_, _ = rand.Read(iv) // 安全随机生成
	return iv
}

func decIV(plaintext []byte) []byte {
	return plaintext[:aes.BlockSize]
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-padding], nil
}

func (a *AESEncrypt) cipherBlock() (cipher.Block, error) {
	var err error
	a.once.Do(func() {
		a.block, err = aes.NewCipher(a.Key)
	})
	return a.block, err
}

func (a *AESEncrypt) Encrypt(raw []byte) ([]byte, error) {
	block, err := a.cipherBlock()
	if err != nil {
		return nil, err
	}
	var ciphertext []byte
	if a.Model == AESModelCBC {
		raw = pkcs7Pad(raw, aes.BlockSize)
		if len(raw)%aes.BlockSize != 0 {
			return nil, errors.New("plaintext is not a multiple of the block size")
		}
		ciphertext = make([]byte, len(raw))
		iv := encIV()
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext, raw)
		ciphertext = append(iv, ciphertext...)
	} else {
		// TODO: 其他模式支持
		return nil, errors.New("not supported aes model")
	}
	return ciphertext, nil
}

func (a *AESEncrypt) EncryptBase64(base64Raw string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	ciphertext, err := a.Encrypt(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *AESEncrypt) Decrypt(cipherText []byte) ([]byte, error) {
	block, err := a.cipherBlock()
	if err != nil {
		return nil, err
	}
	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	plaintext := make([]byte, len(cipherText)-aes.BlockSize)
	if a.Model == AESModelCBC {
		mode := cipher.NewCBCDecrypter(block, decIV(cipherText))
		mode.CryptBlocks(plaintext, cipherText[aes.BlockSize:])
		plaintext, err = pkcs7Unpad(plaintext)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: 其他模式支持
		return nil, errors.New("not supported aes model")
	}
	return plaintext, nil
}

func (a *AESEncrypt) DecryptBase64(base64Cipher string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return "", err
	}
	decryptBase64, err := a.Decrypt(data)
	return base64.StdEncoding.EncodeToString(decryptBase64), nil
}
