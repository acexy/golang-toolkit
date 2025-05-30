package symmetric

import (
	"crypto/aes"
	"fmt"
	"testing"
)

type AESIVCreator struct {
}

func (A AESIVCreator) Encrypt(key, raw []byte) []byte {
	return key[:aes.BlockSize]
}

func (A AESIVCreator) Decrypt(key, cipherText []byte) []byte {
	return key[:aes.BlockSize]
}

func TestAESEncrypt(t *testing.T) {
	key := []byte("1234567890abcdef")     // 16字节key
	raw := []byte("hello aes12345678 明文") // 16字节明文
	encrypt := NewAES(key)
	enc, err := encrypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("AES Encrypt error: %v", err)
	}
	dec, err := encrypt.Decrypt(enc)
	if err != nil {
		t.Fatalf("AES Decrypt error: %v", err)
	}
	fmt.Println(string(dec))

	// 自定义IV
	encrypt = NewAESWithOption(key, AESOption{
		IVCreator: &AESIVCreator{},
	})
	enc, err = encrypt.Encrypt(raw)
	if err != nil {
		t.Fatalf("AES Encrypt error: %v", err)
	}
	dec, err = encrypt.Decrypt(enc)
	if err != nil {
		t.Fatalf("AES Decrypt error: %v", err)
	}
	fmt.Println(string(dec))
}
