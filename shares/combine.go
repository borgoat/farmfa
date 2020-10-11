package shares

import (
	"errors"
	"fmt"
	"github.com/SSSaaS/sssa-golang"
)

type TOTPSecret string

func Combine(tokens []Token) (TOTPSecret, error) {
	noDuplicates := removeDuplicateTokens(tokens)

	if err := validateGroupOfTokens(noDuplicates); err != nil {
		return "", fmt.Errorf("invalid tokens: %w", err)
	}

	shares := getSharesFromTokens(noDuplicates)

	secret, err := sssa.Combine(shares)
	if err != nil {
		return "", fmt.Errorf("failed to combine shares: %w", err)
	}

	return TOTPSecret(secret), nil
}

func removeDuplicateTokens(tokens []Token) []Token {
	encountered := make(map[uint]struct{}, len(tokens))
	j := 0
	for _, token := range tokens {
		if _, exists := encountered[token.Serial]; exists {
			continue
		}
		encountered[token.Serial] = struct{}{}
		tokens[j] = token
		j++
	}
	return tokens[:j]
}

func validateGroupOfTokens(tokens []Token) error {
	var err error
	if err = validateTokensFromSameSecret(tokens); err != nil {
		return err
	}

	if err = validateEnoughTokens(tokens); err != nil {
		return err
	}

	return nil
}

var ErrNotSameSecret = errors.New("the provided tokens are not from the same secret")

func validateTokensFromSameSecret(tokens []Token) error {
	var secretId string
	for _, token := range tokens {
		if secretId != "" && secretId != token.Secret {
			return ErrNotSameSecret
		}
		secretId = token.Secret
	}
	return nil
}

var ErrNotEnoughTokens = errors.New("the provided tokens are not enough to cross the threshold")

func validateEnoughTokens(tokens []Token) error {
	if tokens[0].Threshold > uint(len(tokens)) {
		return ErrNotEnoughTokens
	}
	return nil
}

func getSharesFromTokens(tokens []Token) []string {
	shares := make([]string, len(tokens))
	for i, token := range tokens {
		shares[i] = token.Share
	}
	return shares
}
