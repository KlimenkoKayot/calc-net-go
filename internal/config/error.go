package config

import "fmt"

var (
	ErrLoadEnvironment       = fmt.Errorf("ошибка загрузки переменных среды")
	ErrInvalidVariableType   = fmt.Errorf("неверный тип переменной среды")
	ErrInvalidTime           = fmt.Errorf("время не может быть отрицательным значением")
	ErrInvalidComputingValue = fmt.Errorf("число агентов (горутин) должно быть больше 0")
)
