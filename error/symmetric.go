package error

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidAESKeySize 表示 AES 密钥长度无效
	ErrInvalidAESKeySize = errors.New("invalid AES key size")

	// ErrEmptyEncryptData 表示待加密数据为空
	ErrEmptyEncryptData = errors.New("empty data to encrypt")

	// ErrEmptyCipherData 表示待解密密文为空
	ErrEmptyCipherData = errors.New("empty cipher data")

	// ErrUnsupportedAESMode 表示不支持的 AES 加密模式
	ErrUnsupportedAESMode = errors.New("unsupported AES mode")

	// ErrInvalidIVSize 表示 IV 长度无效
	ErrInvalidIVSize = errors.New("invalid IV size")

	// ErrInvalidNonceSize 表示 nonce 长度无效
	ErrInvalidNonceSize = errors.New("invalid nonce size")

	// ErrCipherDataTooShort 表示密文长度不足
	ErrCipherDataTooShort = errors.New("cipher data too short")

	// ErrInvalidBlockSize 表示 block size 无效
	ErrInvalidBlockSize = errors.New("invalid block size")

	// ErrInvalidPadding 表示 padding 内容无效
	ErrInvalidPadding = errors.New("invalid padding")

	// ErrInvalidBase64Data 表示 Base64 数据无效
	ErrInvalidBase64Data = errors.New("invalid base64 data")

	// ErrCreateCipherBlockFailed 表示创建 cipher block 失败
	ErrCreateCipherBlockFailed = errors.New("failed to create cipher block")

	// ErrCreateGCMFailed 表示创建 GCM 实例失败
	ErrCreateGCMFailed = errors.New("failed to create GCM")

	// ErrCreateIVFailed 表示创建 IV 失败
	ErrCreateIVFailed = errors.New("failed to create IV")

	// ErrCreateNonceFailed 表示创建 nonce 失败
	ErrCreateNonceFailed = errors.New("failed to create nonce")

	// ErrExtractIVFailed 表示提取 IV 失败
	ErrExtractIVFailed = errors.New("failed to extract IV")

	// ErrExtractNonceFailed 表示提取 nonce 失败
	ErrExtractNonceFailed = errors.New("failed to extract nonce")

	// ErrSeparateCipherDataFailed 表示分离密文失败
	ErrSeparateCipherDataFailed = errors.New("failed to separate cipher data")
)

// NewInvalidAESKeySizeError 创建 AES 密钥长度无效错误
func NewInvalidAESKeySizeError(size int) error {
	return fmt.Errorf("%w: must be 16, 24, or 32 bytes, got %d", ErrInvalidAESKeySize, size)
}

// NewInvalidIVSizeError 创建 IV 长度无效错误
func NewInvalidIVSizeError(expected, actual int) error {
	return fmt.Errorf("%w: expected %d, got %d", ErrInvalidIVSize, expected, actual)
}

// NewInvalidNonceSizeError 创建 nonce 长度无效错误
func NewInvalidNonceSizeError(expected, actual int) error {
	return fmt.Errorf("%w: expected %d, got %d", ErrInvalidNonceSize, expected, actual)
}

// WrapSymmetricError 包装对称加密错误
func WrapSymmetricError(base, err error) error {
	if err == nil {
		return base
	}
	return fmt.Errorf("%w: %w", base, err)
}
