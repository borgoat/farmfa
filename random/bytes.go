package random

import "crypto/rand"

func Bytes(n int) ([]byte, error) {
	var b = make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
