.PHONY: migrate

migrate:
	GO_ENV=dev go run migrate/migrate.go

up:
	GO_ENV=dev go run ./cmd/main.go
