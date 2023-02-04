package conf

import (
	"golang.org/x/crypto/bcrypt"
)

func GetDefaultConf() Conf {
	return Conf{
		Log: Log{
			Level:  "warn",
			Stdout: []Stdout{},
			File:   []File{},
		},
		Server: Server{
			Host: "127.0.0.1",
			Port: 7086,
		},
		Auth: Auth{
			User: []User{
				{
					Username: "flydav",
					PasswordHash: (func() string {
						b, _ := bcrypt.GenerateFromPassword([]byte("flydav"), bcrypt.DefaultCost)
						return string(b)
					})(),
					PasswordCrypt: "bcrypt",
				},
			},
		},
	}
}

type Conf struct {
	Log    Log    `toml:"log"`
	Server Server `toml:"server"`
	Auth   Auth   `toml:"auth"`
}

type Server struct {
	Host  string `toml:"host"`
	Port  int    `toml:"port"`
	Path  string `toml:"path"`
	FsDir string `toml:"fs_dir"`
}
type User struct {
	SubPath       string `toml:"sub_path"`
	SubFsDir      string `toml:"sub_fs_dir"`
	Username      string `toml:"username"`
	PasswordHash  string `toml:"password_hash"`
	PasswordCrypt string `toml:"password_crypt"`
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
