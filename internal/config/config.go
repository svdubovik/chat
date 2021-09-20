package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var defaultStrValues = map[string]string{
	"LOG_LEVEL":    "info",
	"LOG_FORMAT":   "json",
	"BIND_ADDRESS": ":8000",
}

type Config struct {
	LogLevel    string
	LogFormat   string
	BindAddress string
	Service     string
}

func getStrEnv(service, envVarName string) string {
	envVar, exists := os.LookupEnv(fmt.Sprintf("%s_%s", strings.ToUpper(service), envVarName))
	if !exists {
		envVar = defaultStrValues[envVarName]
	}

	return envVar
}

func NewConfig(service string) *Config {
	godotenv.Load()

	return &Config{
		LogLevel:    getStrEnv(service, "LOG_LEVEL"),
		LogFormat:   getStrEnv(service, "LOG_FORMAT"),
		BindAddress: getStrEnv(service, "BIND_ADDRESS"),
		Service:     service,
	}
}
