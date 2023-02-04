package app

import (
	"os"

	"example.com/m/pkg/logger"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitLogger(conf Log, verbose bool) {

	if verbose {
		logger.SetLevel(logrus.DebugLevel)
		// enable source code line numbers
		logger.SetReportCaller(true)
	} else {
		logger.SetLevel(levelToLogrusLevel(conf.Level))
	}
	for _, stdout := range conf.Stdout {
		switch stdout.Format {
		case LogFormatJSON:
			logger.SetFormatter(&logrus.JSONFormatter{})
		case LogFormatText:
			logger.SetFormatter(&logrus.TextFormatter{})
		}
		switch stdout.Output {
		case LogOutputStdout:
			logger.SetOutput(os.Stdout)
		case LogOutputStderr:
			logger.SetOutput(os.Stderr)
		}
	}

	logger.Debugf("%d stdout logger loaded", len(conf.Stdout))

	for _, file := range conf.File {
		switch file.Format {
		case LogFormatJSON:
			logger.SetFormatter(&logrus.JSONFormatter{})
		case LogFormatText:
			logger.SetFormatter(&logrus.TextFormatter{})
		}
		logger.SetOutput(&lumberjack.Logger{
			Filename:   file.Path,
			MaxSize:    file.MaxSize,
			MaxAge:     file.MaxAge,
			MaxBackups: 3,
			Compress:   true,
		})
	}

	logger.Debugf("%d file logger loaded", len(conf.File))
}
