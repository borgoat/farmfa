package session_test

import (
	"github.com/borgoat/farmfa/session"
	"testing"
)

func TestMemory_CreateSession(t *testing.T) {
	genericOracleCreateSesssion(t, session.NewInMemoryStore())
}

func TestMemory_AddToc_valid(t *testing.T) {
	genericOracleAddToc_valid(t, session.NewInMemoryStore())
}

func TestMemory_AddToc_empty(t *testing.T) {
	genericOracleAddToc_empty(t, session.NewInMemoryStore())
}

func TestMemory_AddToc_notEncrypted(t *testing.T) {
	genericOracleAddToc_notEncrypted(t, session.NewInMemoryStore())
}

func TestMemory_AddToc_alreadyExists(t *testing.T) {
	genericOracleAddToc_alreadyExists(t, session.NewInMemoryStore())
}

func TestMemory_GenerateTOTP(t *testing.T) {
	genericOracleGenerateTOTP(t, session.NewInMemoryStore())
}
