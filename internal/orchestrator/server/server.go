package orchestrator

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
)

// Структура сервера
type Server struct {
	Config  *config.Config
	handler *handler.OrchestratorHandler
	mux     *mux.Router
}

// Создание нового экземпляра сервера
func NewServer(config *config.Config, handler *handler.OrchestratorHandler) (*Server, error) {
	mux := mux.NewRouter()
	return &Server{
		Config:  config,
		handler: handler,
		mux:     mux,
	}, nil
}

// Запуск сервера, использует роутер gorilla/mux
func (s *Server) Run() error {

	s.setupRoutes()

	log.Printf("Server started at port :%d\n", s.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), s.mux); err != nil {
		return err
	}
	return nil
}

func (s *Server) setupRoutes() error {
	// Разные endpoint`ы
	s.mux.HandleFunc("/", s.handler.Index)
	s.mux.HandleFunc("/api/v1/calculate", s.handler.NewExpression)
	s.mux.HandleFunc("/api/v1/expressions", s.handler.Expressions)
	s.mux.HandleFunc("/api/v1/expressions/{id}", s.handler.Expression)
	s.mux.HandleFunc("/internal/task", s.handler.PostTask).Methods("POST")
	s.mux.HandleFunc("/internal/task", s.handler.GetTask).Methods("GET")

	staticDir := filepath.Join(".", "web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	s.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	return nil
}
