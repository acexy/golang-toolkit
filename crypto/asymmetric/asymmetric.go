package asymmetric

// KeyPairManager creates and loads typed public/private key pairs.
type KeyPairManager[T KeyPair] interface {
	// Create generates a new public/private key pair using the manager configuration.
	Create() (T, error)
	// LoadPublicKey loads a public key from PEM.
	LoadPublicKey(pubPem string) (T, error)
	// LoadPrivateKey loads a private key from PEM and derives its public key.
	LoadPrivateKey(priPem string) (T, error)
	// LoadKeyPair loads PEM keys and validates their match when both are supplied.
	LoadKeyPair(pubPem, priPem string) (T, error)
}

// KeyPair exposes raw keys and PEM serialization.
type KeyPair interface {
	// PrivateKey returns the raw private key, or nil if it is absent.
	PrivateKey() any
	// PublicKey returns the raw public key, or nil if it is absent.
	PublicKey() any
	// ToPublicPem exports the public key in the default PEM format.
	ToPublicPem() (string, error)
	// ToPrivatePem exports the private key in the default PEM format.
	ToPrivatePem() (string, error)
	// ToPKIXPublicPem exports a PKIX public-key PEM.
	ToPKIXPublicPem() (string, error)
	// ToPKCS8PrivatePem exports a PKCS8 private-key PEM.
	ToPKCS8PrivatePem() (string, error)
}

// RsaKeyPair adds RSA PKCS1 serialization to KeyPair.
type RsaKeyPair interface {
	KeyPair
	// ToPKCS1PublicPem exports an RSA PKCS1 public-key PEM.
	ToPKCS1PublicPem() (string, error)
	// ToPKCS1PrivatePem exports an RSA PKCS1 private-key PEM.
	ToPKCS1PrivatePem() (string, error)
}

// RsaKeyPairManager manages RSA key pairs.
type RsaKeyPairManager interface {
	KeyPairManager[RsaKeyPair]
}

// EcdsaKeyPair adds SEC1 EC private-key serialization to KeyPair.
type EcdsaKeyPair interface {
	KeyPair
	// ToECPrivatePem exports a SEC1 EC private-key PEM.
	ToECPrivatePem() (string, error)
}

// EcdsaKeyPairManager manages ECDSA key pairs.
type EcdsaKeyPairManager interface {
	KeyPairManager[EcdsaKeyPair]
}

// CryptEncrypt defines public-key encryption and private-key decryption.
type CryptEncrypt interface {
	// Encrypt encrypts plaintext bytes with the public key.
	Encrypt(keyPair KeyPair, raw []byte) ([]byte, error)
	// EncryptBase64 decodes Base64 plaintext, encrypts it, and returns Base64 ciphertext.
	EncryptBase64(keyPair KeyPair, base64Raw string) (string, error)
	// Decrypt decrypts ciphertext bytes with the private key.
	Decrypt(keyPair KeyPair, cipher []byte) ([]byte, error)
	// DecryptBase64 decodes Base64 ciphertext, decrypts it, and returns Base64 plaintext.
	DecryptBase64(keyPair KeyPair, base64Cipher string) (string, error)
}

// CryptSign defines private-key signing and public-key verification.
type CryptSign interface {
	// Sign signs data with the private key.
	Sign(keyPair KeyPair, raw []byte) ([]byte, error)
	// Verify verifies data with the public key, returning an error on failure.
	Verify(keyPair KeyPair, raw, sign []byte) error
	// SignBase64 decodes Base64 data and returns its signature as Base64.
	SignBase64(keyPair KeyPair, base64Raw string) (string, error)
	// VerifyBase64 decodes Base64 data and signature before verification.
	VerifyBase64(keyPair KeyPair, base64Raw, base64Sign string) error
}
