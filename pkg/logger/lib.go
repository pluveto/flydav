package logger

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var DefaultCombinedLogger = New()

type nothingWriter struct{}

func (nothingWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func init() {
	DefaultCombinedLogger.AddHook(&logHook{})
	var nilWrite nothingWriter
	DefaultCombinedLogger.SetOutput(nilWrite)
}

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

type logHook struct {
}

func (hook *logHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *logHook) Fire(entry *logrus.Entry) error {
	for _, logger := range DefaultCombinedLogger.additionalLoggers {
		// pass entry to all additionalLoggers
		logger.WithFields(entry.Data).Log(entry.Level, entry.Message)
	}
	return nil
}

func AddLogger(logger *logrus.Logger) {
	DefaultCombinedLogger.additionalLoggers = append(DefaultCombinedLogger.additionalLoggers, logger)
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
	DefaultCombinedLogger.ApplyAll(
		func(l *logrus.Logger) {
			l.SetOutput(out)
		},
	)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter logrus.Formatter) {
	DefaultCombinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetFormatter(formatter)
	})
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	DefaultCombinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetReportCaller(include)
	})
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	DefaultCombinedLogger.ApplyAll(func(l *logrus.Logger) {
		l.SetLevel(level)
	})
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	return DefaultCombinedLogger.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level logrus.Level) bool {
	return DefaultCombinedLogger.IsLevelEnabled(level)
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	DefaultCombinedLogger.AddHook(hook)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *logrus.Entry {
	return DefaultCombinedLogger.WithField(logrus.ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func WithContext(ctx context.Context) *logrus.Entry {
	return DefaultCombinedLogger.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return DefaultCombinedLogger.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return DefaultCombinedLogger.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithTime(t time.Time) *logrus.Entry {
	return DefaultCombinedLogger.WithTime(t)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	DefaultCombinedLogger.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	DefaultCombinedLogger.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	DefaultCombinedLogger.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	DefaultCombinedLogger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	DefaultCombinedLogger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	DefaultCombinedLogger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	DefaultCombinedLogger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	DefaultCombinedLogger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	DefaultCombinedLogger.Fatal(args...)
}

// TraceFn logs a message from a func at level Trace on the standard logger.
func TraceFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.TraceFn(fn)
}

// DebugFn logs a message from a func at level Debug on the standard logger.
func DebugFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.DebugFn(fn)
}

// PrintFn logs a message from a func at level Info on the standard logger.
func PrintFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.PrintFn(fn)
}

// InfoFn logs a message from a func at level Info on the standard logger.
func InfoFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.InfoFn(fn)
}

// WarnFn logs a message from a func at level Warn on the standard logger.
func WarnFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.WarnFn(fn)
}

// WarningFn logs a message from a func at level Warn on the standard logger.
func WarningFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.WarningFn(fn)
}

// ErrorFn logs a message from a func at level Error on the standard logger.
func ErrorFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.ErrorFn(fn)
}

// PanicFn logs a message from a func at level Panic on the standard logger.
func PanicFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.PanicFn(fn)
}

// FatalFn logs a message from a func at level Fatal on the standard logger then the process will exit with status set to 1.
func FatalFn(fn logrus.LogFunction) {
	DefaultCombinedLogger.FatalFn(fn)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	DefaultCombinedLogger.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	DefaultCombinedLogger.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	DefaultCombinedLogger.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	DefaultCombinedLogger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	DefaultCombinedLogger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	DefaultCombinedLogger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	DefaultCombinedLogger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	DefaultCombinedLogger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	DefaultCombinedLogger.Fatalf(format, args...)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	DefaultCombinedLogger.Traceln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	DefaultCombinedLogger.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	DefaultCombinedLogger.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	DefaultCombinedLogger.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	DefaultCombinedLogger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	DefaultCombinedLogger.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	DefaultCombinedLogger.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	DefaultCombinedLogger.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	DefaultCombinedLogger.Fatalln(args...)
}
