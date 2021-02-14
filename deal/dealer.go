package deal

import (
	"errors"
	"fmt"
	"time"

	"github.com/SSSaaS/sssa-golang"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/ptr"
	"github.com/giorgioazzinnaro/farmfa/random"
	"github.com/pquerna/otp/totp"
)

// ErrInvalidTOTPSecret is returned when testing the given secret against the TOTP algorithm unsuccessfully
var ErrInvalidTOTPSecret = errors.New("the provided value is not a valid TOTP secret key")

// ErrHighThreshold is returned when the number of players is lower than the threshold
var ErrHighThreshold = errors.New("threshold cannot be higher than the number of players")

// ErrLowThreshold is returned when the provided threshold makes no sense: if less than 2, there's no reason to use farMFA
var ErrLowThreshold = errors.New("threshold cannot be lower than 2")

// CreateTocs generates encrypted len(players) encrypted Tocs from secret; note is a message to remind the purpose of these Tocs
func CreateTocs(note, secret string, players []*Player, threshold int) (map[string]string, error) {
	groupSize := len(players)
	if err := validateThreshold(groupSize, threshold); err != nil {
		return nil, err
	}

	// This is just to validate the provided secret
	_, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return nil, ErrInvalidTOTPSecret
	}

	groupID, err := random.String(5)
	if err != nil {
		return nil, fmt.Errorf("failed to generate group ID: %w", err)
	}

	shares, err := sssa.Create(threshold, groupSize, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to split secret into shares: %w", err)
	}

	playerTocs := make(map[string]string, groupSize)

	for i, p := range players {

		tocID, err := random.String(5)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Toc ID: %w", err)
		}

		toc := &api.Toc{
			GroupId:        groupID,
			GroupSize:      groupSize,
			GroupThreshold: threshold,
			Note:           ptr.String(note),
			TocId:          tocID,
			Share:          shares[i],
		}

		enc, err := p.EncryptToc(toc)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt Toc for %q: %w", p.address, err)
		}

		playerTocs[p.address] = enc
	}

	return playerTocs, nil
}

func validateThreshold(players, threshold int) error {
	if threshold < 2 {
		return ErrLowThreshold
	}
	if players < threshold {
		return ErrHighThreshold
	}
	return nil
}
