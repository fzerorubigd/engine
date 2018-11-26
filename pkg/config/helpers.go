package config

import (
	"time"
)

// GetStruct fill an structure base on the config nested set, this function use reflection, and its not
// good (in my opinion) for frequent call.
// but its best if you need the config to loaded in structure and use that structure after that.
func GetStruct(k string, m interface{}) {
	o.GetStruct(k, m)
}

// GetIntDefault return an int value from Onion, if the value is not exists or its not an
// integer , default is returned
func GetIntDefault(key string, def int) int {
	return o.GetIntDefault(key, def)

}

// GetInt return an int value, if the value is not there, then it return zero value
func GetInt(key string) int {
	return o.GetInt(key)
}

// GetInt64Default return an int64 value from Onion, if the value is not exists or if the value is not
// int64 then return the default
func GetInt64Default(key string, def int64) int64 {
	return o.GetInt64Default(key, def)

}

// GetInt64 return the int64 value from config, if its not there, return zero
func GetInt64(key string) int64 {
	return o.GetInt64(key)
}

// GetFloat32Default return an float32 value from Onion, if the value is not exists or its not a
// float32, default is returned
func GetFloat32Default(key string, def float32) float32 {
	return o.GetFloat32Default(key, def)
}

// GetFloat32 return an float32 value, if the value is not there, then it returns zero value
func GetFloat32(key string) float32 {
	return o.GetFloat32(key)
}

// GetFloat64Default return an float64 value from Onion, if the value is not exists or if the value is not
// float64 then return the default
func GetFloat64Default(key string, def float64) float64 {
	return o.GetFloat64Default(key, def)
}

// GetFloat64 return the float64 value from config, if its not there, return zero
func GetFloat64(key string) float64 {
	return o.GetFloat64(key)
}

// GetStringDefault get a string from Onion. if the value is not exists or if tha value is not
// string, return the default
func GetStringDefault(key string, def string) string {
	return o.GetStringDefault(key, def)
}

// GetString is for getting an string from conig. if the key is not
func GetString(key string) string {
	return o.GetString(key)
}

// GetBoolDefault return bool value from Onion. if the value is not exists or if tha value is not
// boolean, return the default
func GetBoolDefault(key string, def bool) bool {
	return o.GetBoolDefault(key, def)

}

// GetBool is used to get a boolean value fro config, with false as default
func GetBool(key string) bool {
	return o.GetBool(key)
}

// GetDurationDefault is a function to get duration from config. it support both
// string duration (like 1h3m2s) and integer duration
func GetDurationDefault(key string, def time.Duration) time.Duration {
	return o.GetDurationDefault(key, def)
}

// GetDuration is for getting duration from config, it cast both int and string
// to duration
func GetDuration(key string) time.Duration {
	return o.GetDuration(key)
}

// GetStringSlice try to get a slice from the config
func GetStringSlice(key string) []string {
	return o.GetStringSlice(key)

}
