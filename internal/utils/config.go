package utils

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Port          string `env:"PORT" env-default:"80"`
	Host          string `env:"HOST" env-default:"localhost"`
	DBPath        string `env:"DBPATH" env-default:"./app_test.db"`
	MigrationPath string `env:"MIGRATIONPATH" env-defaul:"./migrations"`
}

func LoadConfig() (AppConfig, error) {
	var cfg AppConfig
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return AppConfig{}, err
	}
	slog.Info("config loaded")
	return cfg, nil
}
