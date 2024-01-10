package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

func IsStrongPassword(password string) bool {
	minLength := 8
	var number, upper, lower, special, length bool

	for _, char := range password {
		number = number || strings.ContainsRune("0123456789", char)
		upper = upper || strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", char)
		lower = lower || strings.ContainsRune("abcdefghijclmnopqrstuvwxyz", char)
		special = special || strings.ContainsRune("!@#$%^&*()-_=+[]{};:,.<>?/|\\`~", char)
	}

	length = len(password) >= minLength
	return number && upper && lower && special && length
}

var (
	ErrTooShort    = errors.New("password is too short")
	ErrNoUpperCase = errors.New("no uppercase characters")
	ErrNoLowerCase = errors.New("no lowercase characters")
	ErrNoSpecial   = errors.New("no special characters")
	ErrNoNumber    = errors.New("no numbers")
)

func CheckPassword(password string) error {
	minLength := 8
	var number, upper, lower, special, length bool

	for _, char := range password {
		number = number || strings.ContainsRune("0123456789", char)
		upper = upper || strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", char)
		lower = lower || strings.ContainsRune("abcdefghijclmnopqrstuvwxyz", char)
		special = special || strings.ContainsRune("!@#$%^&*()-_=+[]{};:,.<>?/|\\`~", char)
	}

	length = len(password) >= minLength
	if !length {
		return fmt.Errorf("%w: less than %d", ErrTooShort, minLength)
	}
	if !number {
		return ErrNoNumber
	}
	if !upper {
		return ErrNoUpperCase
	}
	if !lower {
		return ErrNoLowerCase
	}
	if !special {
		return ErrNoSpecial
	}
	return nil
}

func Key(password string) []byte {
	key := make([]byte, 32)
	for i := 0; i < len(password); i++ {
		key[i] = password[i]
	}
	return key
}

func Encrypt(str, password string) (string, error) {
	key := Key(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(str))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(str))
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encodedStr, password string) (string, error) {
	key := Key(password)
	ciphertext, err := base64.URLEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

var ErrWrongPassword = errors.New("wrong password")

const prefix = "-=#PREFIX#=-"

func ReliableEncrypt(str, password string) (string, error) {
	return Encrypt(prefix+str, password)
}

func ReliableDecrypt(encodedStr, password string) (string, error) {
	str, err := Decrypt(encodedStr, password)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(str, prefix) {
		return "", ErrWrongPassword
	}
	return str[len(prefix):], nil
}
