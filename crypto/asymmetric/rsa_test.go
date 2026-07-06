package asymmetric

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func newTestRsaKeyPair(t *testing.T) RsaKeyPair {
	t.Helper()

	keyPair, err := NewRsaKeyManager(2048).Create()
	if err != nil {
		t.Fatalf("create rsa key pair: %v", err)
	}
	return keyPair
}

func TestRsaPKCS1EncryptDecrypt(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	encrypt := NewRsaEncryptWithPKCS1()
	raw := []byte("hello rsa")

	cipher, err := encrypt.Encrypt(keyPair, raw)
	if err != nil {
		t.Fatalf("encrypt pkcs1: %v", err)
	}
	got, err := encrypt.Decrypt(keyPair, cipher)
	if err != nil {
		t.Fatalf("decrypt pkcs1: %v", err)
	}
	if string(got) != string(raw) {
		t.Fatalf("unexpected decrypt result: got %q, want %q", got, raw)
	}
}

func TestRsaOAEPEncryptDecrypt(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	encrypt, err := NewRsaEncryptWithOAEP(sha1.New(), []byte("label"))
	if err != nil {
		t.Fatalf("create oaep encrypt: %v", err)
	}
	raw := []byte("hello oaep")

	cipher, err := encrypt.Encrypt(keyPair, raw)
	if err != nil {
		t.Fatalf("encrypt oaep: %v", err)
	}
	got, err := encrypt.Decrypt(keyPair, cipher)
	if err != nil {
		t.Fatalf("decrypt oaep: %v", err)
	}
	if string(got) != string(raw) {
		t.Fatalf("unexpected decrypt result: got %q, want %q", got, raw)
	}
}

func TestRsaBase64EncryptDecrypt(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	encrypt := NewRsaEncryptWithPKCS1()
	raw := []byte("hello base64")
	base64Raw := base64.StdEncoding.EncodeToString(raw)

	cipher, err := encrypt.EncryptBase64(keyPair, base64Raw)
	if err != nil {
		t.Fatalf("encrypt base64: %v", err)
	}
	got, err := encrypt.DecryptBase64(keyPair, cipher)
	if err != nil {
		t.Fatalf("decrypt base64: %v", err)
	}
	gotRaw, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("decode decrypted base64: %v", err)
	}
	if string(gotRaw) != string(raw) {
		t.Fatalf("unexpected decrypt result: got %q, want %q", gotRaw, raw)
	}
}

func TestRsaPKCS1SignVerify(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	sign := NewRsaSignWithPKCS1AndSHA256()
	raw := []byte("hello sign")

	signature, err := sign.Sign(keyPair, raw)
	if err != nil {
		t.Fatalf("sign pkcs1: %v", err)
	}
	if err := sign.Verify(keyPair, raw, signature); err != nil {
		t.Fatalf("verify pkcs1: %v", err)
	}
}

func TestRsaPSSSignVerify(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	sign := NewRsaSignWithPSSAndSHA256(rsa.PSSSaltLengthAuto)
	raw := []byte("hello pss")

	signature, err := sign.Sign(keyPair, raw)
	if err != nil {
		t.Fatalf("sign pss: %v", err)
	}
	if err := sign.Verify(keyPair, raw, signature); err != nil {
		t.Fatalf("verify pss: %v", err)
	}
}

func TestRsaPemRoundTrip(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	publicPKCS1, err := keyPair.ToPKCS1PublicPem()
	if err != nil {
		t.Fatalf("export pkcs1 public pem: %v", err)
	}
	privatePKCS1, err := keyPair.ToPKCS1PrivatePem()
	if err != nil {
		t.Fatalf("export pkcs1 private pem: %v", err)
	}
	publicPKIX, err := keyPair.ToPKIXPublicPem()
	if err != nil {
		t.Fatalf("export pkix public pem: %v", err)
	}
	privatePKCS8, err := keyPair.ToPKCS8PrivatePem()
	if err != nil {
		t.Fatalf("export pkcs8 private pem: %v", err)
	}

	manager := NewEmptyRsaKeyManager()
	for name, pair := range map[string][2]string{
		"pkcs1": {publicPKCS1, privatePKCS1},
		"pkcs8": {publicPKIX, privatePKCS8},
	} {
		loaded, err := manager.LoadKeyPair(pair[0], pair[1])
		if err != nil {
			t.Fatalf("load %s key pair: %v", name, err)
		}
		if _, err := NewRsaSignWithPKCS1AndSHA256().Sign(loaded, []byte(name)); err != nil {
			t.Fatalf("use loaded %s key pair: %v", name, err)
		}
	}
}

func TestRsaLoadPrivateKeyDerivesPublicKey(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	privatePem, err := keyPair.ToPKCS8PrivatePem()
	if err != nil {
		t.Fatalf("export private pem: %v", err)
	}

	loaded, err := NewEmptyRsaKeyManager().LoadPrivateKey(privatePem)
	if err != nil {
		t.Fatalf("load private key: %v", err)
	}
	if loaded.PublicKey() == nil {
		t.Fatal("expected public key derived from private key")
	}
}

func TestRsaLoadKeyPairMismatch(t *testing.T) {
	keyPair1 := newTestRsaKeyPair(t)
	keyPair2 := newTestRsaKeyPair(t)
	publicPem, err := keyPair1.ToPKIXPublicPem()
	if err != nil {
		t.Fatalf("export public pem: %v", err)
	}
	privatePem, err := keyPair2.ToPKCS8PrivatePem()
	if err != nil {
		t.Fatalf("export private pem: %v", err)
	}

	_, err = NewEmptyRsaKeyManager().LoadKeyPair(publicPem, privatePem)
	if !errors.Is(err, toolkitError.ErrKeyPairMismatch) {
		t.Fatalf("expected ErrKeyPairMismatch, got %v", err)
	}
}

func TestRsaRejectsWrongKeyType(t *testing.T) {
	ecdsaKeyPair, err := NewEcdsaKeyManager(CreateEcdsaSetting{Curve: elliptic.P256()}).Create()
	if err != nil {
		t.Fatalf("create ecdsa key pair: %v", err)
	}

	_, err = NewRsaEncryptWithPKCS1().Encrypt(ecdsaKeyPair, []byte("raw"))
	if !errors.Is(err, toolkitError.ErrNotRsaPublicKey) {
		t.Fatalf("expected ErrNotRsaPublicKey, got %v", err)
	}
}

func TestRsaOAEPRequiresHashAtRuntime(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	encrypt := &RsaEncrypt{paddingType: PaddingTypeOAEP}

	_, err := encrypt.Encrypt(keyPair, []byte("raw"))
	if !errors.Is(err, toolkitError.ErrNilHashFunction) {
		t.Fatalf("expected ErrNilHashFunction, got %v", err)
	}
}

func TestRsaSignRequiresHash(t *testing.T) {
	keyPair := newTestRsaKeyPair(t)
	sign := &RsaSign{paddingType: PaddingTypePKCS1, hashTypeForSign: crypto.SHA256}

	_, err := sign.Sign(keyPair, []byte("raw"))
	if !errors.Is(err, toolkitError.ErrNilHashFunction) {
		t.Fatalf("expected ErrNilHashFunction, got %v", err)
	}
}
