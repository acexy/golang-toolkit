package error

import "errors"

var (
	// ErrNilKeyPair 表示公私钥对为空
	ErrNilKeyPair = errors.New("nil key pair")

	// ErrBadKey 表示密钥内容无效
	ErrBadKey = errors.New("bad key")

	// ErrBadKeyLength 表示密钥长度无效
	ErrBadKeyLength = errors.New("bad key length")

	// ErrKeyPairMismatch 表示公钥与私钥不匹配
	ErrKeyPairMismatch = errors.New("key pair mismatch")

	// ErrBadPublicKey 表示公钥内容无效
	ErrBadPublicKey = errors.New("bad public key")

	// ErrBadPrivateKey 表示私钥内容无效
	ErrBadPrivateKey = errors.New("bad private key")

	// ErrNilPublicKey 表示公钥为空
	ErrNilPublicKey = errors.New("nil public key")

	// ErrNilPrivateKey 表示私钥为空
	ErrNilPrivateKey = errors.New("nil private key")

	// ErrNilCurve 表示椭圆曲线为空
	ErrNilCurve = errors.New("nil curve")

	// ErrNilHashFunction 表示 hash 函数为空
	ErrNilHashFunction = errors.New("nil hash function")

	// ErrNotRsaPublicKey 表示不是 RSA 公钥
	ErrNotRsaPublicKey = errors.New("not an RSA public key")

	// ErrNotRsaPrivateKey 表示不是 RSA 私钥
	ErrNotRsaPrivateKey = errors.New("not an RSA private key")

	// ErrNotEcdsaPublicKey 表示不是 ECDSA 公钥
	ErrNotEcdsaPublicKey = errors.New("not an ECDSA public key")

	// ErrNotEcdsaPrivateKey 表示不是 ECDSA 私钥
	ErrNotEcdsaPrivateKey = errors.New("not an ECDSA private key")

	// ErrUnsupportedPaddingType 表示不支持的填充类型
	ErrUnsupportedPaddingType = errors.New("not supported paddingType")

	// ErrVerifyFailed 表示签名验证失败
	ErrVerifyFailed = errors.New("verify failed")
)
