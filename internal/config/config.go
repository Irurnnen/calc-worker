package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ComputingPower int
}

func NewConfigExample() *Config {
	return &Config{
		Port: 8080,
	}
}

func NewConfigFromEnv() *Config {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Fatalf("Fatal error while getting config from env: COMPUTING_POWER:%s", os.Getenv("COMPUTING_POWER"))
		return NewConfigExample()
	}
	return &Config{
		ComputingPower: computingPower,
	}
}
