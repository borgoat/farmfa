package server

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/session"
)

type Server struct {
	sm session.Store
}

var _ api.ServerInterface = (*Server)(nil)

func New(sm session.Store) *Server {
	return &Server{
		sm: sm,
	}
}
