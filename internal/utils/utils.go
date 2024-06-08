package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandId(n int) (string, error) {
	bs := make([]byte, n)
	if _, err := rand.Read(bs); err != nil {
		return "", err
	}
	return hex.EncodeToString(bs), nil
}
