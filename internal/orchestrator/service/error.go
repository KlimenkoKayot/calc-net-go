package orchestrator

import "fmt"

var (
	ErrInvalidSymbolRPN = fmt.Errorf("неизвестный символ в библиотеке rpn")
	ErrInvalidOperation = fmt.Errorf("неизвестная операция при попытке поиска времени выполнения")
	ErrHaveNoTask       = fmt.Errorf("нет задач")
	ErrEmptyRequestList = fmt.Errorf("список запросов пустой")
	ErrAnswerExpression = fmt.Errorf("попытка поиска задач в выражении, состоящем из ответа")
)
