package valueobjects

import (
	"errors"
	"strings"
	"unicode"
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
