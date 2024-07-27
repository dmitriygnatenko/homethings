package helpers

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] // nolint:gosec
	}
	return string(b)
}
