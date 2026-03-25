package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel     string `env:"LOG_LEVEL" envDefault:"debug"`
	DbUrl        string `env:"DB_URL" envDefault:""`
	RedisUrl     string `env:"REDIS_URL" envDefault:""`
	KafkaBrokers string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	WorkerNumber int    `env:"WORKER_NUMBER" envDefault:"10"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
