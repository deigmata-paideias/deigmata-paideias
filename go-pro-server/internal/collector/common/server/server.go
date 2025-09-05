package server

import "collector-go/internal/collector/common/config"

type Server struct {
	// collector server config
	config *config.Config
	// logger
	// name .....
}

func New(cfg *config.Config) *Server {

	return &Server{config: cfg}
}

func (s *Server) Validate() error {

	return nil
}
