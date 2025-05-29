package symmetric

import (
	"fmt"
	"testing"
)

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
}
