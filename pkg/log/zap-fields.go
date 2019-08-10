package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Err is the error field
func Err(err error) zapcore.Field {
	if err == nil {
		return zap.Skip()
	}
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

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) zapcore.Field {
	return zap.Time(key, val)
}
