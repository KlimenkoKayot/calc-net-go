package agent

import "fmt"

var (
	ErrLoadEnvironment       = fmt.Errorf("ошибка загрузки переменных среды")
	ErrInvalidVariableType   = fmt.Errorf("неверный тип переменной среды")
	ErrInvalidComputingValue = fmt.Errorf("число агентов (горутин) должно быть больше 0")
)
