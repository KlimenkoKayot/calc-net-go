package models

import (
	"time"

	"github.com/klimenkokayot/calc-net-go/internal/shared/customList"
)

type Expression struct {
	Id     string  `json:"id,omitempty"`
	Status string  `json:"status,omitempty"`
	Result float64 `json:"result,omitempty"`
	Value  string  `json:"expression,omitempty"`

	Hash [64]byte               `json:"-"`
	List *customList.LinkedList `json:"-"`
}

type Task struct {
	Id             uint          `json:"id"`
	FirstArgument  float64       `json:"arg1"`
	SecondArgument float64       `json:"arg2"`
	Operation      rune          `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
}

type TaskResult struct {
	Id     uint    `json:"id"`
	Result float64 `json:"result"`
}

type Result struct {
	Id     uint    `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
