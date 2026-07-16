package asymmetric

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func newTestEcdsaKeyPair(t *testing.T) EcdsaKeyPair {
	t.Helper()

	keyPair, err := NewEcdsaKeyManager(elliptic.P256()).Create()
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

func TestEcdsaECPrivatePemRoundTrip(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	privatePem, err := keyPair.ToECPrivatePem()
	if err != nil {
		t.Fatalf("export ec private pem: %v", err)
	}

	loaded, err := NewEmptyEcdsaKeyManager().LoadPrivateKey(privatePem)
	if err != nil {
		t.Fatalf("load ec private key: %v", err)
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

func TestEcdsaLoadPublicKeyOnly(t *testing.T) {
	keyPair := newTestEcdsaKeyPair(t)
	publicPem, err := keyPair.ToPublicPem()
	if err != nil {
		t.Fatalf("export public pem: %v", err)
	}

	loaded, err := NewEmptyEcdsaKeyManager().LoadPublicKey(publicPem)
	if err != nil {
		t.Fatalf("load public key: %v", err)
	}
	if loaded.PublicKey() == nil {
		t.Fatal("expected loaded public key")
	}
	if loaded.PrivateKey() != nil {
		t.Fatal("expected private key to be nil")
	}
	if !sameTestEcdsaPublicKey(keyPair.PublicKey(), loaded.PublicKey()) {
		t.Fatal("expected loaded public key to equal original public key")
	}

	sign := NewEcdsaSign(crypto.SHA256.New())
	raw := []byte("public key only verify")
	signature, err := sign.Sign(keyPair, raw)
	if err != nil {
		t.Fatalf("sign ecdsa: %v", err)
	}
	if err = sign.Verify(loaded, raw, signature); err != nil {
		t.Fatalf("verify with loaded public key: %v", err)
	}
}

func TestEcdsaLoadKeyPairMismatch(t *testing.T) {
	keyPair1 := newTestEcdsaKeyPair(t)
	keyPair2 := newTestEcdsaKeyPair(t)
	publicPem, err := keyPair1.ToPublicPem()
	if err != nil {
		t.Fatalf("export public pem: %v", err)
	}
	privatePem, err := keyPair2.ToPrivatePem()
	if err != nil {
		t.Fatalf("export private pem: %v", err)
	}

	_, err = NewEmptyEcdsaKeyManager().LoadKeyPair(publicPem, privatePem)
	if !errors.Is(err, toolkitError.ErrKeyPairMismatch) {
		t.Fatalf("expected ErrKeyPairMismatch, got %v", err)
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

func TestEcdsaVerifyRejectsWrongKeyType(t *testing.T) {
	ecdsaKeyPair := newTestEcdsaKeyPair(t)
	rsaKeyPair := newTestRsaKeyPair(t)
	sign := NewEcdsaSign(crypto.SHA256.New())
	signature, err := sign.Sign(ecdsaKeyPair, []byte("raw"))
	if err != nil {
		t.Fatalf("sign ecdsa: %v", err)
	}

	err = sign.Verify(rsaKeyPair, []byte("raw"), signature)
	if !errors.Is(err, toolkitError.ErrNotEcdsaPublicKey) {
		t.Fatalf("expected ErrNotEcdsaPublicKey, got %v", err)
	}
}

func sameTestEcdsaPublicKey(expected, actual any) bool {
	expectedKey, ok := expected.(*ecdsa.PublicKey)
	if !ok || expectedKey == nil {
		return false
	}
	actualKey, ok := actual.(*ecdsa.PublicKey)
	if !ok || actualKey == nil {
		return false
	}
	if expectedKey.Curve == nil || actualKey.Curve == nil || expectedKey.X == nil || expectedKey.Y == nil || actualKey.X == nil || actualKey.Y == nil {
		return false
	}
	return expectedKey.Curve == actualKey.Curve && expectedKey.X.Cmp(actualKey.X) == 0 && expectedKey.Y.Cmp(actualKey.Y) == 0
}
