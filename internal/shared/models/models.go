package models

import "time"

type Expression []Symbol

type Symbol struct {
	Value       float64
	Operation   rune
	IsOperation bool
}

type Task struct {
	Id             int           `json:"id"`
	FirstArgument  float64       `json:"arg1"`
	SecondArgument float64       `json:"arg2"`
	Operation      rune          `json:"operation"`
	OperationTime  time.Duration `json:"operation_time"`
}

type Result struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
