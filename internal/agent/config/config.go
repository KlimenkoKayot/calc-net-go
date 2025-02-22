package agent

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	ComputingPower uint64 // количество горутин Агентов
}

func NewConfig() (*Config, error) {
	Port := flag.Int("port", 8080, "Port to run the Agent on")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		return nil, ErrLoadEnvironment
	}

	ComputingPowerString := os.Getenv("COMPUTING_POWER")
	if ComputingPowerString == "" {
		ComputingPowerString = "4"
	}
	ComputingPower, err := strconv.Atoi(ComputingPowerString)
	if err != nil {
		return nil, ErrInvalidVariableType
	}
	if ComputingPower < 0 {
		return nil, ErrInvalidComputingValue
	}

	return &Config{
		*Port,
		uint64(ComputingPower),
	}, nil
}
