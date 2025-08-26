DOCKER_COMPOSE ?= docker compose

ENV_FILE ?= .env
DEV_DOCKER_COMPOSE_FILE ?= setup/dev/docker-compose.yaml
MIGRATIONS_PATH ?= internal/connections/postgres/migrations/

DOCKER_UP = $(DOCKER_COMPOSE) --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d
DOCKER_DOWN = $(DOCKER_COMPOSE) --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

default: ;
-include $(ENV_FILE)
export

DATABASE_URL := postgres://$(POSTGRES_USERNAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

confirm_action = @read -p "🚨 $(1) [y/N] " confirm && [ "$$confirm" = "y" ] || (echo "❌ Ação cancelada."; exit 1)

.PHONY: up build down log migrate-up migrate-down

up: COMPOSE_FILE ?= $(DEV_DOCKER_COMPOSE_FILE)
up:
	@echo "🚀 Iniciando ambiente..."
	@$(DOCKER_UP)

build: COMPOSE_FILE = $(DEV_DOCKER_COMPOSE_FILE)
build:
	@echo "🛠️  Construindo serviços..."
	@$(DOCKER_UP) --build

down: COMPOSE_FILE = $(DEV_DOCKER_COMPOSE_FILE)
down:
	@echo "🛑 Parando ambiente..."
	@$(DOCKER_DOWN)

logs:
	@echo "📜 Exibindo logs da API..."
	docker logs -f docker_restaurant_api

migrate:
	@echo "Running migrations..."
	docker exec docker_restaurant_api migrate -path=internal/connections/postgres/migrations -database ${DATABASE_URL} -verbose up

migrate-down:
	@echo "Reverting migrations..."
	docker exec docker_restaurant_api migrate -path=internal/connections/postgres/migrations -database ${DATABASE_URL} -verbose down
