package orchestrator

import (
	"log"
	"sync"
	"time"

	config "github.com/klimenkokayot/calc-net-go/internal/orchestrator/config"
	"github.com/klimenkokayot/calc-net-go/internal/shared/customList"
	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
	"github.com/klimenkokayot/calc-net-go/internal/shared/utils"
	"github.com/klimenkokayot/calc-net-go/pkg/rpn"
)

// Структура сервиса, для удобства некоторые переменные из конфига перенесены сюда
type OrchestratorService struct {
	TimeAdditionMs        time.Duration // Время решение задачи сложения
	TimeSubtractionMs     time.Duration // Время решения задачи деления
	TimeMultiplicationsMs time.Duration // Время решения задачи умножения
	TimeDivisionsMs       time.Duration // Время решения задачи вычитания

	TaskIdUpdate       map[uint]*customList.Node       // Словарь, который по id подзадачи может найти указатель на элемент RPN forward list
	TaskIdExpression   map[uint]*models.Expression     // Словарь, который по id подзадачи может найти указатель на выражение
	Tasks              []*models.Task                  // Список подзадач для решения
	Answers            map[[64]byte]float64            // Ответы на разные задачи, которые были обработаны ранее
	Expressions        map[[64]byte]*models.Expression // Словарь выражений по хэшу
	RequestExpressions [][64]byte                      // Список всех полученных запросов на подсчет в порядке времени получения запроса

	LastTaskId uint // Счетчик для индексации решения подзадач
	mu         *sync.Mutex
}

// Создает новый экземпляр сервиса оркестартора
func NewOrchestratorService(config config.Config) *OrchestratorService {
	return &OrchestratorService{
		time.Duration(config.TimeAdditionMs),
		time.Duration(config.TimeSubtractionMs),
		time.Duration(config.TimeMultiplicationsMs),
		time.Duration(config.TimeDivisionsMs),

		make(map[uint]*customList.Node, 0),
		make(map[uint]*models.Expression, 0),
		make([]*models.Task, 0),
		make(map[[64]byte]float64),
		make(map[[64]byte]*models.Expression),
		make([][64]byte, 0),

		0,
		&sync.Mutex{},
	}
}

// Создание экземпляра выражения
func (s *OrchestratorService) NewExpression(expression string) (*models.Expression, error) {
	valuesIntergace, err := rpn.ExpressionToRPN(expression)
	if err != nil {
		return nil, err
	}
	// Тут мы создаем новую структуру Linked List,
	// адаптированную под нашу задачу EXPRESSION -> RPN
	// каждый элемент RPN будет отдельным элементом в Linked List
	list := customList.NewLinkedList()
	for _, val := range valuesIntergace {
		// Добавляем новые элементы в Linked List
		// наш стек RPN зареверсится
		switch val.(type) {
		case string:
			list.Add(&customList.NodeData{
				IsOperation: true,
				Operation:   []rune(val.(string))[0],
			})
		case float64:
			list.Add(&customList.NodeData{
				Value: val.(float64),
			})
		default:
			return nil, ErrInvalidSymbolRPN
		}
	}
	// Подсчитываем hash, для добавления нового выражения в сервис
	hash := utils.ExpressionToSHA512(expression)
	return &models.Expression{
		Id:     utils.EncodeToString(hash),
		Value:  expression,
		Hash:   hash,
		List:   list,
		Status: "В очереди.",
	}, nil
}

// Обработка получения новой задачи в сервис
func (s *OrchestratorService) AddExpression(expression string) ([64]byte, error) {
	log.Printf("Получена новая задача: %s\n", expression)
	value, err := s.NewExpression(expression)
	if err != nil {
		log.Printf("error: %v\n", err)
		return [64]byte{}, err
	}
	// Проверка на наличие задачи с сервисе
	// (нужно для того, чтобы не считать подсчитанные запросы)
	s.mu.Lock()
	_, ansFound := s.Answers[value.Hash]
	_, expFound := s.Expressions[value.Hash]
	s.mu.Unlock()
	if !ansFound && !expFound {
		s.mu.Lock()
		s.Expressions[value.Hash] = value
		s.RequestExpressions = append(s.RequestExpressions, value.Hash)
		s.mu.Unlock()
	}
	return value.Hash, nil
}

// Формирование списка выражений (в обработке/выполнено)
func (s *OrchestratorService) GetAllExpressions() []models.Expression {
	expressions := make([]models.Expression, 0)
	s.mu.Lock()
	for _, val := range s.Expressions {
		expressions = append(expressions, models.Expression{
			Id:     utils.EncodeToString(val.Hash),
			Status: val.Status,
		})
	}
	for hash, val := range s.Answers {
		expressions = append(expressions, models.Expression{
			Id:     utils.EncodeToString(hash),
			Status: "Выполнено.",
			Result: val,
		})
	}
	s.mu.Unlock()
	return expressions
}

// Получение времени выполнения операции из сервиса
func (s *OrchestratorService) OperationTime(operation rune) (time.Duration, error) {
	switch operation {
	case '+':
		return s.TimeAdditionMs, nil
	case '-':
		return s.TimeDivisionsMs, nil
	case '*':
		return s.TimeMultiplicationsMs, nil
	case '/':
		return s.TimeSubtractionMs, nil
	default:
		return 0, ErrInvalidOperation
	}
}

// Поиск новых подзадач для решения среди всех выражений
func (s *OrchestratorService) FindNewTasks(expression *models.Expression) ([]*models.Task, error) {
	tasks := []*models.Task{}
	// Если в выражении единственный элемент Linked List
	// то он будет являтся ответом на выражение
	s.mu.Lock()
	status := expression.Status
	s.mu.Unlock()
	if status == "Ошибка." {
		return nil, ErrAnswerExpression
	}
	if expression.List.Root.Next == nil {
		log.Printf("Получен ответ на задачу: %s, ответ: %f\n", expression.Value, expression.List.Root.Data.Value)
		// добавляем ответ на значение
		s.mu.Lock()
		s.Answers[expression.Hash] = expression.List.Root.Data.Value
		// удаляем, т.к. посчитали
		delete(s.Expressions, expression.Hash)
		s.mu.Unlock()
		return nil, ErrAnswerExpression
	}
	cur := expression.List.Root
	var (
		last     *customList.Node
		previous *customList.Node
	)
	// Поиск последовательности в листе формата
	// OPERATION -> FLOAT -> FLOAT
	// Такую операцию можно обработать независимо
	for cur != nil {
		if (last != nil && !last.InAction && previous != nil && !previous.InAction) && last.Data.IsOperation && !previous.Data.IsOperation && !cur.Data.IsOperation {
			operationTime, err := s.OperationTime(last.Data.Operation)
			if err != nil {
				return nil, err
			}
			s.mu.Lock()
			tasks = append(tasks, &models.Task{
				Id:             s.LastTaskId,
				FirstArgument:  cur.Data.Value,
				SecondArgument: previous.Data.Value,
				Operation:      last.Data.Operation,
				OperationTime:  operationTime,
			})
			// Помечаем родительское выражение у подзадачи
			s.TaskIdExpression[s.LastTaskId] = expression
			// Тут происходит переназначение переменных, чтобы
			// не было повторяющихся задач при параллельном запросе
			s.TaskIdUpdate[s.LastTaskId] = last
			last.InAction = true
			previous.InAction = true
			cur.InAction = true
			s.LastTaskId++
			s.mu.Unlock()
			last = nil
			previous = nil
			cur = cur.Next
		} else {
			last = previous
			previous = cur
			cur = cur.Next
		}
	}
	return tasks, nil
}

// Получение новой подзадачи из сервиса
func (s *OrchestratorService) GetTask() (*models.Task, error) {
	task := &models.Task{}
	// Итератор очередности запросов
	reqExpIdx := 0
	// Если задач нет, то их надо найти
	for len(s.Tasks) == 0 && reqExpIdx < len(s.RequestExpressions) {
		// Выражений не было совсем
		if len(s.Expressions) == 0 || reqExpIdx == len(s.Expressions) {
			return nil, ErrHaveNoTask
		} else if len(s.RequestExpressions) == 0 {
			return nil, ErrHaveNoTask
		}
		s.mu.Lock()
		// `reqExpIdx`ый по очередности запрос
		hash := s.RequestExpressions[reqExpIdx]
		s.mu.Unlock()
		tasks, err := s.FindNewTasks(s.Expressions[hash])
		if err == ErrAnswerExpression {
			s.RequestExpressions = append(s.RequestExpressions[:reqExpIdx], s.RequestExpressions[reqExpIdx+1:]...)
			continue
		} else if err != nil {
			return nil, err
		}
		s.mu.Lock()
		s.Tasks = append(s.Tasks, tasks...)
		s.mu.Unlock()
		reqExpIdx++
	}
	s.mu.Lock()
	// Достаем задачу и удаляем из списка
	if len(s.Tasks) == 0 {
		s.mu.Unlock()
		return nil, ErrHaveNoTask
	}
	task = s.Tasks[0]
	s.Tasks = s.Tasks[1:]
	s.mu.Unlock()
	return task, nil
}

// Обработка ответа на подзадачу
func (s *OrchestratorService) ProcessAnswer(taskAnswer *models.TaskResult) {
	s.mu.Lock()
	// Ищем указатель на элемент выражения в Linked List
	node := s.TaskIdUpdate[taskAnswer.Id]
	// Удаление ненужного ключа, делаем сами, т.к. нужно шарить за параллельность
	delete(s.TaskIdUpdate, taskAnswer.Id)
	// Если у нас была последовательность:
	// -> [+] -> [2] -> [2] ->
	// То при подсчете подзадачи, она должна трансформироваться в:
	// -> [4] ->
	node.Data = &customList.NodeData{
		Value: taskAnswer.Result,
	}
	node.Next = node.Next.Next.Next
	node.InAction = false
	s.mu.Unlock()
}

// Обрабатывает подзадачи с ошибками
func (s *OrchestratorService) ProcessErrorAnswer(taskAnswer *models.TaskResult) {
	s.mu.Lock()
	// Ищем указатель на элемент выражение
	expression := s.TaskIdExpression[taskAnswer.Id]
	// Удаление ненужного ключа, делаем сами, т.к. нужно шарить за параллельность
	delete(s.TaskIdUpdate, taskAnswer.Id)
	// Помечаем выражение как ошибочное
	s.Expressions[expression.Hash].Status = "Ошибка."
	s.mu.Unlock()
}
