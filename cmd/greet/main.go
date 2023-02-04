package main

import (
	"os"
	"path/filepath"

	"example.com/m/cmd/greet/app"
	"example.com/m/pkg/logger"
	"github.com/BurntSushi/toml"
	"github.com/alexflint/go-arg"
)

// main Entry point of the application
func main() {
	// main Entry point of the application
	var args app.Args
	var conf app.Conf
	var defaultConf = app.GetDefaultConf()

	args = loadArgsValid()
	conf = loadConfValid(args.Config, defaultConf, "config.toml")

	app.InitLogger(defaultConf.Log, args.Verbose)
	logger.Debug("log level: ", logger.GetLevel())

	app.Run(args, conf)
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

func loadConfValid(path string, defaultConf app.Conf, defaultConfPath string) app.Conf {
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
