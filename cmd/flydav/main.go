package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/alexflint/go-arg"
	"github.com/pluveto/flydav/cmd/flydav/app"
	"github.com/pluveto/flydav/cmd/flydav/conf"
	"github.com/pluveto/flydav/pkg/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// main Entry point of the application
func main() {
	// main Entry point of the application
	var args app.Args
	var cnf conf.Conf
	var defaultConf = conf.GetDefaultConf()

	args = loadArgsValid()
	cnf = loadConfValid(args.Config, defaultConf, "config.toml")
	overrideConf(&cnf, args)
	validateConf(&cnf)
	app.InitLogger(cnf.Log, args.Verbose)
	logger.Debug("log level: ", logger.GetLevel())
	app.Run(cnf)
}

func validateConf(conf *conf.Conf) {
	if len(conf.Auth.User) == 0 {
		logger.Fatal("No user configured")
	}
	if conf.Auth.User[0].Username == "" {
		logger.Fatal("No username configured")
	}
	if conf.Auth.User[0].PasswordHash == "" {
		logger.Fatal("No password configured")
	}
}

func overrideConf(cnf *conf.Conf, args app.Args) {
	if args.Verbose {
		cnf.Log.Level = logrus.DebugLevel.String()
	}
	if args.Port != 0 {
		cnf.Server.Port = args.Port
	}
	if args.Host != "" {
		cnf.Server.Host = args.Host
	}
	if args.Username == "" {
		args.Username = "flydav"
	}
	if args.Config == "" {
		cnf.Auth.User = []conf.User{
			{
				Username:      args.Username,
				PasswordHash:  promptPassword(args.Username),
				PasswordCrypt: "bcrypt",
			},
		}
	}
}

func promptPassword(username string) string {
	var err error
	var password []byte
	MIN_PASS_LEN := 9
	fmt.Printf("Set a temporary password for user %s (at least %d chars): ", username, MIN_PASS_LEN)
	for {
		password, err = term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err == nil && len(password) >= MIN_PASS_LEN {
			break
		}
		fmt.Printf("Invalid password. Must be at least %d chars. Try agin: ", MIN_PASS_LEN)
	}
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}

func loadArgsValid() app.Args {
	var args app.Args
	arg.MustParse(&args)
	return args
}

func getAppDir() string {
	dir, err := os.Executable()
	if err != nil {
		logger.Fatal(err)
	}
	return filepath.Dir(dir)
}

func loadConfValid(path string, defaultConf conf.Conf, defaultConfPath string) conf.Conf {
	if path == "" {
		path = defaultConfPath
	}
	// app executable dir + config.toml has the highest priority
	preferredPath := filepath.Join(getAppDir(), path)
	if _, err := os.Stat(preferredPath); err == nil {
		path = preferredPath
	}
	_, err := toml.DecodeFile(path, &defaultConf)
	if err != nil {
		logger.Warn("failed to load config file: ", err, " using default config")
	}
	logger.WithField("conf", &defaultConf).Debug("configuration loaded")
	return defaultConf
}
