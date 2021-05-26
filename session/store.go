package session

import (
	"errors"
	"github.com/borgoat/farmfa/api"
)

// Store is used to keep track of Sessions, and their Tocs
// ErrSessionNotFound may be returned by any method where a session is referenced by ID
type Store interface {
	// CreateSession is used to persist a new session.
	// It may return ErrSessionAlreadyExists if the ID in session has been already used
	CreateSession(session *api.Session, encryptedTEK []byte, encryptedTocZero string) error

	// GetSession retrieves a session by its ID.
	// ErrSessionNotFound may be returned if the provided ID does not exist.
	GetSession(id string) (*api.Session, error)

	// AddEncryptedToc is used to append a Toc to an existing session.
	// It may return ErrTocAlreadyExists if the provided value has already been seen
	// ErrSessionNotFound may also be returned.
	AddEncryptedToc(id string, encryptedToc string) error

	// GetEncryptedTocs returns a slice of strings with the encrypted Tocs
	GetEncryptedTocs(id string) ([]string, error)

	// GetTEK is used to retrieve the Toc encryption key
	GetTEK(id string) ([]byte, error)

	// GarbageCollect should process all sessions by sending them to the shouldDelete func;
	// if it returns true, remove the given session from storage
	GarbageCollect(shouldDelete func(session *api.Session) bool)
}

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionAlreadyExists = errors.New("a session with the requested ID already exists")
var ErrTocAlreadyExists = errors.New("the provided Toc already exists")
