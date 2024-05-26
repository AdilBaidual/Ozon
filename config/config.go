package config

import (
	"Service/pkg/httpserver"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath string = "./config/config.yaml"
const ServiceName string = "ozon"

type Config struct {
	Valkey   Valkey
	Postgres Postgres `yaml:"Postgres"`
	Paseto   Paseto
	Jaeger   Jaeger            `yaml:"jaeger"`
	Server   httpserver.Config `yaml:"server"`
}

type Valkey struct {
	Host     string `env:"VALKEY_HOST" env-required:"true"`
	Port     int    `env:"VALKEY_PORT" env-required:"true"`
	Password string `env:"VALKEY_PASSWORD" env-required:"true"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `yaml:"SslMode"`
	PgDriver string `yaml:"PgDriver"`
}

type Paseto struct {
	PasetoSecret string `env:"PASETO_SECRET" env-required:"true"`
}

type Jaeger struct {
	LogSpans bool   `yaml:"log_span"`
	Host     string `env:"JAEGER_AGENT_HOST" env-required:"true"`
	Port     int    `env:"JAEGER_AGENT_PORT" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - cleanenv.ReadConfig - %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("NewConfig - cleanenv.UpdateEnv - %w", err)
	}

	return &cfg, nil
}
