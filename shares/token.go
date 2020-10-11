package shares

import (
	"encoding/base64"
	"fmt"
	"github.com/fxamacker/cbor/v2"
)

type Token struct {
	// The secret this token belongs to
	Secret string
	// The serial of this token inside the group
	Serial uint

	// The number of tokens needed to retrieve the secret
	Threshold uint
	// The total secrets that have been created
	Total uint

	// The share needed for SSSS
	Share string
}

func (t *Token) String() (string, error) {
	cborBytes, err := cbor.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("error while marshalling to CBOR: %w", err)
	}

	return base64.StdEncoding.EncodeToString(cborBytes), nil
}

func Parse(rawToken string) (*Token, error) {
	cborBytes, err := base64.StdEncoding.DecodeString(rawToken)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 string: %w", err)
	}

	var t Token
	err = cbor.Unmarshal(cborBytes, &t)
	if err != nil {
		return nil, fmt.Errorf("provided data is not valid CBOR: %w", err)
	}
	return &t, nil
}
