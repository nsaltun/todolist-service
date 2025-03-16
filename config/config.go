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
	User     string `env:"POSTGRES_USER" envDefault:"todolist"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"todolist.8080"`
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	Database string `env:"POSTGRES_DB" envDefault:"todolist"`
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
