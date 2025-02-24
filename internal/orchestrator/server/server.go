package orchestrator

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
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
	handler := handler.NewOrchestratorHandler(*s.Config)
	mux.HandleFunc("/api/v1/calculate", handler.NewExpression)
	log.Printf("Server started at port :%d\n", s.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), mux); err != nil {
		return err
	}
	return nil
}
