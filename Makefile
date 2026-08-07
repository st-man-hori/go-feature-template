.DEFAULT_GOAL := help

.PHONY: help up down build rebuild bash migrate migrate-down fresh seed test format format-test lint openapi ci

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

up: ## Start the local environment
	docker compose up -d

down: ## Stop the local environment
	docker compose down

build: ## Build (or rebuild) the images
	docker compose build

rebuild: ## Rebuild images from scratch and restart
	docker compose down
	docker compose build --no-cache
	docker compose up -d

bash: ## Open a shell in the app container
	docker compose exec app sh

migrate: ## Run migrations
	docker compose exec app migrate -path migrations -database "mysql://$${DB_USERNAME:-app}:$${DB_PASSWORD:-secret}@tcp(mysql:3306)/$${DB_DATABASE:-app}" up

migrate-down: ## Roll back the last migration
	docker compose exec app migrate -path migrations -database "mysql://$${DB_USERNAME:-app}:$${DB_PASSWORD:-secret}@tcp(mysql:3306)/$${DB_DATABASE:-app}" down 1

fresh: ## Drop everything and re-run migrations
	docker compose exec app migrate -path migrations -database "mysql://$${DB_USERNAME:-app}:$${DB_PASSWORD:-secret}@tcp(mysql:3306)/$${DB_DATABASE:-app}" drop -f
	$(MAKE) migrate

seed: ## Run database seeders
	docker compose exec app go run ./cmd/seed

test: ## Run the test suite
	docker compose exec app go test ./...

format: ## Fix formatting with gofmt + goimports
	docker compose exec app sh -c "gofmt -w . && goimports -w ."

format-test: ## Check formatting without fixing
	docker compose exec app sh -c 'test -z "$$(gofmt -l .)"'

lint: ## Run static analysis with golangci-lint
	docker compose exec app golangci-lint run

openapi: ## Regenerate the OpenAPI schema (docs/) via swaggo
	docker compose exec app swag init -g cmd/api/main.go -o docs --parseDependency

ci: format-test lint test ## Run the same checks as CI
