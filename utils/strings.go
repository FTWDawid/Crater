package utils

import "math/rand"

const (
	ALPHANUMERIC_CHARACTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ALPHANUMERIC_LENGTH     = len(ALPHANUMERIC_CHARACTERS)
)

func GenerateRandomCode() string {
	bytes := make([]byte, 6)
	for i := range bytes {
		bytes[i] = ALPHANUMERIC_CHARACTERS[rand.Intn(ALPHANUMERIC_LENGTH)]
	}
	return string(bytes)
}
