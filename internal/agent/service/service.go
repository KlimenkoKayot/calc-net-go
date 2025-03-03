package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	config "github.com/klimenkokayot/calc-net-go/internal/agent/config"
	worker "github.com/klimenkokayot/calc-net-go/internal/agent/worker"
)

type AgentService struct {
	OrchestratorUrl string
	AgentSleepTime  time.Duration
	ComputingPower  uint64
	wg              *sync.WaitGroup
}

func NewAgentService(config config.Config) *AgentService {
	return &AgentService{
		fmt.Sprintf("http://127.0.0.1:%d/internal/task", config.OrchestratorPort),
		config.AgentSleepTime,
		config.ComputingPower,
		&sync.WaitGroup{},
	}
}

func (s *AgentService) StartNewWorker() error {
	log.Println("Запущен новый worker...")
	s.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		worker := worker.NewWorker(s.OrchestratorUrl, s.AgentSleepTime)
		worker.Run()
		wg.Done()
	}(s.wg)
	return nil
}

func (s *AgentService) Run() error {
	for i := 0; i < int(s.ComputingPower); i++ {
		s.StartNewWorker()
	}
	s.wg.Wait()
	return nil
}
