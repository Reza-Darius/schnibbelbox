package utils

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Port          string `env:"PORT" env-default:"80"`
	DBPath        string `env:"DB_PATH" env-default:"./app_test.db"`
	MigrationPath string `env:"MIGRATION_PATH" env-defaul:"./migrations"`
}

func LoadConfig() (AppConfig, error) {
	var cfg AppConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return AppConfig{}, err
	}
	slog.Info("config loaded")
	return cfg, nil
}
