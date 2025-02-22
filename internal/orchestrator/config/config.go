package orchestrator

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  int    // порт для запуска сервера
	TimeAdditionMs        uint64 // время выполнения операции сложения в миллисекундах
	TimeSubtractionMs     uint64 // время выполнения операции вычитания в миллисекундах
	TimeMultiplicationsMs uint64 // время выполнения операции умножения в миллисекундах
	TimeDivisionsMs       uint64 // время выполнения операции деления в миллисекундах
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, ErrLoadEnvironment
	}

	PortString := os.Getenv("PORT")
	if PortString == "" {
		PortString = "8080"
	}
	Port, err := strconv.Atoi(PortString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if Port < 0 {
		return nil, ErrInvalidPort
	}

	TimeAdditionString := os.Getenv("TIME_ADDITION_MS")
	if TimeAdditionString == "" {
		TimeAdditionString = "0"
	}
	TimeAdditionMs, err := strconv.Atoi(TimeAdditionString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeAdditionMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeSubtractionString := os.Getenv("TIME_SUBTRACTION_MS")
	if TimeSubtractionString == "" {
		TimeSubtractionString = "0"
	}
	TimeSubtractionMs, err := strconv.Atoi(TimeSubtractionString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeSubtractionMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeMultiplicationsString := os.Getenv("TIME_MULTIPLICATIONS_MS")
	if TimeMultiplicationsString == "" {
		TimeMultiplicationsString = "0"
	}
	TimeMultiplicationsMs, err := strconv.Atoi(TimeMultiplicationsString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeMultiplicationsMs < 0 {
		return nil, ErrInvalidTime
	}

	TimeDivisionsString := os.Getenv("TIME_DIVISIONS_MS")
	if TimeDivisionsString == "" {
		TimeDivisionsString = "0"
	}
	TimeDivisionsMs, err := strconv.Atoi(TimeDivisionsString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if TimeDivisionsMs < 0 {
		return nil, ErrInvalidTime
	}

	return &Config{
		Port,
		uint64(TimeAdditionMs),
		uint64(TimeSubtractionMs),
		uint64(TimeMultiplicationsMs),
		uint64(TimeDivisionsMs),
	}, nil
}
