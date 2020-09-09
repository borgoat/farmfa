package shares

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	"time"
)

type Group []Token

func (g *Group) Add(token *Token) error {
	newTokens := append(*g, *token)

	if err := validateTokensFromSameSecret(newTokens); err != nil {
		return err
	}

	*g = removeDuplicateTokens(newTokens)

	return nil
}

func (g *Group) GenerateTOTP() (string, error) {
	secret, err := Combine(*g)
	if err != nil {
		return "", fmt.Errorf("shares could not be combined: %w", err)
	}

	code, err := totp.GenerateCode(string(secret), time.Now())
	if err != nil {
		return "", fmt.Errorf("could not use the secret to generate a TOTP code: %w", err)
	}

	return code, nil
}

// IsComplete checks if there are enough shares to reconstruct the secret
func (g *Group) IsComplete() bool {
	return uint(len(*g)) >= (*g)[0].Threshold
}
