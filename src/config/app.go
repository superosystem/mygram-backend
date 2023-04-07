package config

import (
	"os"
)

type AppConfig struct {
	Name    string
	Version string
	Port    string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Name:    os.Getenv("APP_NAME"),
		Version: os.Getenv("APP_VERSION"),
		Port:    os.Getenv("APP_PORT"),
	}
}
