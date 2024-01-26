package authenticator

import (
	"strings"

	"github.com/pluveto/flydav/internal/config"
	"github.com/pluveto/flydav/internal/logger"
)

type Authenticator interface {
	Authenticate(username, password string) (bool, error)
	Authorize(username, path string, permission config.Permission) (bool, error)
	GetRootDir(username string) string
}

// StaticAuthenticator uses static credentials for authentication.
type StaticAuthenticator struct {
	Users map[string]config.StaticUser
	cfg   config.StaticAuthConfig
}

func NewStaticAuthenticator(
	cfg config.StaticAuthConfig,
) *StaticAuthenticator {
	var users = make(map[string]config.StaticUser)
	for _, user := range cfg.Users {
		users[user.Username] = user
	}
	return &StaticAuthenticator{
		Users: users,
		cfg:   cfg,
	}
}

func (sa *StaticAuthenticator) Authenticate(username, password string) (bool, error) {
	// Here you should compare the password with the hashed one in the config.
	// For simplicity, we'll just do a direct comparison.
	user, exists := sa.Users[username]
	if !exists {
		return false, nil
	}
	return user.Password == password, nil
}

func (sa *StaticAuthenticator) Authorize(username, path string, permission config.Permission) (bool, error) {
	user, exists := sa.Users[username]
	logger.Info("authorizing user: ", user, " for path: ", path, " with permission: ", permission)
	if !exists {
		return false, nil
	}
	for _, scope := range user.Scopes {
		if isSubPath(scope.Path, path) {
			for _, perm := range scope.Permissions {
				if perm == permission {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (sa *StaticAuthenticator) GetRootDir(username string) string {
	user, exists := sa.Users[username]
	if !exists {
		return ""
	}
	return user.RootDir
}

func isSubPath(parentPath, subPath string) bool {
	if !strings.HasPrefix(subPath, parentPath) {
		return false
	}
	if len(subPath) == len(parentPath) || subPath[len(parentPath)] == '/' {
		return true
	}
	return false
}
