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
	// AESModeCBC 表示 AES CBC 加密模式
	AESModeCBC AESMode = iota
	// AESModeGCM 表示 AES GCM 加密模式
	AESModeGCM
)

// AESMode 表示 AES 加密模式
type AESMode int

// AESEncrypt 提供 AES 对称加解密能力
type AESEncrypt struct {
	key            []byte         // 私有化密钥
	mode           AESMode        // 加密模式
	ivCreator      IVCreator      // IV规则
	resultCreator  ResultCreator  // 结果返回规则
	paddingCreator PaddingCreator // 填充规则
	block          cipher.Block
	blockErr       error
	once           sync.Once
}

// AESOption 表示创建 AES 实例时使用的可选配置
type AESOption struct {
	Mode           AESMode        // 加密模式
	IVCreator      IVCreator      // IV规则
	ResultCreator  ResultCreator  // 结果返回规则
	PaddingCreator PaddingCreator // 填充规则
}

// NewAES 创建默认 CBC 模式的 AES 加解密实例
func NewAES(key []byte) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, toolkitError.NewInvalidAESKeySizeError(len(key))
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

// NewAESWithOption 使用指定配置创建 AES 加解密实例
func NewAESWithOption(key []byte, option AESOption) (*AESEncrypt, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, toolkitError.NewInvalidAESKeySizeError(len(key))
	}

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	aesInstance := &AESEncrypt{
		key:  keyCopy,
		mode: AESModeCBC,
	}

	// 设置模式
	switch option.Mode {
	case AESModeCBC:
		aesInstance.mode = AESModeCBC
	case AESModeGCM:
		aesInstance.mode = AESModeGCM
	default:
		return nil, toolkitError.ErrUnsupportedAESMode
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

// IVCreator 定义加密和解密时的 IV/nonce 处理规则
type IVCreator interface {
	// CreateForEncrypt 加密时创建IV
	CreateForEncrypt(key, rawData []byte) ([]byte, error)
	// ExtractForDecrypt 解密时提取IV
	ExtractForDecrypt(key, cipherData []byte) ([]byte, error)
}

// ResultCreator 定义 IV/nonce 与密文的组合和分离规则
type ResultCreator interface {
	// CombineResult 加密时组合IV和密文
	CombineResult(iv, cipherData []byte) []byte
	// SeparateResult 解密时分离IV和密文
	SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error)
}

// PaddingCreator 定义分组加密的填充和去填充规则
type PaddingCreator interface {
	// Pad 填充
	Pad(rawData []byte, blockSize int) ([]byte, error)
	// UnPad 去除填充
	UnPad(paddedData []byte) ([]byte, error)
}

// RandomIvCreator 为 CBC 模式随机生成 IV
// 该模式需要配合 AppendResultCreator 使用，解密时从密文前缀解析 IV。
type RandomIvCreator struct{}

// CreateForEncrypt 创建 CBC 加密使用的随机 IV
func (r *RandomIvCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateIVFailed, err)
	}
	return iv, nil
}

// ExtractForDecrypt 从密文前缀提取 CBC 解密使用的 IV
func (r *RandomIvCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < aes.BlockSize {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return cipherData[:aes.BlockSize], nil
}

// RandomGCMNonceCreator 随机生成GCM nonce (GCM模式使用)
type RandomGCMNonceCreator struct{}

// CreateForEncrypt 创建 GCM 加密使用的随机 nonce
func (r *RandomGCMNonceCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	// GCM标准推荐使用12字节的nonce
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateNonceFailed, err)
	}
	return nonce, nil
}

// ExtractForDecrypt 从密文前缀提取 GCM 解密使用的 nonce
func (r *RandomGCMNonceCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	if len(cipherData) < 12 {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return cipherData[:12], nil
}

// AppendResultCreator IV+密文的拼接方式
type AppendResultCreator struct{}

// CombineResult 将 IV/nonce 放在密文前缀
func (a *AppendResultCreator) CombineResult(iv, cipherData []byte) []byte {
	result := make([]byte, len(iv)+len(cipherData))
	copy(result, iv)
	copy(result[len(iv):], cipherData)
	return result
}

// SeparateResult 从组合数据中分离密文
func (a *AppendResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	if len(combinedData) < ivSize {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return combinedData[ivSize:], nil
}

// PureResultCreator 纯密文方式（需要外部管理IV）
type PureResultCreator struct{}

// CombineResult 仅返回密文，调用方需要自行管理 IV/nonce
func (p *PureResultCreator) CombineResult(iv, cipherData []byte) []byte {
	return cipherData
}

// SeparateResult 直接返回密文，调用方需要通过 IVCreator 提供 IV/nonce
func (p *PureResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	// 纯密文模式下无法从数据中提取IV，需要外部提供
	return combinedData, nil
}

// Pkcs7PaddingCreator PKCS7填充
type Pkcs7PaddingCreator struct{}

// Pad 按 PKCS7 规则填充数据
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

// UnPad 按 PKCS7 规则移除填充
func (p *Pkcs7PaddingCreator) UnPad(paddedData []byte) ([]byte, error) {
	if len(paddedData) == 0 {
		return nil, toolkitError.ErrInvalidPadding
	}

	padding := int(paddedData[len(paddedData)-1])
	if padding == 0 || padding > len(paddedData) {
		return nil, toolkitError.ErrInvalidPadding
	}

	// 使用常量时间比较防止padding oracle攻击
	for i := 0; i < padding; i++ {
		if subtle.ConstantTimeByteEq(paddedData[len(paddedData)-1-i], byte(padding)) != 1 {
			return nil, toolkitError.ErrInvalidPadding
		}
	}

	return paddedData[:len(paddedData)-padding], nil
}

// cipherBlock 获取cipher.Block实例
func (a *AESEncrypt) cipherBlock() (cipher.Block, error) {
	a.once.Do(func() {
		a.block, a.blockErr = aes.NewCipher(a.key)
	})
	return a.block, a.blockErr
}

// Encrypt 加密数据
func (a *AESEncrypt) Encrypt(rawData []byte) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, toolkitError.ErrEmptyEncryptData
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateCipherBlockFailed, err)
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

// encryptCBC CBC模式加密
func (a *AESEncrypt) encryptCBC(block cipher.Block, rawData []byte) ([]byte, error) {
	// 填充数据
	paddedData, err := a.paddingCreator.Pad(rawData, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	// 生成IV
	iv, err := a.ivCreator.CreateForEncrypt(a.key, paddedData)
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateIVFailed, err)
	}
	if len(iv) != aes.BlockSize {
		return nil, toolkitError.NewInvalidIVSizeError(aes.BlockSize, len(iv))
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
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateGCMFailed, err)
	}

	// 生成nonce
	nonce, err := a.ivCreator.CreateForEncrypt(a.key, rawData)
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateNonceFailed, err)
	}

	// 确保nonce长度正确
	if len(nonce) != gcm.NonceSize() {
		return nil, toolkitError.NewInvalidNonceSizeError(gcm.NonceSize(), len(nonce))
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
		return nil, toolkitError.ErrEmptyCipherData
	}

	block, err := a.cipherBlock()
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateCipherBlockFailed, err)
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

// decryptCBC CBC模式解密
func (a *AESEncrypt) decryptCBC(block cipher.Block, cipherData []byte) ([]byte, error) {
	// 使用IVCreator提取IV
	iv, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrExtractIVFailed, err)
	}
	if len(iv) != aes.BlockSize {
		return nil, toolkitError.NewInvalidIVSizeError(aes.BlockSize, len(iv))
	}

	// 分离实际密文
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(iv))
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrSeparateCipherDataFailed, err)
	}

	if len(actualCipherData)%aes.BlockSize != 0 {
		return nil, toolkitError.ErrInvalidBlockSize
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
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrCreateGCMFailed, err)
	}

	// 使用IVCreator提取nonce
	nonce, err := a.ivCreator.ExtractForDecrypt(a.key, cipherData)
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrExtractNonceFailed, err)
	}

	// 分离实际密文
	actualCipherData, err := a.resultCreator.SeparateResult(cipherData, len(nonce))
	if err != nil {
		return nil, toolkitError.WrapSymmetricError(toolkitError.ErrSeparateCipherDataFailed, err)
	}

	// 验证nonce长度
	if len(nonce) != gcm.NonceSize() {
		return nil, toolkitError.NewInvalidNonceSizeError(gcm.NonceSize(), len(nonce))
	}

	// GCM解密（包含认证验证）
	return gcm.Open(nil, nonce, actualCipherData, nil)
}

// DecryptBase64 解密Base64字符串
func (a *AESEncrypt) DecryptBase64(base64CipherData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(base64CipherData)
	if err != nil {
		return "", toolkitError.WrapSymmetricError(toolkitError.ErrInvalidBase64Data, err)
	}

	plain, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
