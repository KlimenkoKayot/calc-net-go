package server

import (
	"net/http"

	"github.com/klimenkokayot/calc-net-go/internal/config"
)

type Server struct {
	Config *config.Config
}

func NewServer() (*Server, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	return &Server{
		config,
	}, nil
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	// обработчики
}
