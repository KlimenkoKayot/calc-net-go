package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-net-go/internal/shared/utils"
)

type Expressions struct {
	List []models.Expression `json:"expressions"`
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
		io.Writer(w).Write(utils.ErrorResponse(ErrInternalServer))
		return
	}
	defer r.Body.Close()

	expression := &models.Expression{}
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

	json, err := json.Marshal(models.Expression{
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
	expressions := h.Service.GetAllExpressions()
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
		json, _ := json.Marshal(models.Expression{
			Id:     id,
			Status: val.Status,
		})
		w.WriteHeader(http.StatusOK)
		io.Writer(w).Write(json)
		return
	}

	if val, found := h.Service.Answers[hash]; found {
		json, _ := json.Marshal(models.Expression{
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
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.Writer(w).Write(utils.ErrorResponse(err))
		return
	}
	data, _ := json.Marshal(task)
	w.WriteHeader(http.StatusOK)
	io.Writer(w).Write(data)
}

func (h *OrchestratorHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка при получении результата task: %e\n", err)
		return
	}
	taskAnswer := &models.TaskResult{}
	err = json.Unmarshal(data, taskAnswer)
	if err != nil {
		log.Printf("Ошибка при попытке парсинга json TaskAnswer: %e\n", err)
		return
	}
	log.Printf("Обработка ответа id: %d, ответ: %f", taskAnswer.Id, taskAnswer.Result)
	h.Service.ProcessAnswer(taskAnswer)
}
