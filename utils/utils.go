package utils

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandString generate random string
func GenerateRandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		randNum := rand.Intn(len(letterBytes))
		b[i] = letterBytes[randNum]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
