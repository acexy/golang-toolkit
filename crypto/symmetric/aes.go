package symmetric

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
)

const (
	AESModeCBC AESMode = iota
	AESModeGCM
)

type AESMode int

// AESEncrypt AES加密实现
type AESEncrypt struct {
	key            []byte         // 私有化密钥
	mode           AESMode        // 加密模式
	ivCreator      IVCreator      // IV规则
	resultCreator  ResultCreator  // 结果返回规则
	paddingCreator PaddingCreator // 填充规则
	block          cipher.Block
	once           sync.Once
}

type AESOption struct {
	Mode           AESMode        // 加密模式
	IVCreator      IVCreator      // IV规则
	ResultCreator  ResultCreator  // 结果返回规则
	PaddingCreator PaddingCreator // 填充规则
}

// NewAES 创建AES加密实现
func NewAES(key []byte) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24, or 32 bytes")
	}
	// 复制密钥以避免外部修改
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

// NewAESWithOption 创建AES加密实现
func NewAESWithOption(key []byte, option AESOption) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid key size: must be 16, 24, or 32 bytes")
	}

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	aesInstance := &AESEncrypt{
		key:  keyCopy,
		mode: AESModeCBC,
	}

	// 设置模式
	if option.Mode == AESModeGCM {
		aesInstance.mode = AESModeGCM
	}

	// 设置IV创建器 - GCM模式使用专门的nonce创建器
	if option.IVCreator != nil {
		aesInstance.ivCreator = option.IVCreator
	} else {
		if aesInstance.mode == AESModeGCM {
			aesInstance.ivCreator = &RandomGCMNonceCreator{}
		} else {
			aesInstance.ivCreator = &RandomIvCreator{}
		}
	}

	// 设置结果创建器
	if option.ResultCreator != nil {
		aesInstance.resultCreator = option.ResultCreator
	} else {
		aesInstance.resultCreator = &AppendResultCreator{}
	}

	// 设置填充创建器（仅CBC模式需要）
	if aesInstance.mode == AESModeCBC {
		if option.PaddingCreator != nil {
			aesInstance.paddingCreator = option.PaddingCreator
		} else {
			aesInstance.paddingCreator = &Pkcs7PaddingCreator{}
		}
	}

	return aesInstance, nil
}

// IVCreator IV自定义创建接口
type IVCreator interface {
	// CreateForEncrypt 加密时创建IV
	CreateForEncrypt(key, rawData []byte) ([]byte, error)
	// ExtractForDecrypt 解密时提取IV
	ExtractForDecrypt(key, cipherData []byte) ([]byte, error)
}

// ResultCreator 加密块的结果处理
type ResultCreator interface {
	// CombineResult 加密时组合IV和密文
	CombineResult(iv, cipherData []byte) []byte
	// SeparateResult 解密时分离IV和密文
	SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error)
}

// PaddingCreator 填充接口
type PaddingCreator interface {
	// Pad 填充
	Pad(rawData []byte, blockSize int) ([]byte, error)
	// UnPad 去除填充
	UnPad(paddedData []byte) ([]byte, error)
}

// RandomIvCreator 随机生成IV (CBC模式使用)
// 改模式需要配合 AppendResultCreator 才能正常使用，解秘时需要解析出iv
type RandomIvCreator struct{}

func (r *RandomIvCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate random IV: %w", err)
	}
	return iv, nil
}

func (r *RandomIvCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < aes.BlockSize {
		return nil, errors.New("ciphertext too short to contain IV")
	}
	return cipherData[:aes.BlockSize], nil
}

// RandomGCMNonceCreator 随机生成GCM nonce (GCM模式使用)
type RandomGCMNonceCreator struct{}

func (r *RandomGCMNonceCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	// GCM标准推荐使用12字节的nonce
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %w", err)
	}
	return nonce, nil
}

func (r *RandomGCMNonceCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < 12 {
		return nil, errors.New("ciphertext too short to contain nonce")
	}
	return cipherData[:12], nil
}

// AppendResultCreator IV+密文的拼接方式
type AppendResultCreator struct{}

func (a *AppendResultCreator) CombineResult(iv, cipherData []byte) []byte {
	result := make([]byte, len(iv)+len(cipherData))
	copy(result, iv)
	copy(result[len(iv):], cipherData)
	return result
}

func (a *AppendResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	if len(combinedData) < ivSize {
		return nil, errors.New("combined data too short to contain IV")
	}
	return combinedData[ivSize:], nil
}

// PureResultCreator 纯密文方式（需要外部管理IV）
type PureResultCreator struct{}

func (p *PureResultCreator) CombineResult(iv, cipherData []byte) []byte {
	return cipherData
}

func (p *PureResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	// 纯密文模式下无法从数据中提取IV，需要外部提供
	return combinedData, nil
}

// Pkcs7PaddingCreator PKCS7填充
type Pkcs7PaddingCreator struct{}

func (p *Pkcs7PaddingCreator) Pad(rawData []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 || blockSize > 255 {
		return nil, errors.New("invalid block size")
	}

	padding := blockSize - len(rawData)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	result := make([]byte, len(rawData)+len(padText))
	copy(result, rawData)
	copy(result[len(rawData):], padText)

	return result, nil
}

func (p *Pkcs7PaddingCreator) UnPad(paddedData []byte) ([]byte, error) {
	if len(paddedData) == 0 {
		return nil, errors.New("empty padded data")
	}

	padding := int(paddedData[len(paddedData)-1])
	if padding == 0 || padding > len(paddedData) {
		return nil, errors.New("invalid padding size")
	}

	// 使用常量时间比较防止padding oracle攻击
	for i := 0; i < padding; i++ {
		if subtle.ConstantTimeByteEq(paddedData[len(paddedData)-1-i], byte(padding)) != 1 {
			return nil, errors.New("invalid padding content")
		}
	}

	return paddedData[:len(paddedData)-padding], nil
}

// cipherBlock 获取cipher.Block实例
func (a *AESEncrypt) cipherBlock() (cipher.Block, error) {
	var err error
	a.once.Do(func() {
		a.block, err = aes.NewCipher(a.key)
	})
	return a.block, err
}

// Encrypt 加密数据
func (a *AESEncrypt) Encrypt(rawData []byte) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, errors.New("empty data to encrypt")
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	switch a.mode {
	case AESModeCBC:
		return a.encryptCBC(block, rawData)
	case AESModeGCM:
		return a.encryptGCM(block, rawData)
	default:
		return nil, errors.New("unsupported AES mode")
	}
}

// encryptCBC CBC模式加密
func (a *AESEncrypt) encryptCBC(block cipher.Block, rawData []byte) ([]byte, error) {
	// 填充数据
	paddedData, err := a.paddingCreator.Pad(rawData, aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("padding failed: %w", err)
	}

	// 生成IV
	iv, err := a.ivCreator.CreateForEncrypt(a.key, paddedData)
	if err != nil {
		return nil, fmt.Errorf("IV creation failed: %w", err)
	}

	// 加密
	cipherData := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, paddedData)

	// 组合结果
	return a.resultCreator.CombineResult(iv, cipherData), nil
}

// encryptGCM GCM模式加密
func (a *AESEncrypt) encryptGCM(block cipher.Block, rawData []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成nonce
	nonce, err := a.ivCreator.CreateForEncrypt(a.key, rawData)
	if err != nil {
		return nil, fmt.Errorf("nonce creation failed: %w", err)
	}

	// 确保nonce长度正确
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size: expected %d, got %d", gcm.NonceSize(), len(nonce))
	}

	// GCM加密（包含认证标签）
	cipherData := gcm.Seal(nil, nonce, rawData, nil)

	// 组合结果
	return a.resultCreator.CombineResult(nonce, cipherData), nil
}

// EncryptBase64 加密并返回Base64字符串
func (a *AESEncrypt) EncryptBase64(rawData []byte) (string, error) {
	cipherData, err := a.Encrypt(rawData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// Decrypt 解密数据
func (a *AESEncrypt) Decrypt(cipherData []byte) ([]byte, error) {
	if len(cipherData) == 0 {
		return nil, errors.New("empty cipher data")
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	switch a.mode {
	case AESModeCBC:
		return a.decryptCBC(block, cipherData)
	case AESModeGCM:
		return a.decryptGCM(block, cipherData)
	default:
		return nil, errors.New("unsupported AES mode")
	}
}

// decryptCBC CBC模式解密
func (a *AESEncrypt) decryptCBC(block cipher.Block, cipherData []byte) ([]byte, error) {
	// 使用IVCreator提取IV
	iv, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract IV: %w", err)
	}

	// 分离实际密文
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(iv))
	if err != nil {
		return nil, fmt.Errorf("failed to separate cipher data: %w", err)
	}

	if len(actualCipherData)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext length is not a multiple of block size")
	}

	// 解密
	rawData := make([]byte, len(actualCipherData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(rawData, actualCipherData)

	// 去除填充
	return a.paddingCreator.UnPad(rawData)
}

// decryptGCM GCM模式解密
func (a *AESEncrypt) decryptGCM(block cipher.Block, cipherData []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// 使用IVCreator提取nonce
	nonce, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, fmt.Errorf("failed to extract nonce: %w", err)
	}

	// 分离实际密文
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(nonce))
	if err != nil {
		return nil, fmt.Errorf("failed to separate cipher data: %w", err)
	}

	// 验证nonce长度
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("invalid nonce size: expected %d, got %d", gcm.NonceSize(), len(nonce))
	}

	// GCM解密（包含认证验证）
	return gcm.Open(nil, nonce, actualCipherData, nil)
}

// DecryptBase64 解密Base64字符串
func (a *AESEncrypt) DecryptBase64(base64CipherData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64CipherData)
	if err != nil {
		return "", fmt.Errorf("invalid base64 data: %w", err)
	}

	plain, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
