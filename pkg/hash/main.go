package hashing

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	hashCharset           = "abcdefghijklmnopqrstuvwxyzA"
	defaultRandomLength   = 16
	passwordCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	defaultPasswordLength = 16
)
const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	passwordLength = 16
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
	allChars := lowercaseChars + uppercaseChars + digitChars + specialChars
	password := make([]rune, passwordLength)

	password[0] = secureRandChar(lowercaseChars)
	password[1] = secureRandChar(uppercaseChars)
	password[2] = secureRandChar(digitChars)
	password[3] = secureRandChar(specialChars)

	for i := 4; i < passwordLength; i++ {
		password[i] = secureRandChar(allChars)
	}

	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

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

func secureRandChar(chars string) rune {
	max := big.NewInt(int64(len(chars)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
		return rune(chars[r.Intn(len(chars))])
	}
	return rune(chars[n.Int64()])
}
