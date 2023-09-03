package util

import (
	"math/rand"
)

func IsBlank(str string) bool {
	return len(str) == 0
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	_, _ = rand.Read(randomBytes)
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(randomBytes)
}

func ToPtr[T any](t T) *T {
	return &t
}
