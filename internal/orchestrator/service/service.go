package orchestrator

import (
	"time"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
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
