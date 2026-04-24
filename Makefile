SHELL := /bin/bash
.DEFAULT_GOAL := help

APP_NAME := delayedNotifier
COMPOSE := docker compose
GO := go
APP_PKG := ./cmd/main
GOLANGCI_LINT := golangci-lint

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Dev:"
	@echo "  make up              Start services"
	@echo "  make down            Stop everything"
	@echo "  make logs            Follow logs"
	@echo "  make ps              Show containers"
	@echo "  make restart         Restart service"
	@echo ""
	@echo "Migrations:"
	@echo "  make migrate-up      Run migrations (via migrator service)"
	@echo "  make migrate-down    Rollback 1 migration (via migrator service)"
	@echo ""
	@echo "Go:"
	@echo "  make tidy            go mod tidy"
	@echo "  make fmt             gofmt"
	@echo "  make test            go test ./..."
	@echo "  make build           Build service locally into ./bin/"
	@echo ""
	@echo "Lint/format:"
	@echo "  make lint            Run golangci-lint"
	@echo "  make lint-fix        Run golangci-lint with --fix"
	@echo "  make fmt-ci          Run golangci-lint fmt"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build    Build service image"

.PHONY: up
up:
	$(COMPOSE) up -d --build

.PHONY: down
down:
	$(COMPOSE) down -v

.PHONY: ps
ps:
	$(COMPOSE) ps

.PHONY: logs
logs:
	$(COMPOSE) logs -f --tail=200

.PHONY: restart
restart:
	$(COMPOSE) restart $(APP_NAME)

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

generate:
	go generate ./...

mockery-install:
	go install github.com/vektra/mockery/v3@v3.2.5

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: test
test:
	$(GO) test ./...

.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=0 $(GO) build -o bin/$(APP_NAME) $(APP_PKG)

.PHONY: docker-build
docker-build:
	$(COMPOSE) build $(APP_NAME)

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix:
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: fmt-ci
fmt-ci:
	$(GOLANGCI_LINT) fmt ./...

oapi-install:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
