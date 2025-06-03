package symmetric

import (
	"crypto/aes"
	"fmt"
	"testing"
)

type AESIVCreator struct {
}

func (A AESIVCreator) Encrypt(key, raw []byte) [aes.BlockSize]byte {
	return [16]byte(key[:aes.BlockSize])
}

func (A AESIVCreator) Decrypt(key, cipherText []byte) [aes.BlockSize]byte {
	return [16]byte(key[:aes.BlockSize])
}

type AESResultCreator struct {
}

func (A AESResultCreator) Encrypt(iv [16]byte, rawCipherData []byte) []byte {
	return rawCipherData
}

func (A AESResultCreator) Decrypt(iv [16]byte, cipherData []byte) []byte {
	return cipherData
}

func TestAESEncrypt(t *testing.T) {
	key := []byte("1234567890abcdef")     // 16字节key
	raw := []byte("hello aes12345678 明文") // 16字节明文

	// 默认AES工作模式
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

	// 自定义AES工作模式
	encrypt = NewAESWithOption(key, AESOption{
		IVCreator:     &AESIVCreator{},
		ResultCreator: AESResultCreator{},
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
