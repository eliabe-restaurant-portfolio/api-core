package valueobjects

import (
	"errors"
	"strings"
)

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if !isValidEmail(email) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{
		email: email,
	}, nil
}

func (e Email) Get() string {
	return e.email
}

func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 || at == len(email)-1 {
		return false
	}
	dot := strings.LastIndex(email[at:], ".")
	return dot >= 2
}
