package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port               int    `env:"PORT" envDefault:"8080"`
	GinMode            string `env:"GIN_MODE" envDefault:"debug"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"debug"`
	DbUrl              string `env:"DB_URL" envDefault:""`
	RedisUrl           string `env:"REDIS_URL" envDefault:""`
	JWTSecret          string `env:"JWT_SECRET" envDefault:""`
	JWTDuration        string `env:"JWT_DURATION" envDefault:"1h"`
	RateLimitCapacity  int    `env:"RL_CAPACITY" envDefault:"10"`
	RateLimitPerMinute int    `env:"RL_PER_MINUTE" envDefault:"10"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
