package deal

import (
	"bytes"
	"encoding/json"
	"fmt"

	"filippo.io/age"
	"filippo.io/age/armor"

	"github.com/borgoat/farmfa/api"
)

// EncryptFunc turns a Toc into an encrypted string
type EncryptFunc func(toc *api.Toc) (string, error)

// Player represents an individual who receives an encrypted Toc
type Player struct {
	// address is a string to identify the recipient
	address string

	// enc is a function that can encrypt a Toc for the intended player
	enc EncryptFunc
}

func NewPlayer(address string, encryptionFunc EncryptFunc) (*Player, error) {
	return &Player{
		address: address,
		enc:     encryptionFunc,
	}, nil
}

// EncryptToc simply encrypts a Toc for the intended player
func (p *Player) EncryptToc(toc *api.Toc) (string, error) {
	return p.enc(toc)
}

// EncryptWithAge is used to encrypt a JSON-marshalled Toc with Age for a certain player
func EncryptWithAge(playerKeys ...age.Recipient) EncryptFunc {
	return func(toc *api.Toc) (string, error) {
		var out bytes.Buffer

		armorW := armor.NewWriter(&out)

		w, err := age.Encrypt(armorW, playerKeys...)
		if err != nil {
			return "", fmt.Errorf("failed to create encrypted buffer: %w", err)
		}

		jEnc := json.NewEncoder(w)
		err = jEnc.Encode(toc)
		if err != nil {
			return "", fmt.Errorf("failed to marshal Toc: %w", err)
		}

		err = w.Close()
		if err != nil {
			return "", fmt.Errorf("failed to close encrypted buffer %w", err)
		}

		err = armorW.Close()
		if err != nil {
			return "", fmt.Errorf("failed to close armor: %w", err)
		}
		return out.String(), nil
	}
}
