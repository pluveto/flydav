package app

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/pluveto/flydav/pkg/logger"
	"github.com/sirupsen/logrus"
)

func InitLogger(conf Log, verbose bool) {
	newLoggerCount := len(conf.Stdout) + len(conf.File)
	if newLoggerCount != 0 {
		for i := 0; i < newLoggerCount; i++ {
			logger.AddLogger(logrus.New())
		}
	}
	nextLoggerIndex := 1

	if verbose {
		logger.SetLevel(logrus.DebugLevel)
		// enable source code line numbers
		logger.SetReportCaller(true)
	} else {
		logger.SetLevel(levelToLogrusLevel(conf.Level))
	}

	for _, stdout := range conf.Stdout {
		currentLogger := logger.DefaultCombinedLogger.GetLogger(nextLoggerIndex)
		switch stdout.Format {
		case LogFormatJSON:
			currentLogger.SetFormatter(&logrus.JSONFormatter{})
		case LogFormatText:
			currentLogger.SetFormatter(&logrus.TextFormatter{})
		}
		switch stdout.Output {
		case LogOutputStdout:
			currentLogger.SetOutput(os.Stdout)
		case LogOutputStderr:
			currentLogger.SetOutput(os.Stderr)
		}
		nextLoggerIndex++
	}

	for _, file := range conf.File {
		currentLogger := logger.DefaultCombinedLogger.GetLogger(nextLoggerIndex)

		switch file.Format {
		case LogFormatJSON:
			currentLogger.SetFormatter(&logrus.JSONFormatter{})
		case LogFormatText:
			currentLogger.SetFormatter(&logrus.TextFormatter{})
		}
		currentLogger.SetOutput(&lumberjack.Logger{
			Filename:   file.Path,
			MaxSize:    file.MaxSize,
			MaxAge:     file.MaxAge,
			MaxBackups: 3,
			Compress:   true,
		})
		nextLoggerIndex++
	}

}
