package handler

import (
	"math/rand"
	"os"
	"time"
)

// GetRandomString - Generate random string
func GetRandomString(length int) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randomBytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[r.Intn(len(charset))]
	}

	return string(randomBytes)
}

func DeleteUploadFile(filepath string) {
	if _, err := os.Stat(filepath); err == nil {
		os.Remove(filepath)
	}
}
