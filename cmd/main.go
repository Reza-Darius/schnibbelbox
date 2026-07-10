package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/reza-darius/schnibbel/internal/database"
	"github.com/reza-darius/schnibbel/internal/routes"
	"github.com/reza-darius/schnibbel/internal/utils"
)

func main() {
	utils.InitLogging()

	config, err := utils.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	_, err = database.InitPG()
	if err != nil {
		slog.Error("database init error")
		os.Exit(1)
	}

	routes := routes.InitRoutes()

	addr := "0.0.0.0:" + config.Port
	slog.Info("listening", "addr", addr)

	err = http.ListenAndServe(addr, routes)
	if err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
