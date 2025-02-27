package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
	"github.com/klimenkokayot/calc-net-go/internal/shared/utils"
)

type Expression struct {
	Id     string  `json:"id,omitempty"`
	Status string  `json:"status,omitempty"`
	Result float64 `json:"result,omitempty"`

	Value string `json:"expression,omitempty"`
}

type Expressions struct {
	List []Expression `json:"expressions"`
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
	fmt.Println("ABOBOA")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}
	defer r.Body.Close()

	expression := &Expression{}
	err = json.Unmarshal(data, expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.Writer(w).Write(utils.ErrorResponse(ErrInvalidBodyDecode))
		return
	}

	hash, err := h.Service.AddExpression(expression.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}

	json, err := json.Marshal(Expression{
		Id: utils.EncodeToString(hash),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.Writer(w).Write(json)
}

func (h *OrchestratorHandler) Expressions(w http.ResponseWriter, r *http.Request) {
	expressions := []Expression{}
	for _, val := range h.Service.Expressions {
		expressions = append(expressions, Expression{
			Id:     utils.EncodeToString(val.Hash),
			Status: val.Status,
		})
	}
	for hash, val := range h.Service.Answers {
		expressions = append(expressions, Expression{
			Id:     utils.EncodeToString(hash),
			Status: "Выполнено.",
			Result: val,
		})
	}

	json, err := json.Marshal(expressions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.Writer(w).Write(json)
}

func (h *OrchestratorHandler) Expression(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	hash, err := utils.EncodedToSHA512(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	}

	if val, found := h.Service.Expressions[hash]; found {
		json, _ := json.Marshal(Expression{
			Id:     id,
			Status: val.Status,
		})
		w.WriteHeader(http.StatusOK)
		io.Writer(w).Write(json)
		return
	}

	if val, found := h.Service.Answers[hash]; found {
		json, _ := json.Marshal(Expression{
			Id:     id,
			Status: "Выполнено.",
			Result: val,
		})
		w.WriteHeader(http.StatusOK)
		io.Writer(w).Write(json)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *OrchestratorHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.Service.GetTask()
	if err == service.ErrHaveNoTask {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(task)
	w.WriteHeader(http.StatusOK)
	io.Writer(w).Write(data)
}

func (h *OrchestratorHandler) PostTask(w http.ResponseWriter, r *http.Request) {

}
