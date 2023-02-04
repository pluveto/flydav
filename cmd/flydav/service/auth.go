package service

import (
	"errors"

	"github.com/pluveto/flydav/cmd/flydav/conf"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthService struct {
	Users []conf.User
}

var (
	ErrCrendential = errors.New("invalid username or password")
)

func (s *BasicAuthService) Authenticate(username, password string) error {
	for _, user := range s.Users {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if user.Username == username && err == nil {
			return nil
		}
	}
	return ErrCrendential
}
