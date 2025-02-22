package orchestrator

import (
	"fmt"
	"net/http"

	config "github.com/klimenkokayot/calc-net-go/internal/agent/config"
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
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), mux); err != nil {
		return err
	}
	return nil
}
