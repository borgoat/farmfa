package server

import (
	"crypto/rand"
	"encoding/base32"
)

func generateRandomBytes(n int) ([]byte, error) {
	var b = make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(n int) (string, error) {
	b, err := generateRandomBytes(n)
	return base32.StdEncoding.EncodeToString(b), err
}
