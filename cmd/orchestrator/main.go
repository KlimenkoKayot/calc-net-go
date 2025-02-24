package main

import (
	"log"

	orchestrator "github.com/klimenkokayot/calc-net-go/internal/orchestrator/server"
)

func main() {
	server, err := orchestrator.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
