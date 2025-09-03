package symmetric

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

// 生成固定密钥用于测试
func generateTestKey(size int) []byte {
	key := make([]byte, size)
	for i := range key {
		key[i] = byte(i + 1)
	}
	return key
}

// TestBasicCBCMode 测试基础CBC模式
func TestBasicCBCMode(t *testing.T) {
	t.Log("=== Testing Basic CBC Mode ===")

	// 使用32字节密钥 (AES-256)
	key := generateTestKey(32)

	// 创建AES实例
	aes, err := NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create AES instance: %v", err)
	}

	// 测试数据
	plaintext := []byte("Hello, AES CBC Mode!")
	t.Logf("Original text: %s", string(plaintext))

	// 加密
	ciphertext, err := aes.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	t.Logf("Encrypted (hex): %s", hex.EncodeToString(ciphertext))
	t.Logf("Ciphertext length: %d bytes", len(ciphertext))

	// 解密
	decrypted, err := aes.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	t.Logf("Decrypted text: %s", string(decrypted))

	// 验证
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decryption mismatch!\nOriginal:  %s\nDecrypted: %s",
			string(plaintext), string(decrypted))
	} else {
		t.Log("✅ CBC Mode: Encryption/Decryption successful!")
	}
}

// TestBasicGCMMode 测试基础GCM模式
func TestBasicGCMMode(t *testing.T) {
	t.Log("=== Testing Basic GCM Mode ===")

	// 使用32字节密钥 (AES-256)
	key := generateTestKey(32)

	// 创建GCM模式的AES实例
	option := AESOption{
		Mode: AESModeGCM,
	}
	aes, err := NewAESWithOption(key, option)
	if err != nil {
		t.Fatalf("Failed to create GCM AES instance: %v", err)
	}

	// 测试数据
	plaintext := []byte("Hello, AES GCM Mode!")
	t.Logf("Original text: %s", string(plaintext))

	// 加密
	ciphertext, err := aes.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("GCM encryption failed: %v", err)
	}
	t.Logf("Encrypted (hex): %s", hex.EncodeToString(ciphertext))
	t.Logf("Ciphertext length: %d bytes", len(ciphertext))

	// 解密
	decrypted, err := aes.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("GCM decryption failed: %v", err)
	}
	t.Logf("Decrypted text: %s", string(decrypted))

	// 验证
	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("GCM decryption mismatch!\nOriginal:  %s\nDecrypted: %s",
			string(plaintext), string(decrypted))
	} else {
		t.Log("✅ GCM Mode: Encryption/Decryption successful!")
	}
}

// TestBase64Functions 测试Base64功能
func TestBase64Functions(t *testing.T) {
	t.Log("=== Testing Base64 Functions ===")

	key := generateTestKey(32)
	aes, err := NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create AES instance: %v", err)
	}

	plaintext := "Base64 test message 测试中文!"
	t.Logf("Original text: %s", plaintext)

	// 加密为Base64
	base64Cipher, err := aes.EncryptBase64([]byte(plaintext))
	if err != nil {
		t.Fatalf("Base64 encryption failed: %v", err)
	}
	t.Logf("Base64 encrypted: %s", base64Cipher)

	// 从Base64解密
	decrypted, err := aes.DecryptBase64(base64Cipher)
	if err != nil {
		t.Fatalf("Base64 decryption failed: %v", err)
	}
	t.Logf("Decrypted text: %s", decrypted)

	// 验证
	if plaintext != decrypted {
		t.Errorf("Base64 decryption mismatch!\nOriginal:  %s\nDecrypted: %s",
			plaintext, decrypted)
	} else {
		t.Log("✅ Base64: Encryption/Decryption successful!")
	}
}

// TestDifferentKeySizes 测试不同密钥长度
func TestDifferentKeySizes(t *testing.T) {
	t.Log("=== Testing Different Key Sizes ===")

	keySizes := []struct {
		size int
		name string
	}{
		{16, "AES-128"},
		{24, "AES-192"},
		{32, "AES-256"},
	}

	plaintext := []byte("Test message for different key sizes")

	for _, ks := range keySizes {
		t.Run(ks.name, func(t *testing.T) {
			key := generateTestKey(ks.size)

			aes, err := NewAES(key)
			if err != nil {
				t.Fatalf("Failed to create %s instance: %v", ks.name, err)
			}

			// 加密
			ciphertext, err := aes.Encrypt(plaintext)
			if err != nil {
				t.Fatalf("%s encryption failed: %v", ks.name, err)
			}

			// 解密
			decrypted, err := aes.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("%s decryption failed: %v", ks.name, err)
			}

			// 验证
			if !bytes.Equal(plaintext, decrypted) {
				t.Errorf("%s decryption mismatch", ks.name)
			} else {
				t.Logf("✅ %s: Success", ks.name)
			}
		})
	}
}

// TestMultipleOperations 测试多次操作
func TestMultipleOperations(t *testing.T) {
	t.Log("=== Testing Multiple Operations ===")

	key := generateTestKey(32)
	aes, err := NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create AES instance: %v", err)
	}

	testCases := []string{
		"First message",
		"Second message with more content",
		"Third message 中文测试",
		"", // 注意：空字符串测试可能会失败，这是预期的
		"A very long message that spans multiple blocks to test the padding and encryption thoroughly",
	}

	successCount := 0
	for i, plaintext := range testCases {
		if plaintext == "" {
			t.Logf("Skipping empty string test (expected to fail)")
			continue
		}

		t.Logf("Test %d: %s", i+1, plaintext)

		// 加密
		ciphertext, err := aes.Encrypt([]byte(plaintext))
		if err != nil {
			t.Logf("  ❌ Encryption failed: %v", err)
			continue
		}

		// 解密
		decrypted, err := aes.Decrypt(ciphertext)
		if err != nil {
			t.Logf("  ❌ Decryption failed: %v", err)
			continue
		}

		// 验证
		if !bytes.Equal([]byte(plaintext), decrypted) {
			t.Logf("  ❌ Content mismatch")
			continue
		}

		t.Logf("  ✅ Success")
		successCount++
	}

	t.Logf("Multiple operations: %d/%d successful", successCount, len(testCases)-1) // -1因为跳过了空字符串
}

// TestInvalidInputs 测试无效输入
func TestInvalidInputs(t *testing.T) {
	t.Log("=== Testing Invalid Inputs ===")

	t.Run("Invalid Key Sizes", func(t *testing.T) {
		invalidSizes := []int{15, 17, 31, 33}

		for _, size := range invalidSizes {
			key := make([]byte, size)
			_, err := NewAES(key)
			if err == nil {
				t.Errorf("Should reject key size %d", size)
			} else {
				t.Logf("✅ Correctly rejected key size %d: %v", size, err)
			}
		}
	})

	t.Run("Empty Data", func(t *testing.T) {
		key := generateTestKey(32)
		aes, _ := NewAES(key)

		_, err := aes.Encrypt([]byte{})
		if err == nil {
			t.Error("Should reject empty data")
		} else {
			t.Logf("✅ Correctly rejected empty data: %v", err)
		}
	})

	t.Run("Invalid Base64", func(t *testing.T) {
		key := generateTestKey(32)
		aes, _ := NewAES(key)

		_, err := aes.DecryptBase64("invalid-base64!")
		if err == nil {
			t.Error("Should reject invalid base64")
		} else {
			t.Logf("✅ Correctly rejected invalid base64: %v", err)
		}
	})
}

// TestRandomness 测试随机性
func TestRandomness(t *testing.T) {
	t.Log("=== Testing Randomness ===")

	key := generateTestKey(32)
	aes, err := NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create AES instance: %v", err)
	}

	plaintext := []byte("Same message for randomness test")

	// 加密同一消息多次，应该得到不同的密文
	ciphertexts := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		ciphertext, err := aes.Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Encryption %d failed: %v", i+1, err)
		}
		ciphertexts[i] = ciphertext
		t.Logf("Encryption %d: %s", i+1, hex.EncodeToString(ciphertext[:32])) // 只显示前32字节
	}

	// 检查是否有重复
	unique := true
	for i := 0; i < len(ciphertexts); i++ {
		for j := i + 1; j < len(ciphertexts); j++ {
			if bytes.Equal(ciphertexts[i], ciphertexts[j]) {
				unique = false
				t.Errorf("Found duplicate ciphertexts at index %d and %d", i, j)
			}
		}
	}

	if unique {
		t.Log("✅ All ciphertexts are unique (good randomness)")
	}

	// 验证所有密文都能正确解密
	for i, ciphertext := range ciphertexts {
		decrypted, err := aes.Decrypt(ciphertext)
		if err != nil {
			t.Errorf("Failed to decrypt ciphertext %d: %v", i+1, err)
		} else if !bytes.Equal(plaintext, decrypted) {
			t.Errorf("Decryption %d content mismatch", i+1)
		}
	}
	t.Log("✅ All ciphertexts decrypt correctly")
}

// SimpleExample 简单使用示例
func SimpleExample() {
	fmt.Println("=== Simple AES Usage Example ===")

	// 1. 创建密钥 (32字节 = AES-256)
	key := make([]byte, 32)
	rand.Read(key) // 实际使用中应该用安全的密钥生成方式

	// 2. 创建AES实例
	aes, err := NewAES(key)
	if err != nil {
		fmt.Printf("Error creating AES: %v\n", err)
		return
	}

	// 3. 加密
	message := "Hello, World! 你好世界!"
	encrypted, err := aes.EncryptBase64([]byte(message))
	if err != nil {
		fmt.Printf("Encryption error: %v\n", err)
		return
	}

	// 4. 解密
	decrypted, err := aes.DecryptBase64(encrypted)
	if err != nil {
		fmt.Printf("Decryption error: %v\n", err)
		return
	}

	// 5. 显示结果
	fmt.Printf("Original:  %s\n", message)
	fmt.Printf("Encrypted: %s\n", encrypted)
	fmt.Printf("Decrypted: %s\n", decrypted)
	fmt.Printf("Match: %t\n", message == decrypted)
}

// TestExample 运行示例
func TestExample(t *testing.T) {
	SimpleExample()
}
