package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"time"
)

type Config struct {
	Host string `env:"HOST" env-required:"true"`
	Port string `env:"PORT" env-required:"true"`

	DBUrl string `env:"DATABASE_URL" env-required:"true"`

	JWTSignedKey          string        `env:"JWT_SIGNED_KEY" env-required:"true"`
	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN" env-required:"true"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN" env-required:"true"`
}

func New(configPath string, l *slog.Logger) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		l.Error("Failed to read config", "error", err)
		return nil, err
	}

	return cfg, nil
}
