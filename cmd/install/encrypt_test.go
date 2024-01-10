package main

import "testing"

func TestEncrypt(t *testing.T) {
	str := "test string"
	password := "testPasswordtestPassword"
	_, err := Encrypt(str, password)
	if err != nil {
		t.Fatalf("Failed to encrypt: %s", err)
	}
}

func TestDecrypt(t *testing.T) {
	str := "test string"
	password := "testPasswordtestPassword"
	encryptedStr, _ := Encrypt(str, password)
	t.Logf("len: %d", len(encryptedStr))
	decryptedStr, err := Decrypt(encryptedStr, password)
	if err != nil {
		t.Fatalf("Failed to decrypt: %s", err)
	}
	if decryptedStr != str {
		t.Errorf("Expected: %v, got: %v", str, decryptedStr)
	}
}
