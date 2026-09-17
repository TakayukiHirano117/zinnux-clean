.PHONY: migrate

migrate:
	GO_ENV=dev go run migrate/migrate.go
