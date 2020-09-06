package sessions

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/random"
	"github.com/giorgioazzinnaro/farmfa/shares"
	"sync"
	"time"
)

type inMemSession struct {
	session *api.PrivateSession
	tokens  shares.Group

	mu sync.Mutex
}

type InMemory struct {
	sessions map[SessionIdentifier]*inMemSession
}

func NewInMemory() SessionManager {
	return &InMemory{
		sessions: map[SessionIdentifier]*inMemSession{},
	}
}

func (i *InMemory) Start(firstShare *shares.Token) (*api.PrivateSession, error) {
	var resp api.PrivateSession
	var err error

	*resp.Id, err = random.String(25)
	if err != nil {
		return nil, err
	}
	*resp.Private, err = random.String(25)
	if err != nil {
		return nil, err
	}

	*resp.CreatedAt = time.Now()
	*resp.ShareGroup = firstShare.Secret
	*resp.Shares = int(firstShare.Total)
	*resp.Threshold = int(firstShare.Threshold)
	*resp.Complete = false
	*resp.Closed = false

	i.sessions[SessionIdentifier(*resp.Id)] = &inMemSession{
		session: &resp,
		tokens:  []shares.Token{
			*firstShare,
		},
	}

	return &resp, nil
}

func (i *InMemory) AddShare(id SessionIdentifier, share *shares.Token) error {
	// TODO Handle if i.sessions[id] is not there
	session := i.sessions[id]

	session.mu.Lock()
	defer session.mu.Unlock()

	err := session.tokens.Add(share)
	if err != nil {
		return err
	}

	return nil
}

func (i *InMemory) Status(id SessionIdentifier) (*api.Session, error) {
	// TODO Handle if i.sessions[id] is not there
	return &i.sessions[id].session.Session, nil
}
