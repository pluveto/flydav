package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/pluveto/flydav/cmd/flydav/conf"
	"github.com/pluveto/flydav/pkg/logger"
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
		logger.Debug("no such user: ", username)
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
		logger.Debug("bcrypt compare error: ", err)
	case conf.SHA256Hash:
		gen := sha256.New()
		gen.Write([]byte(password))
		expectedHash := hex.EncodeToString(gen.Sum(nil))
		if user.Username == username && user.PasswordHash == expectedHash {
			return nil
		}
		logger.Debug("sha256 compare error, expected hash: ", user.PasswordHash, ", actual hash: ", expectedHash)
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
