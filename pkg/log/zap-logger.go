package log

import (
	"context"

	"go.uber.org/zap"

	"elbix.dev/engine/pkg/initializer"
)

var (
	logger *zap.Logger
)

type loggerInit struct {
}

func (loggerInit) Initialize(ctx context.Context) {
	go func() {
		<-ctx.Done()
		_ = logger.Sync()
	}()
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, f ...zap.Field) {
	logger.Debug(msg, f...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, f ...zap.Field) {
	logger.Info(msg, f...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, f ...zap.Field) {
	logger.Error(msg, f...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, f ...zap.Field) {
	logger.Fatal(msg, f...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, f ...zap.Field) {
	logger.Panic(msg, f...)
}

// Logger TODO: this is not good , return new logger for use in other places
func Logger() *zap.Logger {
	return logger
}

// SwapLogger used just for testing, do not use it in any other place
func SwapLogger(l *zap.Logger) {
	logger = l
}

func init() {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction(zap.AddCallerSkip(2))
		if err != nil {
			panic(err)
		}
	}

	initializer.Register(loggerInit{}, 0)
}
