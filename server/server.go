package server

import (
	"github.com/giorgioazzinnaro/farmfa/api"
	"github.com/giorgioazzinnaro/farmfa/session"
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
