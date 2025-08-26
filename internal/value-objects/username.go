package valueobjects

import (
	"strings"
)

type Username struct {
	username string
}

func NewUsername(username string) (Username, error) {
	username = strings.TrimSpace(username)
	return Username{username: username}, nil
}

func (u Username) Get() string {
	return u.username
}
