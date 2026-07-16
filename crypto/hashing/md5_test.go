package hashing

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"
)

func TestMd5Hex(t *testing.T) {
	expected := "c4ca4238a0b923820dcc509a6f75849b"
	if actual := Md5Hex("1"); actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestMd5Base64(t *testing.T) {
	sum := md5.Sum([]byte("1"))
	expected := base64.StdEncoding.EncodeToString(sum[:])
	if actual := Md5Base64("1"); actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestMd5Bytes(t *testing.T) {
	expected := md5.Sum([]byte("1"))
	if actual := Md5Bytes([]byte("1")); actual != expected {
		t.Fatalf("expected %x, got %x", expected, actual)
	}
}

func TestMd5File(t *testing.T) {
	filePath := writeHashTestFile(t, "1")
	sum := md5.Sum([]byte("1"))

	expectedHex := hex.EncodeToString(sum[:])
	actualHex, err := Md5FileHex(filePath)
	if err != nil {
		t.Fatalf("md5 file hex: %v", err)
	}
	if actualHex != expectedHex {
		t.Fatalf("expected %s, got %s", expectedHex, actualHex)
	}

	expectedBase64 := base64.StdEncoding.EncodeToString(sum[:])
	actualBase64, err := Md5FileBase64(filePath)
	if err != nil {
		t.Fatalf("md5 file base64: %v", err)
	}
	if actualBase64 != expectedBase64 {
		t.Fatalf("expected %s, got %s", expectedBase64, actualBase64)
	}
}

func writeHashTestFile(t *testing.T, content string) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "hashing-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err = file.WriteString(content); err != nil {
		_ = file.Close()
		t.Fatalf("write temp file: %v", err)
	}
	if err = file.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}
	return file.Name()
}
