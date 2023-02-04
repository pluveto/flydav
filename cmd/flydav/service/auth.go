package service

import (
	"errors"

	"github.com/pluveto/flydav/cmd/flydav/conf"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthService struct {
	UserMap map[string]conf.User
}

func NewBasicAuthService(users []conf.User) *BasicAuthService {
	ret := &BasicAuthService{}
	ret.UserMap = make(map[string]conf.User)
	for _, user := range users {
		ret.UserMap[user.Username] = user
	}
	return ret
}

var (
	ErrCrendential = errors.New("invalid username or password")
)

func (s *BasicAuthService) Authenticate(username, password string) error {
	user, ok := s.UserMap[username]
	if !ok {
		return ErrCrendential
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if user.Username == username && err == nil {
		return nil
	}
	return ErrCrendential
}

func (s *BasicAuthService) GetAuthorizedSubDir(username string) (string, error) {
	user, ok := s.UserMap[username]
	if !ok {
		return "", errors.New("no such user")
	}
	return user.SubFsDir, nil
}
func (s *BasicAuthService) GetPathPrefix(username string) (string, error) {
	user, ok := s.UserMap[username]
	if !ok {
		return "", errors.New("no such user")
	}
	return user.SubPath, nil
}
