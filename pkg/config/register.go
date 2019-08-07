package config

import (
	"time"

	"github.com/goraz/onion/configwatch"
)

// RegisterString add an string to config
func RegisterString(key string, def string, description string) configwatch.String {
	setDescription(key, description)
	return configwatch.RegisterString(key, def)
}

// RegisterInt add an int to config
func RegisterInt(key string, def int, description string) configwatch.Int {
	setDescription(key, description)
	return configwatch.RegisterInt(key, def)
}

// RegisterInt64 add an int to config
func RegisterInt64(key string, def int64, description string) configwatch.Int {
	setDescription(key, description)
	return configwatch.RegisterInt64(key, def)
}

// RegisterFloat64 add an int to config
func RegisterFloat64(key string, def float64, description string) configwatch.Float {
	setDescription(key, description)
	return configwatch.RegisterFloat64(key, def)
}

// RegisterBoolean add a boolean to config
func RegisterBoolean(key string, def bool, description string) configwatch.Bool {
	setDescription(key, description)
	return configwatch.RegisterBool(key, def)
}

// RegisterDuration add an duration to config
func RegisterDuration(key string, def time.Duration, description string) configwatch.Int {
	setDescription(key, description)
	return configwatch.RegisterDuration(key, def)
}
