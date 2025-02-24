package orchestrator

import "fmt"

var (
	ErrInvalidSymbolRPN = fmt.Errorf("неизвестный символ в библиотеке rpn")
	ErrInvalidOperation = fmt.Errorf("неизвестная операция при попытке поиска времени выполнения")
)
