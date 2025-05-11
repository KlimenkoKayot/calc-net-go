package main

import (
	"log"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	orchestrator "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server"
	handler "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server/handler"
	service "github.com/klimenkokayot/calc-net-go/internal/orchestrator/service"
)

func main() {
	// Создаем конфиг
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Создаем новый сервер aka оркестратора
	orchestratorService := service.NewOrchestratorService(config)
	orchestratorHandler := handler.NewOrchestratorHandler(config, orchestratorService)
	server, err := orchestrator.NewServer(config, orchestratorHandler)
	if err != nil {
		log.Fatal(err)
	}
	// Запуск созданного сервера
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
