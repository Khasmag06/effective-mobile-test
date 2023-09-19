package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP      HTTPConfig
	PG        PGConfig
	Redis     RedisConfig
	Logger    LoggerConfig
	PersonApi PersonApiConfig
	Kafka     KafkaConfig
}

type (
	HTTPConfig struct {
		Host    string `env:"HTTP_HOST" yaml:"host"`
		Port    string `env:"HTTP_PORT" yaml:"port"`
		Address string
	}

	PGConfig struct {
		URL      string
		User     string `env:"POSTGRES_USER"     yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD" yaml:"password"`
		Host     string `env:"POSTGRES_HOST"     yaml:"host"`
		Port     uint16 `env:"POSTGRES_PORT"     yaml:"port"`
		DB       string `env:"POSTGRES_DB"       yaml:"DB"`
		SSLMode  string `env:"POSTGRES_SSL_MODE" yaml:"SSLMode"`
	}

	RedisConfig struct {
		Host     string `env:"REDIS_HOST"     yaml:"host"`
		Port     string `env:"REDIS_PORT"     yaml:"port"`
		Password string `env:"REDIS_PASSWORD" yaml:"password"`
		DB       int    `env:"REDIS_DB"       yaml:"DB"`
	}

	LoggerConfig struct {
		LogFilePath string `env:"LOG_FILE_PATH" yaml:"logFilePath"`
		Level       string `env:"LOG_LVL"       yaml:"level"`
	}

	PersonApiConfig struct {
		AgeURL         string `env:"AGE_API_URL"    yaml:"ageURL"`
		GenderURL      string `env:"GENDER_API_URL" yaml:"genderURL"`
		NationalityURL string `env:"NATION_API_URL" yaml:"nationalityURL"`
	}

	KafkaConfig struct {
		BrokerURLs     []string
		Broker         string `env:"KAFKA_BROKER"             yaml:"broker"`
		FioTopic       string `env:"KAFKA_FIO_TOPIC"          yaml:"fioTopic"`
		FioFailedTopic string `env:"KAFKA_FIO_FAILED_TOPIC"   yaml:"fioFailedTopic"`
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}
	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.User, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DB, cfg.PG.SSLMode)
	cfg.Kafka.BrokerURLs = []string{cfg.Kafka.Broker}

	return cfg, nil
}
