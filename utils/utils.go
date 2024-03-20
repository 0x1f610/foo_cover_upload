package utils

import (
	"math/rand"
	"time"
)

const letters = "abcdefghjkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ23456789"

var source = rand.NewSource(time.Now().UnixNano())
var rng = rand.New(source)

func GenerateString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}
