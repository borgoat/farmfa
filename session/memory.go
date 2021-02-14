package session

import (
	"sync"

	"github.com/giorgioazzinnaro/farmfa/api"
)

type inMemSession struct {
	session       *api.Session
	encryptedTek  []byte
	encryptedTocs []string

	mu sync.Mutex
}

// InMemoryStore simply stores session in the process memory.
// If HA or reliability in the event of failure is needed, this is clearly not a great option...
type InMemoryStore struct {
	sessions map[string]*inMemSession
}

func NewInMemoryStore() Store {
	return &InMemoryStore{sessions: map[string]*inMemSession{}}
}

func (i *InMemoryStore) CreateSession(session *api.Session, encryptedTEK []byte, encryptedTocZero string) error {

	if _, found := i.sessions[session.Id]; found {
		return ErrSessionAlreadyExists
	}

	var s inMemSession
	s.session = session
	s.encryptedTek = encryptedTEK
	s.encryptedTocs = make([]string, 1)
	s.encryptedTocs[0] = encryptedTocZero

	i.sessions[session.Id] = &s

	return nil
}

func (i *InMemoryStore) GetSession(id string) (*api.Session, error) {
	s, found := i.sessions[id]
	if !found {
		return nil, ErrSessionNotFound
	}

	return s.session, nil
}

func (i *InMemoryStore) AddEncryptedToc(id string, encryptedToc string) error {
	s, found := i.sessions[id]
	if !found {
		return ErrSessionNotFound
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, et := range s.encryptedTocs {
		if et == encryptedToc {
			return ErrTocAlreadyExists
		}
	}

	s.encryptedTocs = append(s.encryptedTocs, encryptedToc)

	return nil
}

func (i *InMemoryStore) GetEncryptedTocs(id string) ([]string, error) {
	s, found := i.sessions[id]
	if !found {
		return nil, ErrSessionNotFound
	}

	return s.encryptedTocs, nil
}

func (i *InMemoryStore) GetTEK(id string) ([]byte, error) {
	s, found := i.sessions[id]
	if !found {
		return nil, ErrSessionNotFound
	}

	return s.encryptedTek, nil
}
