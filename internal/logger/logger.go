package logger

import (
	"fmt"
	"os"

	logrus_filename "github.com/exgalibas/logrus-filename"
	"github.com/pluveto/flydav/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	fileLogger    *logrus.Logger
	consoleLogger *logrus.Logger
}

// NewLogger creates a new Logger instance.
func NewLogger(level logrus.Level, enableFileLogger bool, filePath string, wrapLevel int) (*Logger, error) {
	logger := &Logger{}

	// Set up the console logger with colors and caller reporting
	logger.consoleLogger = logrus.New()
	logger.consoleLogger.AddHook(logrus_filename.NewHook(logrus_filename.WithSkip(2 + wrapLevel)))
	logger.consoleLogger.SetLevel(level)
	logger.consoleLogger.Formatter = &CustomTextFormatter{}

	// If logging to a file is enabled, set up the file logger
	if enableFileLogger {
		fileLogger := logrus.New()
		fileLogger.AddHook(logrus_filename.NewHook(logrus_filename.WithSkip(2 + wrapLevel)))
		fileLogger.SetLevel(level)
		fileLogger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileLogger.Out = file
		logger.fileLogger = fileLogger
	}

	return logger, nil
}

type CustomTextFormatter struct{}

var colors = map[logrus.Level]string{
	logrus.DebugLevel: "\033[37m",
	logrus.InfoLevel:  "\033[32m",
	logrus.WarnLevel:  "\033[33m",
	logrus.ErrorLevel: "\033[31m",
	logrus.FatalLevel: "\033[31m",
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := entry.Level
	color, ok := colors[level]
	if !ok {
		color = "\033[37m"
	}

	levelText := level.String()
	if len(levelText) < 5 {
		levelText += " "
	}

	file := entry.Data["file"]
	s := fmt.Sprintf("%s[%s] \033[30m %s \033[0m%s\n", color, levelText, file, entry.Message)
	return []byte(s), nil
}

func (l *Logger) log(level logrus.Level, args ...interface{}) {
	if l.fileLogger != nil {
		l.fileLogger.Log(level, args...)
	}
	if l.consoleLogger != nil {
		l.consoleLogger.Log(level, args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(logrus.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(logrus.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(logrus.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(logrus.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(logrus.FatalLevel, args...)
	// Consider proper resource cleanup before exiting
	// or return an error and let the caller handle it.
}

var DefaultLogger *Logger

func Init(cfg config.LogConfig) {
	var err error
	DefaultLogger, err = NewLogger(logrus.InfoLevel, cfg.Enabled, cfg.Path, 1)
	if err != nil {
		panic(err)
	}
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}
