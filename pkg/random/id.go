package random

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

// ID Channel is for a new unique string,
// Used mainly at generating payment token
var ID = make(chan string)

func init() {
	// Make sure random generator is a bit fair random :)
	rand.Seed(int64(time.Now().Nanosecond()))

	go func() {
		h := sha1.New()
		c := []byte(time.Now().String())
		for {
			_, _ = h.Write(c)
			ID <- fmt.Sprintf("%x", h.Sum(nil))
		}
	}()
}
