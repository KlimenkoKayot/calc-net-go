package orchestrator

import (
	"time"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-net-go/pkg/rpn"
)

type OrchestratorService struct {
	TimeAdditionMs        time.Duration
	TimeSubtractionMs     time.Duration
	TimeMultiplicationsMs time.Duration
	TimeDivisionsMs       time.Duration
	Expressions           []models.Expression
	Tasks                 []models.Task
}

func NewOrchestratorService(config config.Config) *OrchestratorService {
	return &OrchestratorService{
		time.Duration(config.TimeAdditionMs),
		time.Duration(config.TimeSubtractionMs),
		time.Duration(config.TimeMultiplicationsMs),
		time.Duration(config.TimeDivisionsMs),
		make([]models.Expression, 0),
		make([]models.Task, 0),
	}
}

func (s *OrchestratorService) ConvertExpression(expression string) (*models.Expression, error) {
	valuesIntergace, err := rpn.ExpressionToRPN(expression)
	if err != nil {
		return nil, err
	}
	result := models.Expression{}
	for _, val := range valuesIntergace {
		switch val.(type) {
		case string:
			result = append(result, models.Symbol{IsOperation: true, Operation: []rune(val.(string))[0]})
		case float64:
			result = append(result, models.Symbol{Value: val.(float64)})
		default:
			return nil, ErrInvalidSymbolRPN
		}
	}
	return &result, nil
}
