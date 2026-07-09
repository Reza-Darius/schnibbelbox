addr := "127.0.0.1:4000"

run:
    go run ./cmd/main.go {{ addr }}

migrate:
    goose up
