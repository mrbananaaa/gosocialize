-include .env
export

API_ENTRY_POINT := cmd/api/main.go

DOCKER_COMPOSE_FILE := docker-compose.yaml

.PHONY: build
build:
	@echo "Building api.."
	@go build -o bin/api $(API_ENTRY_POINT)

.PHONY: run
run:
	@go run $(API_ENTRY_POINT)

.PHONY: clean
clean:
	@echo "cleaning tmp dir..."
	@rm -rf tmp

.PHONY: watch
watch:
	@air -c .air.toml

.PHONY: dev
dev: clean watch

.PHONY: windows-dev
windows-dev:
	@air -c .air.windows.toml

.PHONY: compose-up
compose-up:
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: compose-down
compose-down:
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs:
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

.PHONY: psql
psql:
	@docker compose -f $(DOCKER_COMPOSE_FILE) exec -it postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run"
	@echo "  make build"
	@echo "  make dev"
	@echo "  make windows-dev"
	@echo "  make compose-up"
	@echo "  make compose-down"
	@echo "  make compose-logs"
	@echo "  make psql"
