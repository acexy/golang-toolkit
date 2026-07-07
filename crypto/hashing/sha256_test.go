package hashing

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSha256Hex(t *testing.T) {
	sum := sha256.Sum256([]byte("sha256"))
	expected := hex.EncodeToString(sum[:])
	if actual := Sha256Hex("sha256"); actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestSha256Base64(t *testing.T) {
	sum := sha256.Sum256([]byte("sha256"))
	expected := base64.StdEncoding.EncodeToString(sum[:])
	if actual := Sha256Base64("sha256"); actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestSha256Bytes(t *testing.T) {
	sum := sha256.Sum256([]byte("sha256"))
	expected := sum[:]
	actual := Sha256Bytes([]byte("sha256"))
	if hex.EncodeToString(actual) != hex.EncodeToString(expected) {
		t.Fatalf("expected %x, got %x", expected, actual)
	}
}

func TestSha256File(t *testing.T) {
	filePath := writeHashTestFile(t, "sha256")
	sum := sha256.Sum256([]byte("sha256"))

	expectedHex := hex.EncodeToString(sum[:])
	actualHex, err := Sha256FileHex(filePath)
	if err != nil {
		t.Fatalf("sha256 file hex: %v", err)
	}
	if actualHex != expectedHex {
		t.Fatalf("expected %s, got %s", expectedHex, actualHex)
	}

	expectedBase64 := base64.StdEncoding.EncodeToString(sum[:])
	actualBase64, err := Sha256FileBase64(filePath)
	if err != nil {
		t.Fatalf("sha256 file base64: %v", err)
	}
	if actualBase64 != expectedBase64 {
		t.Fatalf("expected %s, got %s", expectedBase64, actualBase64)
	}
}
