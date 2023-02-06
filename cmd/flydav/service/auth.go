package service

import (
	"crypto/sha256"
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
	ErrCrendential           = errors.New("invalid username or password")
	ErrUnsupportedHashMethod = errors.New("unsupported hash method")
)

func (s *BasicAuthService) Authenticate(username, password string) error {
	user, ok := s.UserMap[username]
	if !ok {
		return ErrCrendential
	}
	hashMethod := user.PasswordCrypt
	// todo: sha256
	switch hashMethod {
	case conf.BcryptHash:
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if user.Username == username && err == nil {
			return nil
		}
	case conf.SHA256Hash:
		sha256 := sha256.New()
		sha256.Write([]byte(password))
		if user.Username == username && user.PasswordHash == string(sha256.Sum(nil)) {
			return nil
		}
	default:
		return ErrUnsupportedHashMethod
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
