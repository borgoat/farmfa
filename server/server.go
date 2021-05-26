package server

import (
	"github.com/borgoat/farmfa/api"
	"github.com/borgoat/farmfa/session"
)

type Server struct {
	oracle *session.Oracle
}

var _ api.ServerInterface = (*Server)(nil)

func New(oracle *session.Oracle) *Server {
	return &Server{
		oracle: oracle,
	}
}
