package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Err is the error field
func Err(err error) zapcore.Field {
	return zap.String("error", err.Error())
}

// String is a wrapper for zap string field
func String(key, val string) zapcore.Field {
	return zap.String(key, val)
}

// Any return zapcore field based on type of val
func Any(key string, val interface{}) zapcore.Field {
	return zap.Any(key, val)
}
