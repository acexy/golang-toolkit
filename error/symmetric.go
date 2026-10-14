package error

import (
	"errors"
)

var (
	// ErrInvalidAESKeySize indicates invalid AES key size.
	ErrInvalidAESKeySize = errors.New("invalid AES key size")

	// ErrEmptyEncryptData indicates empty data to encrypt.
	ErrEmptyEncryptData = errors.New("empty data to encrypt")

	// ErrEmptyCipherData indicates empty cipher data.
	ErrEmptyCipherData = errors.New("empty cipher data")

	// ErrUnsupportedAESMode indicates unsupported AES mode.
	ErrUnsupportedAESMode = errors.New("unsupported AES mode")

	// ErrInvalidIVSize indicates invalid IV size.
	ErrInvalidIVSize = errors.New("invalid IV size")

	// ErrInvalidNonceSize indicates invalid nonce size.
	ErrInvalidNonceSize = errors.New("invalid nonce size")

	// ErrCipherDataTooShort indicates cipher data too short.
	ErrCipherDataTooShort = errors.New("cipher data too short")

	// ErrInvalidBlockSize indicates invalid block size.
	ErrInvalidBlockSize = errors.New("invalid block size")

	// ErrInvalidPadding indicates invalid padding.
	ErrInvalidPadding = errors.New("invalid padding")

	// ErrInvalidBase64Data indicates invalid base64 data.
	ErrInvalidBase64Data = errors.New("invalid base64 data")

	// ErrCreateCipherBlockFailed indicates failed to create cipher block.
	ErrCreateCipherBlockFailed = errors.New("failed to create cipher block")

	// ErrCreateGCMFailed indicates failed to create GCM.
	ErrCreateGCMFailed = errors.New("failed to create GCM")

	// ErrCreateIVFailed indicates failed to create IV.
	ErrCreateIVFailed = errors.New("failed to create IV")

	// ErrCreateNonceFailed indicates failed to create nonce.
	ErrCreateNonceFailed = errors.New("failed to create nonce")

	// ErrExtractIVFailed indicates failed to extract IV.
	ErrExtractIVFailed = errors.New("failed to extract IV")

	// ErrExtractNonceFailed indicates failed to extract nonce.
	ErrExtractNonceFailed = errors.New("failed to extract nonce")

	// ErrSeparateCipherDataFailed indicates failed to separate cipher data.
	ErrSeparateCipherDataFailed = errors.New("failed to separate cipher data")
)
