package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"filippo.io/age/armor"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"filippo.io/age"
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/random"
)

type inMemSession struct {
	session       *api.Session
	encryptedTek  []byte
	encryptedTocs []string

	mu sync.RWMutex
}

type InMemory struct {
	sessions map[string]*inMemSession
}

func NewInMemory() Store {
	return &InMemory{sessions: map[string]*inMemSession{}}
}

func (i *InMemory) CreateSession(tocZero *api.Toc) (*api.SessionCredentials, error) {
	var (
		resp api.SessionCredentials
		sess inMemSession
	)

	sessID, err := random.String(10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a session ID: %w", err)
	}
	sess.session = &api.Session{
		Id:            sessID,
		CreatedAt:     time.Now(),
		Status:        "pending",
		Complete:      false,
		TocGroupId:    tocZero.GroupId,
		TocsInGroup:   tocZero.GroupSize,
		TocsThreshold: tocZero.GroupThreshold,
		TocsProvided:  1,
	}
	resp.Session = *sess.session

	// A TEK is generated and used to encrypt Toc zero
	tek, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tek: %w", err)
	}
	// The session public key is returned as Tek
	resp.Tek = tek.Recipient().String()

	// the TEK private key is kept in encrypted form (encrypted by KEK)
	encTek, kek, err := encryptedTek(tek)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kek: %w", err)
	}
	sess.encryptedTek = encTek

	// KEK is then returned to the applicant, which will provide it again when generating the TOTP
	resp.Kek = kek

	// Toc zero needs to be encrypted like the others and then stored in the current session
	encTocZero, err := encryptToc(tek.Recipient(), tocZero)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt toc zero: %w", err)
	}
	sess.encryptedTocs = make([]string, 1, tocZero.GroupThreshold)
	sess.encryptedTocs[0] = encTocZero

	i.sessions[sessID] = &sess
	return &resp, nil
}

// encryptedTek is used to encrypt the Toc encryption key (private key), with a key encryption key (one-time pad)
func encryptedTek(identity *age.X25519Identity) (encryptedTek, kek []byte, err error) {
	tek := []byte(identity.String())
	kek, err = random.Bytes(len(tek))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate one time pad: %w", err)
	}

	encryptedTek = make([]byte, len(tek))
	for i := range kek {
		encryptedTek[i] = tek[i] ^ kek[i]
	}

	return encryptedTek, kek, nil
}

// encryptToc is used just for tocZero that needs to be stored encrypted in memory
func encryptToc(id age.Recipient, toc *api.Toc) (string, error) {
	var out bytes.Buffer

	armOut := armor.NewWriter(&out)
	ageOut, err := age.Encrypt(armOut, id)
	if err != nil {
		return "", fmt.Errorf("failed to create encrypted buffer: %w", err)
	}

	jEnc := json.NewEncoder(ageOut)
	err = jEnc.Encode(toc)
	if err != nil {
		return "", fmt.Errorf("failed to encode as JSON: %w", err)
	}

	err = ageOut.Close()
	if err != nil {
		return "", fmt.Errorf("failed to flush buffer: %w", err)
	}

	err = armOut.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close armored buffer: %w", err)
	}

	return out.String(), nil
}

var ErrSessionNotFound = errors.New("session not found")

func (i *InMemory) GetSession(id string) (*api.Session, error) {
	s, found := i.sessions[id]
	if !found {
		return nil, fmt.Errorf("id %s invalid: %w", id, ErrSessionNotFound)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.session, nil
}

var ErrEmptyToc = errors.New("the provided Toc is empty")
var ErrTocIsNotEncrypted = errors.New("the provided Toc is not a valid age armored string")
var ErrTocAlreadyExists = errors.New("the provided Toc already exists")

func (i *InMemory) AddToc(id string, encryptedToc string) error {
	s, found := i.sessions[id]
	if !found {
		return fmt.Errorf("id %s invalid: %w", id, ErrSessionNotFound)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if encryptedToc == "" {
		return ErrEmptyToc
	}

	if !isValidAgeArmoredString(encryptedToc) {
		return ErrTocIsNotEncrypted
	}

	for _, et := range s.encryptedTocs {
		if et == encryptedToc {
			return ErrTocAlreadyExists
		}
	}

	s.encryptedTocs = append(s.encryptedTocs, encryptedToc)

	return nil
}

func isValidAgeArmoredString(armored string) bool {
	r := strings.NewReader(armored)
	ar := armor.NewReader(r)
	w, err := io.Copy(ioutil.Discard, ar)

	if err != nil {
		return false
	}

	if w == 0 {
		return false
	}

	return true
}

func (i *InMemory) DecryptTocs(id string, key *api.SessionKeyEncryptionKey) ([]api.Toc, error) {
	s, found := i.sessions[id]
	if !found {
		return nil, fmt.Errorf("id %s invalid: %w", id, ErrSessionNotFound)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO Ensure there's enough Tocs

	tek, err := decryptTek(key, s.encryptedTek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt TEK: %w", err)
	}

	decryptedTocs := make([]api.Toc, len(s.encryptedTocs))
	for i, encToc := range s.encryptedTocs {
		decToc, err := decryptToc(tek, encToc)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt a Toc: %w", err)
		}
		decryptedTocs[i] = *decToc
	}

	return decryptedTocs, nil
}

var ErrKekInvalidLength = errors.New("the provided kek has an invalid length")

func decryptTek(kek *api.SessionKeyEncryptionKey, encryptedTek []byte) (age.Identity, error) {
	tekLen := len(encryptedTek)
	if len(kek.Kek) != tekLen {
		return nil, ErrKekInvalidLength
	}

	decryptedTek := make([]byte, tekLen)
	for i := 0; i < tekLen; i++ {
		decryptedTek[i] = kek.Kek[i] ^ encryptedTek[i]
	}

	id, err := age.ParseX25519Identity(string(decryptedTek))
	if err != nil {
		return nil, fmt.Errorf("the parsed TEK is not a valid age identity: %w", err)
	}

	return id, nil
}

func decryptToc(identity age.Identity, encryptedToc string) (*api.Toc, error) {
	sr := strings.NewReader(encryptedToc)
	ar := armor.NewReader(sr)
	r, err := age.Decrypt(ar, identity)
	if err != nil {
		return nil, fmt.Errorf("failed to create decryptor: %w", err)
	}

	jDec := json.NewDecoder(r)

	var toc api.Toc
	err = jDec.Decode(&toc)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return &toc, nil
}
