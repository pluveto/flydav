package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pluveto/flydav/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type HashMethond string

const BcryptHash HashMethond = "bcrypt"
const SHA256Hash HashMethond = "sha256"

func GetDefaultConf() Conf {
	defaultFsDir, _ := os.Getwd()
	if !strings.HasPrefix(defaultFsDir, "/home") {
		webdavDir := filepath.Join(os.TempDir(), "flydav")
		err := os.MkdirAll(webdavDir, 0755)
		if err != nil {
			logger.Fatal("Failed to create webdav tmp dir", err)
		}
		defaultFsDir = webdavDir
	}
	return Conf{
		Log: Log{
			Level:  "warn",
			Stdout: []Stdout{},
			File:   []File{},
		},
		Server: Server{
			Host:  "127.0.0.1",
			Port:  7086,
			Path:  "/webdav",
			FsDir: defaultFsDir,
		},
		Auth: Auth{
			User: []User{
				{
					Username: "flydav",
					PasswordHash: (func() string {
						b, _ := bcrypt.GenerateFromPassword([]byte("flydavdefaultpassword"), bcrypt.DefaultCost)
						return string(b)
					})(),
					PasswordCrypt: BcryptHash,
				},
			},
		},
		UI: UI{
			Enabled: false,
			Path:    "/ui",
			Source:  "",
		},
		CORS: CORS{
			Enabled: false,
		},
	}
}

type Conf struct {
	Log    Log    `toml:"log" yaml:"log"`
	Server Server `toml:"server" yaml:"server"`
	Auth   Auth   `toml:"auth" yaml:"auth"`
	UI     UI     `toml:"ui" yaml:"ui"`
	CORS   CORS   `toml:"cors" yaml:"cors"`
}

type CORS struct {
	Enabled          bool     `toml:"enabled" yaml:"enabled"`
	AllowedOrigins   []string `toml:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `toml:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `toml:"allowed_headers" yaml:"allowed_headers"`
	ExposedHeaders   []string `toml:"exposed_headers" yaml:"exposed_headers"`
	AllowCredentials bool     `toml:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           int
}

type Server struct {
	Host  string `toml:"host" yaml:"host"`
	Port  int    `toml:"port" yaml:"port"`
	Path  string `toml:"path" yaml:"path"`
	FsDir string `toml:"fs_dir" yaml:"fs_dir"`
}

type UI struct {
	Enabled bool   `toml:"enabled" yaml:"enabled"`
	Path    string `toml:"path" yaml:"path"`   // Path prefix. TODO: ui.path cannot equals to server.path
	Source  string `toml:"source" yaml:"source"` // Source location of the UI
}

type User struct {
	SubPath       string      `toml:"sub_path" yaml:"sub_path"`
	SubFsDir      string      `toml:"sub_fs_dir" yaml:"sub_fs_dir"`
	Username      string      `toml:"username" yaml:"username"`
	PasswordHash  string      `toml:"password_hash" yaml:"password_hash"`
	PasswordCrypt HashMethond `toml:"password_crypt" yaml:"password_crypt"`
}
type Auth struct {
	User []User `toml:"user" yaml:"user"`
}

type File struct {
	Format  LogFormat `toml:"format" yaml:"format"`
	Path    string    `toml:"path" yaml:"path"`
	MaxSize int       `toml:"max_size" yaml:"max_size"`
	MaxAge  int       `toml:"max_age" yaml:"max_age"`
}

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

type LogOutput string

const (
	LogOutputStdout LogOutput = "stdout"
	LogOutputStderr LogOutput = "stderr"
	LogOutputFile   LogOutput = "file"
)

type Stdout struct {
	Format LogFormat `toml:"format" yaml:"format"`
	Output LogOutput `toml:"output" yaml:"output"`
}
type Log struct {
	Level  string   `toml:"level" yaml:"level"`
	File   []File   `toml:"file" yaml:"file"`
	Stdout []Stdout `toml:"stdout" yaml:"stdout"`
}
