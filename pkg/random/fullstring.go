package random

import "math/rand"

var chars = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

// String is the random string generator from all lowercase string
func String(l int) string {
	var res string
	ll := len(chars)
	for i := 0; i < l; i++ {
		res = res + string(chars[rand.Intn(ll)])
	}

	return res
}
