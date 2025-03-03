package agent

import (
	"log"

	config "github.com/klimenkokayot/calc-net-go/internal/agent/config"
	"github.com/klimenkokayot/calc-net-go/internal/agent/service"
)

type Agent struct {
	Service          *service.AgentService
	OrchestratorPort int
}

func NewAgent() (*Agent, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	service := service.NewAgentService(*config)
	return &Agent{
		service,
		config.OrchestratorPort,
	}, nil
}

func (a *Agent) Run() error {
	log.Println("start new agent")
	return a.Service.Run()
}
