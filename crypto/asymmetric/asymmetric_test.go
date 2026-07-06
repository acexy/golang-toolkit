package asymmetric

import (
	"crypto"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func newTestEcdsaKeyPair(t *testing.T) EcdsaKeyPair {
	t.Helper()

	keyPair, err := NewEcdsaKeyManager(CreateEcdsaSetting{Curve: elliptic.P256()}).Create()
	if err != nil {
		t.Fatalf("create ecdsa key pair: %v", err)
	}
	return keyPair
}

func TestEcdsaPemRoundTrip(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	publicPem, err := keyPair.ToPublicPem()
	if err != nil {
		t.Fatalf("export public pem: %v", err)
	}
	privatePem, err := keyPair.ToPrivatePem()
	if err != nil {
		t.Fatalf("export private pem: %v", err)
	}

	loaded, err := NewEmptyEcdsaKeyManager().LoadKeyPair(publicPem, privatePem)
	if err != nil {
		t.Fatalf("load ecdsa key pair: %v", err)
	}
	if loaded.PublicKey() == nil || loaded.PrivateKey() == nil {
		t.Fatal("expected loaded public and private keys")
	}
}

func TestEcdsaLoadPrivateKeyDerivesPublicKey(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	privatePem, err := keyPair.ToPrivatePem()
	if err != nil {
		t.Fatalf("export private pem: %v", err)
	}

	loaded, err := NewEmptyEcdsaKeyManager().LoadPrivateKey(privatePem)
	if err != nil {
		t.Fatalf("load private key: %v", err)
	}
	if loaded.PublicKey() == nil {
		t.Fatal("expected public key derived from private key")
	}
}

func TestEcdsaSignVerifyASN1(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())
	raw := []byte("hello ecdsa")

	signature, err := sign.Sign(keyPair, raw)
	if err != nil {
		t.Fatalf("sign ecdsa asn1: %v", err)
	}
	if err := sign.Verify(keyPair, raw, signature); err != nil {
		t.Fatalf("verify ecdsa asn1: %v", err)
	}
}

func TestEcdsaSignVerifyRawDigest(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	sign := NewEcdsaSign(nil)
	raw := []byte("raw digest")

	signature, err := sign.Sign(keyPair, raw)
	if err != nil {
		t.Fatalf("sign raw digest: %v", err)
	}
	if err := sign.Verify(keyPair, raw, signature); err != nil {
		t.Fatalf("verify raw digest: %v", err)
	}
}

func TestEcdsaSignVerifyRS(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())
	raw := []byte("hello rs")

	r, s, err := sign.SignRS(keyPair, raw)
	if err != nil {
		t.Fatalf("sign rs: %v", err)
	}
	ok, err := sign.VerifyRS(keyPair, raw, r, s)
	if err != nil {
		t.Fatalf("verify rs: %v", err)
	}
	if !ok {
		t.Fatal("expected rs signature verification success")
	}
}

func TestEcdsaBase64SignVerify(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())
	base64Raw := base64.StdEncoding.EncodeToString([]byte("hello base64"))

	signature, err := sign.SignBase64(keyPair, base64Raw)
	if err != nil {
		t.Fatalf("sign base64: %v", err)
	}
	if err := sign.VerifyBase64(keyPair, base64Raw, signature); err != nil {
		t.Fatalf("verify base64: %v", err)
	}
}

func TestEcdsaVerifyFailure(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())
	signature, err := sign.Sign(keyPair, []byte("raw"))
	if err != nil {
		t.Fatalf("sign ecdsa: %v", err)
	}

	err = sign.Verify(keyPair, []byte("changed"), signature)
	if !errors.Is(err, toolkitError.ErrVerifyFailed) {
		t.Fatalf("expected ErrVerifyFailed, got %v", err)
	}
}

func TestEcdsaRejectsWrongKeyType(t *testing.T) {
	rsaKeyPair := newTestRsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())

	_, err := sign.Sign(rsaKeyPair, []byte("raw"))
	if !errors.Is(err, toolkitError.ErrNotEcdsaPrivateKey) {
		t.Fatalf("expected ErrNotEcdsaPrivateKey, got %v", err)
	}
}
