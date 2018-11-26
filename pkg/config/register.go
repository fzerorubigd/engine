package config

import (
	"sync"
	"time"

	"gopkg.in/fzerorubigd/onion.v3"
)

var (
	lock = &sync.Mutex{}
)

// RegisterString add an string to config
func RegisterString(key string, def string, description string) onion.String {
	setDescription(key, description)
	return o.RegisterString(key, def)
}

// RegisterInt add an int to config
func RegisterInt(key string, def int, description string) onion.Int {
	setDescription(key, description)
	return o.RegisterInt(key, def)
}

// RegisterInt64 add an int to config
func RegisterInt64(key string, def int64, description string) onion.Int {
	setDescription(key, description)
	return o.RegisterInt64(key, def)
}

// RegisterFloat64 add an int to config
func RegisterFloat64(key string, def float64, description string) onion.Float {
	setDescription(key, description)
	return o.RegisterFloat64(key, def)
}

// RegisterBoolean add a boolean to config
func RegisterBoolean(key string, def bool, description string) onion.Bool {
	setDescription(key, description)
	return o.RegisterBool(key, def)
}

// RegisterDuration add an duration to config
func RegisterDuration(key string, def time.Duration, description string) onion.Int {
	setDescription(key, description)
	return o.RegisterDuration(key, def)
}
