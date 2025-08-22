package hashing

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	hashCharset           = "abcdefghijklmnopqrstuvwxyzA"
	defaultRandomLength   = 16
	passwordCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	defaultPasswordLength = 16
)

func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Compare(first, second string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(first), []byte(second))
	return err == nil
}

func GeneratePassword() (string, string, error) {
	password := make([]byte, defaultPasswordLength)
	charsetLen := big.NewInt(int64(len(passwordCharset)))

	for i := 0; i < defaultPasswordLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate random index: %w", err)
		}
		password[i] = passwordCharset[randomIndex.Int64()]
	}

	encrypted, err := Hash(string(password))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random index: %w", err)
	}

	return encrypted, string(password), nil
}

func GenerateRandom() (string, string, error) {
	random := make([]byte, defaultRandomLength)
	charsetLen := big.NewInt(int64(len(hashCharset)))

	for i := 0; i < defaultRandomLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate random index: %w", err)
		}
		random[i] = hashCharset[randomIndex.Int64()]
	}

	encrypted, err := Hash(string(random))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate random index: %w", err)
	}

	return encrypted, string(random), nil
}
