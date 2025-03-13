package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
	PostgresConfig
}

type PostgresConfig struct {
	PostgresUrl string `env:"POSTGRES_DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/todolist?sslmode=disable"`
}

// NewAppConfig returns a new instance of AppConfig
func NewAppConfig() *AppConfig {
	var c AppConfig
	err := env.Parse(&c)
	if err != nil {
		log.Fatalf("error loading app config: %v", err)
	}

	return &c
}
