package symmetric

import (
	"bytes"
	"crypto/aes"
	"errors"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

func generateTestKey(size int) []byte {
	key := make([]byte, size)
	for i := range key {
		key[i] = byte(i + 1)
	}
	return key
}

func TestAESCBCRoundTrip(t *testing.T) {
	crypt, err := NewAES(generateTestKey(32))
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	raw := []byte("hello aes cbc")
	cipherData, err := crypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("encrypt cbc: %v", err)
	}
	decrypted, err := crypt.Decrypt(cipherData)
	if err != nil {
		t.Fatalf("decrypt cbc: %v", err)
	}
	if !bytes.Equal(raw, decrypted) {
		t.Fatalf("expected %q, got %q", raw, decrypted)
	}
}

func TestAESGCMRoundTrip(t *testing.T) {
	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{Mode: AESModeGCM})
	if err != nil {
		t.Fatalf("new aes gcm: %v", err)
	}

	raw := []byte("hello aes gcm")
	cipherData, err := crypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("encrypt gcm: %v", err)
	}
	decrypted, err := crypt.Decrypt(cipherData)
	if err != nil {
		t.Fatalf("decrypt gcm: %v", err)
	}
	if !bytes.Equal(raw, decrypted) {
		t.Fatalf("expected %q, got %q", raw, decrypted)
	}
}

func TestAESBase64RoundTrip(t *testing.T) {
	crypt, err := NewAES(generateTestKey(32))
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	raw := "base64 message"
	cipherText, err := crypt.EncryptBase64([]byte(raw))
	if err != nil {
		t.Fatalf("encrypt base64: %v", err)
	}
	decrypted, err := crypt.DecryptBase64(cipherText)
	if err != nil {
		t.Fatalf("decrypt base64: %v", err)
	}
	if decrypted != raw {
		t.Fatalf("expected %q, got %q", raw, decrypted)
	}
}

func TestAESKeySizes(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		if _, err := NewAES(generateTestKey(size)); err != nil {
			t.Fatalf("expected key size %d to be valid: %v", size, err)
		}
	}
	for _, size := range []int{0, 15, 17, 31, 33} {
		_, err := NewAES(generateTestKey(size))
		if !errors.Is(err, toolkitError.ErrInvalidAESKeySize) {
			t.Fatalf("expected ErrInvalidAESKeySize for %d, got %v", size, err)
		}
	}
}

func TestAESRejectsUnsupportedMode(t *testing.T) {
	_, err := NewAESWithOption(generateTestKey(32), AESOption{Mode: AESMode(99)})
	if !errors.Is(err, toolkitError.ErrUnsupportedAESMode) {
		t.Fatalf("expected ErrUnsupportedAESMode, got %v", err)
	}
}

func TestAESErrors(t *testing.T) {
	crypt, err := NewAES(generateTestKey(32))
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	_, err = crypt.Encrypt(nil)
	if !errors.Is(err, toolkitError.ErrEmptyEncryptData) {
		t.Fatalf("expected ErrEmptyEncryptData, got %v", err)
	}

	_, err = crypt.Decrypt(nil)
	if !errors.Is(err, toolkitError.ErrEmptyCipherData) {
		t.Fatalf("expected ErrEmptyCipherData, got %v", err)
	}

	_, err = crypt.DecryptBase64("invalid-base64!")
	if !errors.Is(err, toolkitError.ErrInvalidBase64Data) {
		t.Fatalf("expected ErrInvalidBase64Data, got %v", err)
	}

	crypt.mode = AESMode(99)
	_, err = crypt.Encrypt([]byte("raw"))
	if !errors.Is(err, toolkitError.ErrUnsupportedAESMode) {
		t.Fatalf("expected ErrUnsupportedAESMode, got %v", err)
	}
}

func TestAESCBCRejectsInvalidIVSize(t *testing.T) {
	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{
		IVCreator: &fixedIVCreator{iv: []byte("short")},
	})
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	_, err = crypt.Encrypt([]byte("raw"))
	if !errors.Is(err, toolkitError.ErrInvalidIVSize) {
		t.Fatalf("expected ErrInvalidIVSize, got %v", err)
	}
}

func TestAESGCMRejectsInvalidNonceSize(t *testing.T) {
	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{
		Mode:      AESModeGCM,
		IVCreator: &fixedIVCreator{iv: []byte("short")},
	})
	if err != nil {
		t.Fatalf("new aes gcm: %v", err)
	}

	_, err = crypt.Encrypt([]byte("raw"))
	if !errors.Is(err, toolkitError.ErrInvalidNonceSize) {
		t.Fatalf("expected ErrInvalidNonceSize, got %v", err)
	}
}

func TestAESGCMRejectsTamperedCipherData(t *testing.T) {
	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{Mode: AESModeGCM})
	if err != nil {
		t.Fatalf("new aes gcm: %v", err)
	}

	cipherData, err := crypt.Encrypt([]byte("authenticated data"))
	if err != nil {
		t.Fatalf("encrypt gcm: %v", err)
	}
	cipherData[len(cipherData)-1] ^= 1
	if _, err = crypt.Decrypt(cipherData); err == nil {
		t.Fatal("expected tampered GCM cipher data to fail")
	}
}

func TestPkcs7PaddingRejectsInvalidContent(t *testing.T) {
	padding := &Pkcs7PaddingCreator{}

	_, err := padding.Pad([]byte("raw"), 0)
	if !errors.Is(err, toolkitError.ErrInvalidBlockSize) {
		t.Fatalf("expected ErrInvalidBlockSize, got %v", err)
	}

	_, err = padding.UnPad([]byte{1, 2, 3, 4})
	if !errors.Is(err, toolkitError.ErrInvalidPadding) {
		t.Fatalf("expected ErrInvalidPadding, got %v", err)
	}
}

func TestAESCustomComponents(t *testing.T) {
	fixedIV := make([]byte, aes.BlockSize)
	for i := range fixedIV {
		fixedIV[i] = byte(i)
	}

	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{
		IVCreator:      &fixedIVCreator{iv: fixedIV},
		ResultCreator:  &reverseAppendResultCreator{},
		PaddingCreator: &zeroPaddingCreator{},
	})
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	raw := []byte("custom components")
	cipherData, err := crypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("encrypt custom: %v", err)
	}
	decrypted, err := crypt.Decrypt(cipherData)
	if err != nil {
		t.Fatalf("decrypt custom: %v", err)
	}
	if !bytes.Equal(raw, decrypted) {
		t.Fatalf("expected %q, got %q", raw, decrypted)
	}
}

func TestAESPureResultCreator(t *testing.T) {
	fixedIV := make([]byte, aes.BlockSize)
	for i := range fixedIV {
		fixedIV[i] = byte(i)
	}

	crypt, err := NewAESWithOption(generateTestKey(32), AESOption{
		IVCreator:     &fixedIVCreator{iv: fixedIV},
		ResultCreator: &PureResultCreator{},
	})
	if err != nil {
		t.Fatalf("new aes: %v", err)
	}

	raw := []byte("pure result")
	cipherData, err := crypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("encrypt pure: %v", err)
	}
	decrypted, err := crypt.Decrypt(cipherData)
	if err != nil {
		t.Fatalf("decrypt pure: %v", err)
	}
	if !bytes.Equal(raw, decrypted) {
		t.Fatalf("expected %q, got %q", raw, decrypted)
	}
}

type fixedIVCreator struct {
	iv []byte
}

func (f *fixedIVCreator) CreateForEncrypt(key, rawData []byte) ([]byte, error) {
	return f.iv, nil
}

func (f *fixedIVCreator) ExtractForDecrypt(key, cipherData []byte) ([]byte, error) {
	return f.iv, nil
}

type reverseAppendResultCreator struct{}

func (r *reverseAppendResultCreator) CombineResult(iv, cipherData []byte) []byte {
	result := make([]byte, len(iv)+len(cipherData))
	copy(result, cipherData)
	copy(result[len(cipherData):], iv)
	return result
}

func (r *reverseAppendResultCreator) SeparateResult(combinedData []byte, ivSize int) (cipherData []byte, err error) {
	if len(combinedData) < ivSize {
		return nil, toolkitError.ErrCipherDataTooShort
	}
	return combinedData[:len(combinedData)-ivSize], nil
}

type zeroPaddingCreator struct{}

func (z *zeroPaddingCreator) Pad(rawData []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, toolkitError.ErrInvalidBlockSize
	}

	padding := blockSize - len(rawData)%blockSize
	if padding == blockSize {
		padding = 0
	}
	result := make([]byte, len(rawData)+padding)
	copy(result, rawData)
	return result, nil
}

func (z *zeroPaddingCreator) UnPad(paddedData []byte) ([]byte, error) {
	for i := len(paddedData) - 1; i >= 0; i-- {
		if paddedData[i] != 0 {
			return paddedData[:i+1], nil
		}
	}
	return []byte{}, nil
}
