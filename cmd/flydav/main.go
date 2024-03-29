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
	"github.com/pluveto/flydav/pkg/misc"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// main Entry point of the application
func main() {
	// main Entry point of the application
	var args app.Args
	var cnf conf.Conf
	var defaultConf = conf.GetDefaultConf()

	args = loadArgsValid()
	if args.Verbose {
		fmt.Printf("args: %+v\n", args)
	}

	cnf = loadConfValid(args.Verbose, args.Config, defaultConf, "config.toml")
	if args.Verbose {
		fmt.Printf("conf: %+v\n", cnf)
	}
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
	if args.EnabledUI {
		cnf.UI.Enabled = true
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

func loadConfValid(verbose bool, path string, defaultConf conf.Conf, defaultConfPath string) conf.Conf {
	if path == "" {
		path = defaultConfPath
		if verbose {
			fmt.Println("no config file specified, using default config file: ", path)
		}
	}
	// app executable dir + config.toml has the highest priority
	preferredPath := filepath.Join(getAppDir(), path)
	if _, err := os.Stat(preferredPath); err == nil {
		path = preferredPath
		if verbose {
			fmt.Println("using config file: ", path)
		}
	}

	err := decode(path, &defaultConf)
	if err != nil && verbose {
		os.Stderr.WriteString(fmt.Sprintf("Failed to load config file: %s\n", err))
	}else
	{
		logger.WithField("conf", &defaultConf).Debug("configuration loaded")
	}
	return defaultConf
}

func decode(path string, conf *conf.Conf) (error) {
	ext, err := misc.MustGetFileExt(path)
	if err != nil {
		return err
	}

	switch ext {
	case "toml":
		_, err = toml.DecodeFile(path, conf)
	case "yaml", "yml":
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read config file: %s", err)
		}
		
		err = yaml.Unmarshal([]byte(content), conf)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %s", err)
		}
	default:
		err = fmt.Errorf("unsupported config file extension: %s", ext)
	}

	return err
}
