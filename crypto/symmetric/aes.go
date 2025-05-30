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
	Key            []byte         // 加密密钥
	Model          AESModel       //  加密模式
	IVCreator      IVCreator      //  IV规则 不指定时使用默认方式
	ResultCreator  ResultCreator  //  最终结果返回规则 不指定则使用默认方式
	PaddingCreator PaddingCreator //  填充规则
	block          cipher.Block
	once           sync.Once
}

type AESOption struct {
	IVCreator      IVCreator      //  IV规则 不指定时使用默认方式
	Model          AESModel       //  加密模式
	ResultCreator  ResultCreator  //  最终结果返回规则 不指定则使用默认方式
	PaddingCreator PaddingCreator //  填充规则
}

// NewAES 创建AES加密实现
func NewAES(key []byte) *AESEncrypt {
	return &AESEncrypt{
		Key:            key,
		Model:          AESModelCBC,
		IVCreator:      randomIvCreator{},
		ResultCreator:  appendResultCreator{},
		PaddingCreator: pkcs7PaddingCreator{},
	}
}

// NewAESWithOption 创建AES加密实现
func NewAESWithOption(key []byte, option AESOption) *AESEncrypt {
	return &AESEncrypt{
		Key: key,
		Model: func() AESModel {
			if option.Model == AESModelCBC {
				return AESModelCBC
			}
			return option.Model
		}(),
		IVCreator: func() IVCreator {
			if option.IVCreator == nil {
				return randomIvCreator{}
			} else {
				return option.IVCreator
			}
		}(),
		ResultCreator: func() ResultCreator {
			if option.ResultCreator == nil {
				return appendResultCreator{}
			} else {
				return option.ResultCreator
			}
		}(),
		PaddingCreator: func() PaddingCreator {
			if option.PaddingCreator == nil {
				return pkcs7PaddingCreator{}
			} else {
				return option.PaddingCreator
			}
		}(),
	}
}

// IVCreator IV自定义创建接口
type IVCreator interface {
	// Encrypt 加密时创建IV
	Encrypt(key, paddedRawData []byte) [aes.BlockSize]byte
	// Decrypt 解密时创建IV
	Decrypt(key, cipherText []byte) [aes.BlockSize]byte
}

type randomIvCreator struct {
}

func (d randomIvCreator) Encrypt(key, paddedRawData []byte) [aes.BlockSize]byte {
	iv := make([]byte, aes.BlockSize)
	_, _ = rand.Read(iv) // 安全随机生成
	return [aes.BlockSize]byte(iv)
}
func (d randomIvCreator) Decrypt(key, cipherText []byte) [aes.BlockSize]byte {
	return [aes.BlockSize]byte(cipherText[:aes.BlockSize])
}

type ResultCreator interface {
	// Encrypt 加密时创建最终返回内容
	Encrypt(iv [aes.BlockSize]byte, rawCipherData []byte) []byte
	// Decrypt 解密时通过原始的密文创建解密的实际密文
	Decrypt(iv [aes.BlockSize]byte, cipherData []byte) []byte
}

type appendResultCreator struct {
}

func (a appendResultCreator) Encrypt(iv [aes.BlockSize]byte, rawCipherData []byte) []byte {
	return append(iv[:], rawCipherData...)
}

func (a appendResultCreator) Decrypt(iv [aes.BlockSize]byte, cipherText []byte) []byte {
	return cipherText[len(iv):]
}

// PaddingCreator 填充接口
type PaddingCreator interface {
	// Pad 填充
	Pad(rawData []byte, blockSize int) []byte
	// UnPad 去除填充
	UnPad(rawData []byte) ([]byte, error)
}

type pkcs7PaddingCreator struct {
}

func (p pkcs7PaddingCreator) Pad(rawData []byte, blockSize int) []byte {
	padding := blockSize - len(rawData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(rawData, padText...)
}

func (p pkcs7PaddingCreator) UnPad(rawData []byte) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(rawData[len(rawData)-1])
	if padding == 0 || padding > len(rawData) {
		return nil, errors.New("invalid padding size")
	}
	for i := 0; i < padding; i++ {
		if rawData[len(rawData)-1-i] != byte(padding) {
			return nil, errors.New("invalid padding content")
		}
	}
	return rawData[:len(rawData)-padding], nil
}

func (a *AESEncrypt) cipherBlock() (cipher.Block, error) {
	var err error
	a.once.Do(func() {
		a.block, err = aes.NewCipher(a.Key)
	})
	return a.block, err
}

func (a *AESEncrypt) Encrypt(rawData []byte) ([]byte, error) {
	block, err := a.cipherBlock()
	if err != nil {
		return nil, err
	}
	var rawCipherData []byte
	if a.Model == AESModelCBC {
		paddedRawData := a.PaddingCreator.Pad(rawData, aes.BlockSize)
		if len(paddedRawData)%aes.BlockSize != 0 {
			return nil, errors.New("plaintext is not a multiple of the block size")
		}
		rawCipherData = make([]byte, len(paddedRawData))
		iv := a.IVCreator.Encrypt(a.Key, paddedRawData)
		mode := cipher.NewCBCEncrypter(block, iv[:])
		mode.CryptBlocks(rawCipherData, paddedRawData)
		rawCipherData = a.ResultCreator.Encrypt(iv, rawCipherData)
	} else {
		// TODO: 其他模式支持
		return nil, errors.New("not supported aes model")
	}
	return rawCipherData, nil
}

func (a *AESEncrypt) EncryptBase64(base64RawData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64RawData)
	if err != nil {
		return "", err
	}
	cipherData, err := a.Encrypt(data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

func (a *AESEncrypt) Decrypt(cipherData []byte) ([]byte, error) {
	block, err := a.cipherBlock()
	if err != nil {
		return nil, err
	}

	var rawData []byte
	if a.Model == AESModelCBC {
		iv := a.IVCreator.Decrypt(a.Key, cipherData)
		mode := cipher.NewCBCDecrypter(block, iv[:])
		rawCipherData := a.ResultCreator.Decrypt(iv, cipherData)
		if len(rawCipherData)%aes.BlockSize != 0 {
			return nil, errors.New("ciphertext is not a multiple of the block size")
		}
		rawData = make([]byte, len(rawCipherData))
		mode.CryptBlocks(rawData, rawCipherData)
		rawData, err = a.PaddingCreator.UnPad(rawData)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO: 其他模式支持
		return nil, errors.New("not supported aes model")
	}
	return rawData, nil
}

func (a *AESEncrypt) DecryptBase64(base64CipherData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64CipherData)
	if err != nil {
		return "", err
	}
	plain, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
