package sessions

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/shares"
)

type SessionIdentifier string

type SessionManager interface {
	Start(firstShare *shares.Token) (*api.PrivateSession, error)
	AddShare(id SessionIdentifier, share *shares.Token) error
	Status(id SessionIdentifier) (*api.Session, error)
}
