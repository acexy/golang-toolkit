package error

import "errors"

var (
	// ErrNilKeyPair indicates nil key pair.
	ErrNilKeyPair = errors.New("nil key pair")

	// ErrBadKey indicates bad key.
	ErrBadKey = errors.New("bad key")

	// ErrBadKeyLength indicates bad key length.
	ErrBadKeyLength = errors.New("bad key length")

	// ErrKeyPairMismatch indicates key pair mismatch.
	ErrKeyPairMismatch = errors.New("key pair mismatch")

	// ErrBadPublicKey indicates bad public key.
	ErrBadPublicKey = errors.New("bad public key")

	// ErrBadPrivateKey indicates bad private key.
	ErrBadPrivateKey = errors.New("bad private key")

	// ErrUnsupportedPublicKeyType indicates unsupported public key type.
	ErrUnsupportedPublicKeyType = errors.New("unsupported public key type")

	// ErrUnsupportedPrivateKeyType indicates unsupported private key type.
	ErrUnsupportedPrivateKeyType = errors.New("unsupported private key type")

	// ErrNilPublicKey indicates nil public key.
	ErrNilPublicKey = errors.New("nil public key")

	// ErrNilPrivateKey indicates nil private key.
	ErrNilPrivateKey = errors.New("nil private key")

	// ErrNilCurve indicates nil curve.
	ErrNilCurve = errors.New("nil curve")

	// ErrNilHashFunction indicates nil hash function.
	ErrNilHashFunction = errors.New("nil hash function")

	// ErrNotRsaPublicKey indicates not an RSA public key.
	ErrNotRsaPublicKey = errors.New("not an RSA public key")

	// ErrNotRsaPrivateKey indicates not an RSA private key.
	ErrNotRsaPrivateKey = errors.New("not an RSA private key")

	// ErrNotEcdsaPublicKey indicates not an ECDSA public key.
	ErrNotEcdsaPublicKey = errors.New("not an ECDSA public key")

	// ErrNotEcdsaPrivateKey indicates not an ECDSA private key.
	ErrNotEcdsaPrivateKey = errors.New("not an ECDSA private key")

	// ErrUnsupportedPaddingType indicates not supported paddingType.
	ErrUnsupportedPaddingType = errors.New("not supported paddingType")

	// ErrVerifyFailed indicates verify failed.
	ErrVerifyFailed = errors.New("verify failed")
)
