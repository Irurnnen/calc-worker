package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ComputingPower  int
	OrchestratorURL string
}

func NewConfigExample() *Config {
	return &Config{
		ComputingPower:  1,
		OrchestratorURL: "127.0.0.1:8080",
	}
}

func NewConfigFromEnv() *Config {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	orchestratorPath := os.Getenv("ORCHESTRATOR_URL")
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: COMPUTING_POWER:%s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	return &Config{
		ComputingPower:  computingPower,
		OrchestratorURL: orchestratorPath,
	}
}
