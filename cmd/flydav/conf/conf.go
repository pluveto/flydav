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
						b, _ := bcrypt.GenerateFromPassword([]byte("flydav"), bcrypt.DefaultCost)
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
	}
}

type Conf struct {
	Log    Log    `toml:"log"`
	Server Server `toml:"server"`
	Auth   Auth   `toml:"auth"`
	UI     UI     `toml:"ui"`
}

type Server struct {
	Host  string `toml:"host"`
	Port  int    `toml:"port"`
	Path  string `toml:"path"`
	FsDir string `toml:"fs_dir"`
}

type UI struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`   // Path prefix. TODO: ui.path cannot equals to server.path
	Source  string `toml:"source"` // Source location of the UI
}

type User struct {
	SubPath       string      `toml:"sub_path"`
	SubFsDir      string      `toml:"sub_fs_dir"`
	Username      string      `toml:"username"`
	PasswordHash  string      `toml:"password_hash"`
	PasswordCrypt HashMethond `toml:"password_crypt"`
}
type Auth struct {
	User []User `toml:"user"`
}

type File struct {
	Format  LogFormat `toml:"format"`
	Path    string    `toml:"path"`
	MaxSize int       `toml:"max_size"`
	MaxAge  int       `toml:"max_age"`
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
	Format LogFormat `toml:"format"`
	Output LogOutput `toml:"output"`
}
type Log struct {
	Level  string   `toml:"level"`
	File   []File   `toml:"file"`
	Stdout []Stdout `toml:"stdout"`
}
