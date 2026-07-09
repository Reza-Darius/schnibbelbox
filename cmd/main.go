package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/reza-darius/schnibbel/cmd/config"
	utils "github.com/reza-darius/schnibbel/internal"
	"github.com/reza-darius/schnibbel/internal/database"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("hello world"))
	if err != nil {
		log.Println(err)
	}
}

const MigrationPath = "./migrations"
const DBPath = "./database.db"

func main() {
	utils.InitLogging()
	config := config.LoadConfig()

	database.Init(config.DBPath, config.MigrationPath)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("please provide an ip address")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloWorld)

	slog.Info("listening", "addr", args[1])

	err := http.ListenAndServe(args[1], mux)
	log.Fatal(err)
}
