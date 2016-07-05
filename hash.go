package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// randInt returns a random int in the specified range
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// randomString returns a random string of the specified length
func randomString(l int) string {
	chars := make([]byte, l)
	for i := 0; i < l; i++ {
		r := rand.Intn(14)
		switch {
		case r < 5:
			chars[i] = byte(randInt(65, 90))
		case r > 4 && r < 10:
			chars[i] = byte(randInt(97, 122))
		case r > 9:
			chars[i] = byte(randInt(48, 57))
		}
	}
	return string(chars)
}

// NewHash returns a unique sha1 hash of a randomly generated 32char string
func NewHash() string {
	h := sha1.New()
	u := fmt.Sprintf("%s", randomString(32))

	h.Write([]byte(u))
	return hex.EncodeToString(h.Sum(nil))
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
