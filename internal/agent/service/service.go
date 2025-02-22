package service

import (
	config "github.com/klimenkokayot/calc-net-go/internal/agent/config"
)

type AgentService struct {
	ComputingPower uint64
}

func NewAgentService(config config.Config) *AgentService {
	return &AgentService{
		config.ComputingPower,
	}
}
