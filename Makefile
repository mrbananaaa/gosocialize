-include .env
export

DOCKER_COMPOSE_FILE := docker-compose.yaml

.PHONY: compose-up
compose-up:
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: compose-up-build
compose-up-build:
	@docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

.PHONY: compose-restart
compose-restart:
	@docker compose -f $(DOCKER_COMPOSE_FILE) restart
.PHONY: compose-down
compose-down:
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: compose-logs
compose-logs:
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

.PHONY: psql
psql:
	@docker compose -f $(DOCKER_COMPOSE_FILE) exec -it postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: api/migrate
api/migrate:
	@echo "Migrating api"
	@goose reset && goose up

.PHONY: api/seed
api/seed:
	@echo "Seeding api"
	@cd ./backend && go run ./cmd/seed

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make compose-up"
	@echo "  make compose-up-build"
	@echo "  make compose-down"
	@echo "  make compose-logs"
	@echo "  make psql"
	@echo "  make api/migrate"
	@echo "  make api/seed"
