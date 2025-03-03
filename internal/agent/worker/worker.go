package agent

import (
	"log"
	"net/http"
	"time"

	transport "github.com/klimenkokayot/calc-net-go/internal/agent/worker/transport"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

type Worker struct {
	Client    *http.Client
	URL       string
	SleepTime time.Duration
}

func NewWorker(url string, sleepTime time.Duration) *Worker {
	return &Worker{
		&http.Client{},
		url,
		sleepTime,
	}
}

func (w *Worker) Solve(task *models.Task) *models.TaskResult {
	log.Printf("Worker получил новую подзадачу c id: %d\n", task.Id)
	time.Sleep(time.Millisecond * task.OperationTime)
	var answer float64
	switch task.Operation {
	case '+':
		answer = task.FirstArgument + task.SecondArgument
	case '-':
		answer = task.FirstArgument - task.SecondArgument
	case '*':
		answer = task.FirstArgument * task.SecondArgument
	case '/':
		answer = task.FirstArgument / task.SecondArgument
	}
	log.Printf("Получен ответ на подзадачу с id: %d, ответ: %f\n", task.Id, answer)
	return &models.TaskResult{
		Id:     task.Id,
		Result: answer,
	}
}

func (w *Worker) Process() error {
	for {
		task, err := transport.GetTask(w.Client, w.URL)
		if err != nil {
			return err
		}
		result := w.Solve(task)
		err = transport.PostTask(w.Client, w.URL, result)
		if err != nil {
			return err
		}
	}
}

func (w *Worker) Run() error {
	for {
		err := w.Process()
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(w.SleepTime)
	}
}
