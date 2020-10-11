package server

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/sessions"
)

type Server struct {
	sm sessions.SessionManager
}

var _ api.ServerInterface = (*Server)(nil)

func New(sm sessions.SessionManager) *Server {
	return &Server{
		sm: sm,
	}
}
