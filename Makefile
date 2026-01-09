DC := docker compose
PROJECT := delayed notifier
APP_SERVICE := delayed notifier
MIGRATIONS_DIR := ./migrations
LOCAL_DSN := "postgres://postgres:pass@localhost:5432/delayed_notifier?sslmode=disable"

.PHONY: help
help: ## Список команд
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Запуск unit-тестов
	go test -v ./internal/... ./pkg/... ./cmd/... ./config/...

.PHONY: start
start: test ## Запуск всего проекта в Docker (DB + Migrations + App)
	$(DC) up -d postgres redis --wait
	$(MAKE) migrate-up
	$(DC) --profile app up -d --build

.PHONY: start-local
start-local: up migrate-up run ## Запуск локально (DB в Docker, App через go run)

.PHONY: run
run: ## Запуск go run
	go run ./cmd/main.go

# --- Database & Migrations ---

.PHONY: migrate-up
migrate-up: ## Накатить миграции (up)
	migrate -path $(MIGRATIONS_DIR) -database $(LOCAL_DSN) up

.PHONY: migrate-down
migrate-down: ## Откатить миграции (down)
	migrate -path $(MIGRATIONS_DIR) -database $(LOCAL_DSN) down

.PHONY: migrate-force
migrate-force: ## Форсировать версию миграции (make migrate-force v=1)
	migrate -path $(MIGRATIONS_DIR) -database $(LOCAL_DSN) force $(v)

# --- Docker Control ---

.PHONY: up
up: ## Поднять инфраструктуру (db, redis)
	$(DC) -p $(PROJECT) up -d --wait

.PHONY: down
down: ## Остановить и удалить контейнеры
	$(DC) -p $(PROJECT) down

.PHONY: stop
stop: ## Остановить контейнеры (без удаления)
	$(DC) -p $(PROJECT) stop

.PHONY: start-containers
start-containers: ## Запустить остановленные контейнеры
	$(DC) -p $(PROJECT) start

.PHONY: clean
clean: ## Удалить всё (volumes, images)
	$(DC) down -v && docker system prune -f

# --- Logs & Utils ---

.PHONY: logs
logs: ## Видеть логи всех сервисов
	$(DC) logs -f

.PHONY: logs-app
logs-app: ## Логи сервиса приложения
	$(DC) logs -f $(APP_SERVICE)

.PHONY: ps
ps: ## Статус контейнеров
	$(DC) ps

.PHONY: exec-app
exec-app: ## Зайти в контейнер приложения (bash)
	$(DC) exec $(APP_SERVICE) /bin/sh
