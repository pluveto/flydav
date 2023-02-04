package app

import (
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/pluveto/flydav/cmd/flydav/conf"
	"github.com/pluveto/flydav/pkg/logger"
	"github.com/sirupsen/logrus"
)

func InitLogger(cnf conf.Log, verbose bool) {
	newLoggerCount := len(cnf.Stdout) + len(cnf.File)
	if newLoggerCount != 0 {
		for i := 0; i < newLoggerCount; i++ {
			logger.AddLogger(logrus.New())
		}
	} else {
		// if no logger configured, use default logger
		logger.SetOutput(os.Stdout)
		return
	}
	nextLoggerIndex := 1

	if verbose {
		logger.SetLevel(logrus.DebugLevel)
		// enable source code line numbers
		logger.SetReportCaller(true)
	} else {
		logger.SetLevel(levelToLogrusLevel(cnf.Level))
	}

	for _, stdout := range cnf.Stdout {
		currentLogger := logger.DefaultCombinedLogger.GetLogger(nextLoggerIndex)
		switch stdout.Format {
		case conf.LogFormatJSON:
			currentLogger.SetFormatter(&logrus.JSONFormatter{})
		case conf.LogFormatText:
			currentLogger.SetFormatter(&logrus.TextFormatter{})
		}
		switch stdout.Output {
		case conf.LogOutputStdout:
			currentLogger.SetOutput(os.Stdout)
		case conf.LogOutputStderr:
			currentLogger.SetOutput(os.Stderr)
		}
		nextLoggerIndex++
	}

	for _, file := range cnf.File {
		currentLogger := logger.DefaultCombinedLogger.GetLogger(nextLoggerIndex)

		switch file.Format {
		case conf.LogFormatJSON:
			currentLogger.SetFormatter(&logrus.JSONFormatter{})
		case conf.LogFormatText:
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

// levelToLogrusLevel converts a string to a logrus.Level
func levelToLogrusLevel(level string) logrus.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
