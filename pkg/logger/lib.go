package logger

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var combinedLogger = New()

type CombinedLogger struct {
	*logrus.Logger
	additionalLoggers []*logrus.Logger
}

func New() *CombinedLogger {
	return &CombinedLogger{
		Logger:            logrus.New(),
		additionalLoggers: []*logrus.Logger{},
	}
}

func AddLogger(logger *logrus.Logger) {
	combinedLogger.additionalLoggers = append(combinedLogger.additionalLoggers, logger)
}

func (logger *CombinedLogger) GetLogger(index int) *logrus.Logger {
	if index == 0 {
		return logger.Logger
	}
	if index > len(logger.additionalLoggers) || index < 0 {
		return nil
	}
	return logger.additionalLoggers[index-1]
}

func (logger *CombinedLogger) Apply(index int, fn func(*logrus.Logger)) {
	if index == 0 {
		fn(logger.Logger)
		return
	}
	if index > len(logger.additionalLoggers) || index < 0 {
		panic("index of logger out of range")
	}
	fn(logger.additionalLoggers[index-1])
}

func (logger *CombinedLogger) ApplyAll(fn func(*logrus.Logger)) {
	fn(logger.Logger)
	for _, l := range logger.additionalLoggers {
		fn(l)
	}
}

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	combinedLogger.ApplyAll(
		func(l *logrus.Logger) {
			l.SetOutput(out)
		},
	)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	combinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetFormatter(formatter)
	})
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	combinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetReportCaller(include)
	})
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	combinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetLevel(level)
	})
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	return combinedLogger.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level logrus.Level) bool {
	return combinedLogger.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	combinedLogger.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return combinedLogger.WithField(logrus.ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return combinedLogger.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return combinedLogger.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return combinedLogger.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *logrus.Entry {
	return combinedLogger.WithTime(t)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	combinedLogger.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	combinedLogger.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	combinedLogger.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	combinedLogger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	combinedLogger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	combinedLogger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	combinedLogger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	combinedLogger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	combinedLogger.Fatal(args...)
}

// TraceFn logs a message from a func at level Trace on the standard logger.
func TraceFn(fn logrus.LogFunction) {
	combinedLogger.TraceFn(fn)
}

// DebugFn logs a message from a func at level Debug on the standard logger.
func DebugFn(fn logrus.LogFunction) {
	combinedLogger.DebugFn(fn)
}

// PrintFn logs a message from a func at level Info on the standard logger.
func PrintFn(fn logrus.LogFunction) {
	combinedLogger.PrintFn(fn)
}

// InfoFn logs a message from a func at level Info on the standard logger.
func InfoFn(fn logrus.LogFunction) {
	combinedLogger.InfoFn(fn)
}

// WarnFn logs a message from a func at level Warn on the standard logger.
func WarnFn(fn logrus.LogFunction) {
	combinedLogger.WarnFn(fn)
}

// WarningFn logs a message from a func at level Warn on the standard logger.
func WarningFn(fn logrus.LogFunction) {
	combinedLogger.WarningFn(fn)
}

// ErrorFn logs a message from a func at level Error on the standard logger.
func ErrorFn(fn logrus.LogFunction) {
	combinedLogger.ErrorFn(fn)
}

// PanicFn logs a message from a func at level Panic on the standard logger.
func PanicFn(fn logrus.LogFunction) {
	combinedLogger.PanicFn(fn)
}

// FatalFn logs a message from a func at level Fatal on the standard logger then the process will exit with status set to 1.
func FatalFn(fn logrus.LogFunction) {
	combinedLogger.FatalFn(fn)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	combinedLogger.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	combinedLogger.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	combinedLogger.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	combinedLogger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	combinedLogger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	combinedLogger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	combinedLogger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	combinedLogger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	combinedLogger.Fatalf(format, args...)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	combinedLogger.Traceln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	combinedLogger.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	combinedLogger.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	combinedLogger.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	combinedLogger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	combinedLogger.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	combinedLogger.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	combinedLogger.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	combinedLogger.Fatalln(args...)
}
