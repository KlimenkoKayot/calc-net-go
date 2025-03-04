package main

import (
	"log"

	orchestrator "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server"
)

func main() {
	// Создаем новый сервер aka оркестратора
	server, err := orchestrator.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	// Запуск созданного сервера
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
