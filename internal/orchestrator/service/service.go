package orchestrator

import (
	"fmt"
	"log"
	"sync"
	"time"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	"github.com/klimenkokayot/calc-net-go/internal/shared/customList"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-net-go/internal/shared/utils"
	"github.com/klimenkokayot/calc-net-go/pkg/rpn"
)

type OrchestratorService struct {
	TimeAdditionMs        time.Duration
	TimeSubtractionMs     time.Duration
	TimeMultiplicationsMs time.Duration
	TimeDivisionsMs       time.Duration

	TaskIdUpdate       map[uint]*customList.Node
	Tasks              []*models.Task
	Answers            map[[64]byte]float64
	Expressions        map[[64]byte]*models.Expression
	RequestExpressions [][64]byte

	LastTaskId uint
	mu         *sync.Mutex
}

func NewOrchestratorService(config config.Config) *OrchestratorService {
	return &OrchestratorService{
		time.Duration(config.TimeAdditionMs),
		time.Duration(config.TimeSubtractionMs),
		time.Duration(config.TimeMultiplicationsMs),
		time.Duration(config.TimeDivisionsMs),

		make(map[uint]*customList.Node, 0),
		make([]*models.Task, 0),
		make(map[[64]byte]float64),
		make(map[[64]byte]*models.Expression),
		make([][64]byte, 0),

		0,
		&sync.Mutex{},
	}
}

func (s *OrchestratorService) NewExpression(expression string) (*models.Expression, error) {
	valuesIntergace, err := rpn.ExpressionToRPN(expression)
	if err != nil {
		return nil, err
	}
	list := customList.NewLinkedList()
	for _, val := range valuesIntergace {
		switch val.(type) {
		case string:
			list.Add(&customList.NodeData{
				IsOperation: true,
				Operation:   []rune(val.(string))[0],
			})
		case float64:
			list.Add(&customList.NodeData{
				Value: val.(float64),
			})
		default:
			return nil, ErrInvalidSymbolRPN
		}
	}
	hash := utils.ExpressionToSHA512(expression)
	return &models.Expression{
		Id:     utils.EncodeToString(hash),
		Value:  expression,
		Hash:   hash,
		List:   list,
		Status: "В очереди.",
	}, nil
}

func (s *OrchestratorService) AddExpression(expression string) ([64]byte, error) {
	log.Printf("Получена новая задача: %s\n", expression)
	value, err := s.NewExpression(expression)
	if err != nil {
		log.Printf("error: %v\n", err)
		return [64]byte{}, err
	}
	s.mu.Lock()
	_, ansFound := s.Answers[value.Hash]
	_, expFound := s.Expressions[value.Hash]
	s.mu.Unlock()
	if !ansFound && !expFound {
		s.mu.Lock()
		s.Expressions[value.Hash] = value
		s.RequestExpressions = append(s.RequestExpressions, value.Hash)
		s.mu.Unlock()
	}
	return value.Hash, nil
}

func (s *OrchestratorService) GetAllExpressions() []models.Expression {
	expressions := make([]models.Expression, 0)
	fmt.Println("ABOBA1")
	s.mu.Lock()
	fmt.Println("ABOBA3")
	for _, val := range s.Expressions {
		expressions = append(expressions, models.Expression{
			Id:     utils.EncodeToString(val.Hash),
			Status: val.Status,
		})
	}
	for hash, val := range s.Answers {
		expressions = append(expressions, models.Expression{
			Id:     utils.EncodeToString(hash),
			Status: "Выполнено.",
			Result: val,
		})
	}
	s.mu.Unlock()
	fmt.Println("ABOBA2")
	return expressions
}

func (s *OrchestratorService) OperationTime(operation rune) (time.Duration, error) {
	switch operation {
	case '+':
		return s.TimeAdditionMs, nil
	case '-':
		return s.TimeDivisionsMs, nil
	case '*':
		return s.TimeMultiplicationsMs, nil
	case '/':
		return s.TimeSubtractionMs, nil
	default:
		return 0, ErrInvalidOperation
	}
}

func (s *OrchestratorService) FindNewTasks(expression *models.Expression) ([]*models.Task, error) {
	tasks := []*models.Task{}
	if expression.List.Root.Next == nil {
		log.Printf("Получени ответ на задачу: %s, ответ: %f\n", expression.Value, expression.List.Root.Data.Value)
		// добавляем ответ на значение
		s.mu.Lock()
		s.Answers[expression.Hash] = expression.List.Root.Data.Value
		// удаляем, т.к. посчитали
		delete(s.Expressions, expression.Hash)
		s.mu.Unlock()
		return nil, ErrAnswerExpression
	}
	cur := expression.List.Root
	var (
		last     *customList.Node
		previous *customList.Node
	)
	for cur != nil {
		if cur.Data.IsOperation {
			fmt.Printf("%s ", string(rune(cur.Data.Operation)))
		} else {
			fmt.Printf("%.0f ", cur.Data.Value)
		}
		if (last != nil && !last.InAction &&
			previous != nil && !previous.InAction) && last.Data.IsOperation && !previous.Data.IsOperation && !cur.Data.IsOperation {
			operationTime, err := s.OperationTime(last.Data.Operation)
			if err != nil {
				return nil, err
			}
			s.mu.Lock()
			tasks = append(tasks, &models.Task{
				Id:             s.LastTaskId,
				FirstArgument:  cur.Data.Value,
				SecondArgument: previous.Data.Value,
				Operation:      last.Data.Operation,
				OperationTime:  operationTime,
			})
			s.TaskIdUpdate[s.LastTaskId] = last
			last.InAction = true
			previous.InAction = true
			cur.InAction = true
			s.LastTaskId++
			s.mu.Unlock()
			last = nil
			previous = nil
			cur = cur.Next
		} else {
			last = previous
			previous = cur
			cur = cur.Next
		}
	}
	fmt.Println()
	return tasks, nil
}

func (s *OrchestratorService) GetTask() (*models.Task, error) {
	task := &models.Task{}
	reqExpIdx := 0
	for len(s.Tasks) == 0 {
		if len(s.Expressions) == 0 || reqExpIdx == len(s.Expressions) {
			return nil, ErrHaveNoTask
		} else if len(s.RequestExpressions) == 0 {
			return nil, ErrEmptyRequestList
		}
		s.mu.Lock()
		hash := s.RequestExpressions[reqExpIdx]
		s.mu.Unlock()
		tasks, err := s.FindNewTasks(s.Expressions[hash])
		if err == ErrAnswerExpression {
			s.RequestExpressions = append(s.RequestExpressions[:reqExpIdx], s.RequestExpressions[reqExpIdx+1:]...)
			continue
		} else if err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.Tasks = append(s.Tasks, tasks...)
		s.mu.Unlock()
		reqExpIdx++
	}
	s.mu.Lock()
	task = s.Tasks[0]
	s.Tasks = s.Tasks[1:]
	s.mu.Unlock()
	return task, nil
}

func (s *OrchestratorService) ProcessAnswer(taskAnswer *models.TaskResult) {
	s.mu.Lock()
	node := s.TaskIdUpdate[taskAnswer.Id]
	// удаление ненужного ключа
	delete(s.TaskIdUpdate, taskAnswer.Id)
	node.Data = &customList.NodeData{
		Value: taskAnswer.Result,
	}
	node.Next = node.Next.Next.Next
	node.InAction = false
	s.mu.Unlock()
}
