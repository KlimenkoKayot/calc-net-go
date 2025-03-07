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
	Config *config.Config
}

// Создание нового экземпляра сервера
func NewServer() (*Server, error) {
	// Создаем конфиг
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	return &Server{
		config,
	}, nil
}

// Запуск сервера, использует роутер gorilla/mux
func (s *Server) Run() error {
	mux := mux.NewRouter()
	handler := handler.NewOrchestratorHandler(*s.Config)

	// Разные endpoint`ы
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/api/v1/calculate", handler.NewExpression)
	mux.HandleFunc("/api/v1/expressions", handler.Expressions)
	mux.HandleFunc("/api/v1/expressions/{id}", handler.Expression)
	mux.HandleFunc("/internal/task", handler.PostTask).Methods("POST")
	mux.HandleFunc("/internal/task", handler.GetTask).Methods("GET")

	staticDir := filepath.Join(".", "web", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Server started at port :%d\n", s.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), mux); err != nil {
		return err
	}
	return nil
}
