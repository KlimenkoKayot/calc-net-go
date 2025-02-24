package orchestrator

import (
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

	Tasks       []*models.Task
	Answers     map[[64]byte]float64
	Expressions map[[64]byte]*models.Expression

	LastTaskId uint
}

func NewOrchestratorService(config config.Config) *OrchestratorService {
	return &OrchestratorService{
		time.Duration(config.TimeAdditionMs),
		time.Duration(config.TimeSubtractionMs),
		time.Duration(config.TimeMultiplicationsMs),
		time.Duration(config.TimeDivisionsMs),
		make([]*models.Task, 0),
		make(map[[64]byte]float64),
		make(map[[64]byte]*models.Expression),
		0,
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
			list.Add(customList.NodeData{
				IsOperation: true,
				Operation:   []rune(val.(string))[0],
			})
		case float64:
			list.Add(customList.NodeData{
				Value: val.(float64),
			})
		default:
			return nil, ErrInvalidSymbolRPN
		}
	}
	return &models.Expression{
		Hash: utils.ExpressionToSHA512(expression),
		List: list,
	}, nil
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

func (s *OrchestratorService) FindNewTasks(expression models.Expression) (*[]models.Task, error) {
	tasks := []models.Task{}
	cur := expression.List.Root
	var (
		last     *customList.Node
		previous *customList.Node
	)
	for cur != nil {
		if last != nil && previous != nil && last.Data.IsOperation && !previous.Data.IsOperation && !cur.Data.IsOperation {
			operationTime, err := s.OperationTime(last.Data.Operation)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, models.Task{
				Id:             s.LastTaskId,
				FirstArgument:  cur.Data.Value,
				SecondArgument: previous.Data.Value,
				Operation:      last.Data.Operation,
				OperationTime:  operationTime,
				StartListNode:  last,
				ExpressionId:   expression.Hash,
			})
			s.LastTaskId++
			last = nil
			previous = nil
		} else {
			last = previous
			previous = cur
		}
	}
	return &tasks, nil
}
