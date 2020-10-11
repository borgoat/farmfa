package shares

import (
	"fmt"
	"github.com/SSSaaS/sssa-golang"
	"github.com/giorgioazzinnaro/farmfa/random"
)

func Split(secret TOTPSecret, threshold, total uint) ([]Token, error) {
	secretId, err := random.String(5)
	if err != nil {
		return nil, fmt.Errorf("could not generate a name for the secret: %w", err)
	}

	shares, err := sssa.Create(int(threshold), int(total), string(secret))
	if err != nil {
		return nil, fmt.Errorf("could not split shares: %w", err)
	}

	tokens := make([]Token, total)
	for i, share := range shares {
		tokens[i] = Token{
			Secret:    secretId,
			Serial:    uint(i),
			Threshold: threshold,
			Total:     total,
			Share:     share,
		}
	}

	return tokens, nil
}
