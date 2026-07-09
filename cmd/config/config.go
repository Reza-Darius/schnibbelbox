// Package config provides config utilites
package config

import (
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Port          string `env:"PORT" env-default:"80"`
	Host          string `env:"HOST" env-default:"localhost"`
	DBPath        string `env:"DBPATH" env-default:"./app_test.db"`
	MigrationPath string `env:"MIGRATIONPATH" env-defaul:"./migrations"`
}

func LoadConfig() AppConfig {
	var cfg AppConfig
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	slog.Info("config loaded")
	return cfg
}
