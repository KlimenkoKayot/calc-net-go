package models

import (
	"time"

	"github.com/klimenkokayot/calc-net-go/internal/shared/customList"
)

type Expression struct {
	Hash [64]byte
	List *customList.LinkedList
}

type Task struct {
	Id             uint          `json:"id"`
	FirstArgument  float64       `json:"arg1"`
	SecondArgument float64       `json:"arg2"`
	Operation      rune          `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
	StartListNode  *customList.Node
	ExpressionId   [64]byte
}

type Result struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
