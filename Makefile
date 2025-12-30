include .env
export

MIGRATIONS_PATH=./migrations
DC := docker compose 
CURR_DIR_NAME  := $(notdir $(CURDIR))
PROJECT := delayed-notifier  # Имя проекта (опционально)

.PHONY: migrate-up migrate-down migrate-force run-app dc-up dc-stop dc-start

migrate-up:
	docker run --rm \
		-v $(MIGRATIONS_PATH):/migrations \
		--network host migrate/migrate \
		-path=/migrations/ -database "$(POSTGRES_DSN)" up

migrate-down:
	docker run --rm \
		-v $(MIGRATIONS_PATH):/migrations \
		--network host migrate/migrate \
		-path=$(MIGRATIONS_PATH) -database "$(POSTGRES_DSN)" down 1

migrate-force-1:
	docker run --rm \
		-v $(MIGRATIONS_PATH):/migrations \
		--network host migrate/migrate \
		-path=$(MIGRATIONS_PATH) -database "$(POSTGRES_DSN)" force 1

migrate-force-0:
	docker run --rm -v $(MIGRATIONS_PATH):/migrations \
		--network host migrate/migrate \
		-path=$(MIGRATIONS_PATH) -database "$(POSTGRES_DSN)" force 0

dc-up:
	docker compose up -d

dc-down:
	docker compose down

dc-stop:
	docker compose stop

dc-start:
	docker compose start

open-db:
	docker exec -it delayed_pg psql -U postgres -d delayed_notifier

run-app:
	go run cmd/main.go