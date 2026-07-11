set dotenv-load := true

run:
    go run ./cmd/*

migrate:
    goose up

dblogin:
  docker compose exec postgres psql -U {{env('POSTGRES_USER')}} -d {{env('POSTGRES_DB')}}
