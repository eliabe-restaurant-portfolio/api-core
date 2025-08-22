package valueobjects

import (
	"crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"
	"strings"
	"time"
	"unicode"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	passwordLength = 16
)

type Password struct {
	password string
}

func NewPassword(pwd string) (Password, error) {
	pwd = strings.TrimSpace(pwd)
	if !isValidPassword(pwd) {
		return Password{}, errors.New("weak password: must be at least 8 characters long and include uppercase, lowercase, digit and special character")
	}
	return Password{password: pwd}, nil
}

func NewRandonPassword() (Password, error) {
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

	return Password{password: string(password)}, nil
}

func (p Password) Get() string {
	return p.password
}

func isValidPassword(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range pwd {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
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
