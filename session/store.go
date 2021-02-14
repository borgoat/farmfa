package session

import "github.com/giorgioazzinnaro/farmfa/api"

// Store is used to keep track of Sessions, and their Tocs
type Store interface {
	CreateSession(tocZero *api.Toc) (*api.SessionCredentials, error)
	GetSession(id string) (*api.Session, error)
	AddToc(id string, encryptedToc string) error
	DecryptTocs(id string, key *api.SessionKeyEncryptionKey) ([]api.Toc, error)
}
