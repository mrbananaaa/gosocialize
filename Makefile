-include .env
export

API_ENTRY_POINT := cmd/api/main.go

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

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run"
	@echo "  make build"
	@echo "  make dev"
	@echo "  make windows-dev"
