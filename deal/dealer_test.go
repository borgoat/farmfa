package deal_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/borgoat/farmfa/deal"
	"github.com/borgoat/farmfa/random"
	"github.com/stretchr/testify/assert"
)

const ValidTOTPSecret = "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ"
const InvalidTOTPSecret = "Ceci n'est pas un TOTP secret"

const playersNum = 20

type playerIdentities struct {
	pl *deal.Player
	id age.Identity
}

var samplePlayerIdentities = func() map[string]*playerIdentities {
	players := make(map[string]*playerIdentities, playersNum)
	for i := 0; i < playersNum; i++ {
		r, _ := random.String(5)
		addr := fmt.Sprintf("%s@example.com", r)
		id, _ := age.GenerateX25519Identity()
		p, _ := deal.NewPlayer(addr, deal.EncryptWithAge(id.Recipient()))
		players[addr] = &playerIdentities{
			pl: p,
			id: id,
		}
	}
	return players
}()

var samplePlayers = func() []*deal.Player {
	players := make([]*deal.Player, len(samplePlayerIdentities))
	i := 0
	for _, identities := range samplePlayerIdentities {
		players[i] = identities.pl
		i++
	}
	return players
}()

func TestCreateTocs_basic(t *testing.T) {
	playerTocs, err := deal.CreateTocs("test_basic", ValidTOTPSecret, samplePlayers, playersNum/2)
	assert.NoError(t, err)
	assert.Len(t, playerTocs, playersNum)

	for i, toc := range playerTocs {
		plid := samplePlayerIdentities[i]

		r := strings.NewReader(toc)
		armorR := armor.NewReader(r)
		decrypter, err := age.Decrypt(armorR, plid.id)
		assert.NoError(t, err)

		w, err := io.Copy(ioutil.Discard, decrypter)
		assert.NoError(t, err)
		assert.Greater(t, w, int64(0), "encrypted toc is", toc)
	}
}

func TestCreateTocs_lowThreshold(t *testing.T) {
	_, err := deal.CreateTocs("test_low_threshold", ValidTOTPSecret, samplePlayers, 1)
	assert.ErrorIs(t, err, deal.ErrLowThreshold)
}

func TestCreateTocs_highThreshold(t *testing.T) {
	_, err := deal.CreateTocs("test_high_threshold", ValidTOTPSecret, samplePlayers, playersNum*2)
	assert.ErrorIs(t, err, deal.ErrHighThreshold)
}

func TestCreateTocs_invalidSecret(t *testing.T) {
	_, err := deal.CreateTocs("test_invalid_secret", InvalidTOTPSecret, samplePlayers, 3)
	assert.ErrorIs(t, err, deal.ErrInvalidTOTPSecret)
}
