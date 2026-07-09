package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/reza-darius/schnibbel/internal/database"
	"github.com/reza-darius/schnibbel/internal/utils"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	utils.InitLogging()

	config, err := utils.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	_, err = database.Init(config.DBPath, config.MigrationPath)
	if err != nil {
		slog.Error("database init error")
		os.Exit(1)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide an ip address")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloWorld)

	slog.Info("listening", "addr", args[1])

	err = http.ListenAndServe(args[1], mux)
	log.Fatal(err)
}
