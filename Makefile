# Подключение переменных окружения
include .env
export $(shell sed 's/=.*//' .env)

# Константы
POSTGRES_URL := postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=$(PG_SSLMODE)
CONTAINER_NAME := postgres17
POSTGRES_IMAGE := postgres:17.5-alpine3.22

# Запуск
dev-up:
	@echo "Starting development environment with migrations..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	docker-compose up migrate
	@echo "Development environment ready!"

dev-down:
	@echo "Stopping development environment..."
	docker-compose down

dev-restart:
	@echo "Restarting development environment..."
	docker-compose restart postgres

dev-logs:
	docker-compose logs -f postgres

dev-clean:
	@echo "Cleaning development environment (removes volumes)..."
	docker-compose down -v

# PostgreSQL контейнер (старые команды, сохранены для совместимости)
postgres:
	docker run --name $(CONTAINER_NAME) -p $(PG_PORT):5432 \
		-e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) \
		-d $(POSTGRES_IMAGE)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' started successfully."

stop-postgres:
	docker stop $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' stopped successfully."

start-postgres:
	docker start $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' started successfully."

restart-postgres: stop-postgres start-postgres
	@echo "PostgreSQL container '$(CONTAINER_NAME)' restarted successfully."

remove-postgres:
	docker rm -f $(CONTAINER_NAME)
	@echo "PostgreSQL container '$(CONTAINER_NAME)' removed successfully."

# Управление базой данных
createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=$(PG_USER) --owner=$(PG_USER) $(PG_DB)
	@echo "Database '$(PG_DB)' created successfully."

dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb --username=$(PG_USER) --if-exists $(PG_DB)
	@echo "Database '$(PG_DB)' dropped successfully."

# Миграции
migrate-postgres:
	@echo "Running PostgreSQL migrations..."
	migrate -path internal/db/migrations -database "$(POSTGRES_URL)" up

rollback-postgres:
	@echo "Rolling back PostgreSQL migrations..."
	migrate -path internal/db/migrations -database "$(POSTGRES_URL)" down 1

migrate-compose:
	@echo "Running migrations via docker-compose..."
	docker-compose up migrate

# Генерация SQLC кода
sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

# Полный цикл разработки
dev-reset: dev-clean dev-up
	@echo "Development environment reset complete!"

.PHONY: dev-up dev-down dev-restart dev-logs dev-clean dev-reset \
		postgres stop-postgres start-postgres restart-postgres remove-postgres \
		createdb dropdb migrate-postgres rollback-postgres migrate-compose sqlc