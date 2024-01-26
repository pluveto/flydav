package config

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
)

type UserScope struct {
	Path        string       `yaml:"path"`
	Permissions []Permission `yaml:"permissions"`
}

type StaticUser struct {
	Username string      `yaml:"username"`
	Password string      `yaml:"password"`
	RootDir  string      `yaml:"root_dir"`
	Scopes   []UserScope `yaml:"scopes"`
}

type StaticAuthConfig struct {
	Enabled   bool         `yaml:"enabled"`
	Superuser StaticUser   `yaml:"superuser"`
	Users     []StaticUser `yaml:"users"`
}

type LdapAuthConfig struct {
	Enabled bool `yaml:"enabled"`
}

type LocalAuthConfig struct {
	Enabled bool `yaml:"enabled"`
}

type AuthConfig struct {
	Path     string    `yaml:"path"`
	Log      LogConfig `yaml:"log"`
	Backends struct {
		Ldap   LdapAuthConfig   `yaml:"ldap"`
		Local  LocalAuthConfig  `yaml:"local"`
		Static StaticAuthConfig `yaml:"static"`
	} `yaml:"backends"`
}

type AuthBackend string

var LdapAuthBackend AuthBackend = "ldap"
var LocalAuthBackend AuthBackend = "local"
var StaticAuthBackend AuthBackend = "static"

func (ac *AuthConfig) GetEnabledAuthBackend() AuthBackend {
	if ac.Backends.Ldap.Enabled {
		return LdapAuthBackend
	}
	if ac.Backends.Local.Enabled {
		return LocalAuthBackend
	}
	if ac.Backends.Static.Enabled {
		return StaticAuthBackend
	}
	return ""
}
