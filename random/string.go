package random

import "encoding/base32"

func String(n int) (string, error) {
	b, err := Bytes(n)
	return base32.StdEncoding.EncodeToString(b), err
}
