package orchestrator

import (
	"encoding/json"
	"io"
	"net/http"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
)

type Expression struct {
	Value string `json:"expression"`
}

type OrchestratorHandler struct {
	Service *service.OrchestratorService
}

func NewOrchestratorHandler(config config.Config) *OrchestratorHandler {
	return &OrchestratorHandler{
		Service: service.NewOrchestratorService(config),
	}
}

func (h *OrchestratorHandler) NewExpression(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write([]byte(ErrInternalServer.Error()))
		return
	}
	defer r.Body.Close()

	expression := &Expression{}
	err = json.Unmarshal(data, expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.Writer(w).Write([]byte(ErrInvalidBodyDecode.Error()))
		return
	}

	err = h.Service.AddExpression(expression.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write([]byte(ErrInternalServer.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
