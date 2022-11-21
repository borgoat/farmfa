package session_test

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/borgoat/farmfa/api"
	"github.com/borgoat/farmfa/session"
	"github.com/stretchr/testify/assert"
)

func genericOracleCreateSesssion(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "TTXAIGO4",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          []byte{0x33, 0x5D, 0x58, 0x8C, 0x29, 0x7A, 0x81, 0x43, 0x5D, 0x9B, 0x8E, 0x50, 0xC3, 0xEF, 0xB0, 0x7D, 0x62, 0xE5, 0xA7, 0xD8, 0xA3, 0x12, 0x85, 0x59, 0x8C, 0xE4, 0x7, 0xA9, 0x4F, 0x20, 0xD3, 0xB8, 0x12},
		TocId:          []byte{0x26, 0x7, 0x2B, 0x5, 0xF3, 0x7D, 0x10, 0x38, 0x44},
	})
	assert.NoError(t, err)
	assert.Equal(t, false, creds.Complete)

	// assert that TEK is a valid age recipient
	_, err = age.ParseX25519Recipient(creds.Tek)
	assert.NoError(t, err)
}

func genericOracleAddToc_valid(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "ZUX44STM",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          []byte{0x90, 0x8D, 0xA3, 0xC6, 0xB0, 0x9A, 0xEE, 0x7E, 0x85, 0xA3, 0xB3, 0xE, 0x19, 0xEE, 0x25, 0x28, 0xD9, 0x37, 0x25, 0xA7, 0xB2, 0xBE, 0x65, 0xF6, 0xAD, 0x5B, 0xED, 0xD5, 0x9B, 0x69, 0xF9, 0xC6, 0xA},
		TocId:          []byte{0xB4, 0xF, 0x68, 0x1B, 0x5D, 0xB8, 0x42, 0x87, 0x6D},
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBaSXR1aVIyTFFoUjFTMS9q\nRDBPTm05V3JxV3Rvc2tqdGFaYXdEQ3ZZeVhzCjRLRE96OVFkNFE3WUVpMUxmdCth\nTlpJbVFLaCtteDhMZ0pVVU5wanI2L3MKLS0tIFdiSXlrSGFSdTdDb2pRenJCN3dW\ndXVNN3hhTWJLdm05Vy9WRlF6cHdDdzgKlaMq9Aiaz98nHMzmGoRoyjU+DxxR3tVh\n/B981Y+P2ZHIaS9Dzr4pR2SbX/Y/SdTdYkACvET3uXNvVAtFMgETQcQwan2sEkty\nVmb82mem0r9/WyWQ6N1psebsEkicOo2LGJkcDgeNlSWmRPrT830NRim4LLqVkWpe\npqaUdRcIu0H4eP4yAdHwXAyUtuualkLZuax/t7D4W5mzu+UOdqYkOqvh84wpX9Ml\nQ2CzgL+KACOmE/lMh5qGxfydnxw0z19p\n-----END AGE ENCRYPTED FILE-----\n")
	assert.NoError(t, err)
}

func genericOracleAddToc_empty(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "ZUX44STM",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          []byte{0x90, 0x8D, 0xA3, 0xC6, 0xB0, 0x9A, 0xEE, 0x7E, 0x85, 0xA3, 0xB3, 0xE, 0x19, 0xEE, 0x25, 0x28, 0xD9, 0x37, 0x25, 0xA7, 0xB2, 0xBE, 0x65, 0xF6, 0xAD, 0x5B, 0xED, 0xD5, 0x9B, 0x69, 0xF9, 0xC6, 0xA},
		TocId:          []byte{0xB4, 0xF, 0x68, 0x1B, 0x5D, 0xB8, 0x42, 0x87, 0x6D},
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "")
	assert.ErrorIs(t, err, session.ErrEmptyToc)
}

func genericOracleAddToc_notEncrypted(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "ZUX44STM",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          []byte{0x90, 0x8D, 0xA3, 0xC6, 0xB0, 0x9A, 0xEE, 0x7E, 0x85, 0xA3, 0xB3, 0xE, 0x19, 0xEE, 0x25, 0x28, 0xD9, 0x37, 0x25, 0xA7, 0xB2, 0xBE, 0x65, 0xF6, 0xAD, 0x5B, 0xED, 0xD5, 0x9B, 0x69, 0xF9, 0xC6, 0xA},
		TocId:          []byte{0xB4, 0xF, 0x68, 0x1B, 0x5D, 0xB8, 0x42, 0x87, 0x6D},
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "this is not an age armored string")
	assert.ErrorIs(t, err, session.ErrTocIsNotEncrypted)
}

func genericOracleAddToc_alreadyExists(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "ZUX44STM",
		GroupSize:      20,
		GroupThreshold: 10,
		Note:           nil,
		Share:          []byte{0x90, 0x8D, 0xA3, 0xC6, 0xB0, 0x9A, 0xEE, 0x7E, 0x85, 0xA3, 0xB3, 0xE, 0x19, 0xEE, 0x25, 0x28, 0xD9, 0x37, 0x25, 0xA7, 0xB2, 0xBE, 0x65, 0xF6, 0xAD, 0x5B, 0xED, 0xD5, 0x9B, 0x69, 0xF9, 0xC6, 0xA},
		TocId:          []byte{0xB4, 0xF, 0x68, 0x1B, 0x5D, 0xB8, 0x42, 0x87, 0x6D},
	})
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBaSXR1aVIyTFFoUjFTMS9q\nRDBPTm05V3JxV3Rvc2tqdGFaYXdEQ3ZZeVhzCjRLRE96OVFkNFE3WUVpMUxmdCth\nTlpJbVFLaCtteDhMZ0pVVU5wanI2L3MKLS0tIFdiSXlrSGFSdTdDb2pRenJCN3dW\ndXVNN3hhTWJLdm05Vy9WRlF6cHdDdzgKlaMq9Aiaz98nHMzmGoRoyjU+DxxR3tVh\n/B981Y+P2ZHIaS9Dzr4pR2SbX/Y/SdTdYkACvET3uXNvVAtFMgETQcQwan2sEkty\nVmb82mem0r9/WyWQ6N1psebsEkicOo2LGJkcDgeNlSWmRPrT830NRim4LLqVkWpe\npqaUdRcIu0H4eP4yAdHwXAyUtuualkLZuax/t7D4W5mzu+UOdqYkOqvh84wpX9Ml\nQ2CzgL+KACOmE/lMh5qGxfydnxw0z19p\n-----END AGE ENCRYPTED FILE-----\n")
	assert.NoError(t, err)

	err = m.AddToc(creds.Id, "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBaSXR1aVIyTFFoUjFTMS9q\nRDBPTm05V3JxV3Rvc2tqdGFaYXdEQ3ZZeVhzCjRLRE96OVFkNFE3WUVpMUxmdCth\nTlpJbVFLaCtteDhMZ0pVVU5wanI2L3MKLS0tIFdiSXlrSGFSdTdDb2pRenJCN3dW\ndXVNN3hhTWJLdm05Vy9WRlF6cHdDdzgKlaMq9Aiaz98nHMzmGoRoyjU+DxxR3tVh\n/B981Y+P2ZHIaS9Dzr4pR2SbX/Y/SdTdYkACvET3uXNvVAtFMgETQcQwan2sEkty\nVmb82mem0r9/WyWQ6N1psebsEkicOo2LGJkcDgeNlSWmRPrT830NRim4LLqVkWpe\npqaUdRcIu0H4eP4yAdHwXAyUtuualkLZuax/t7D4W5mzu+UOdqYkOqvh84wpX9Ml\nQ2CzgL+KACOmE/lMh5qGxfydnxw0z19p\n-----END AGE ENCRYPTED FILE-----\n")
	assert.ErrorIs(t, err, session.ErrTocAlreadyExists)
}

func genericOracleGenerateTOTP(t *testing.T, store session.Store) {
	m := session.NewOracle(store)
	creds, err := m.CreateSession(&api.Toc{
		GroupId:        "3OROTSEU",
		GroupSize:      5,
		GroupThreshold: 2,
		Note:           nil,
		Share:          []byte{0x71, 0x26, 0x3, 0x97, 0xFD, 0x6F, 0x8C, 0x44, 0x4A, 0x61, 0x40, 0x16, 0xEE, 0x3E, 0x24, 0x66, 0x9F, 0x12, 0xE, 0x98, 0x67, 0x94, 0x6D, 0xDF, 0x38, 0x4C, 0x3, 0x6, 0x7C, 0xE8, 0xC9, 0xBF, 0xA0},
		TocId:          []byte{0xB0, 0x9F, 0x1C, 0x9B, 0x3F, 0x4B, 0xA2, 0xC3, 0x9B},
	})
	assert.NoError(t, err)

	tek, err := age.ParseX25519Recipient(creds.Tek)
	assert.NoError(t, err)

	var out bytes.Buffer
	aw := armor.NewWriter(&out)
	w, err := age.Encrypt(aw, tek)
	assert.NoError(t, err)
	jEnc := json.NewEncoder(w)

	err = jEnc.Encode(&api.Toc{
		GroupId:        "3OROTSEU",
		GroupSize:      5,
		GroupThreshold: 2,
		Note:           nil,
		Share:          []byte{0x95, 0xC5, 0x4, 0xD2, 0x67, 0xFB, 0x85, 0x58, 0x4A, 0xCD, 0x1C, 0x19, 0x28, 0xD5, 0x78, 0xDE, 0x92, 0x65, 0x45, 0x9D, 0x30, 0xCD, 0xC1, 0xC1, 0xDF, 0x64, 0x44, 0x31, 0xEC, 0x76, 0xE7, 0x6, 0x16},
		TocId:          []byte{0xB0, 0x9F, 0x1C, 0x9B, 0x3F, 0x4B, 0xA2, 0xC3, 0x9B},
	})
	assert.NoError(t, err)

	assert.NoError(t, w.Close())
	assert.NoError(t, aw.Close())

	err = m.AddToc(creds.Id, out.String())
	assert.NoError(t, err)

	key := &api.SessionKeyEncryptionKey{Kek: creds.Kek}

	totp, err := m.GenerateTOTP(creds.Id, key)
	assert.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile(`^\d{6}$`), totp)
}
