.PHONY: build run test clean docker-build docker-up docker-down migrate generate-gql

# Переменные
BINARY_NAME=server
DOCKER_COMPOSE=docker-compose

# Сборка приложения
build:
	go build -o $(BINARY_NAME) ./cmd/app

# Запуск приложения
run: build
	./$(BINARY_NAME)

# Запуск тестов
test:
	go test -v ./...

# Запуск только unit тестов (без интеграционных)
test-unit:
	go test -v -short ./...

# Запуск с покрытием
test-coverage:
	go test -v -cover ./...

# Обновление покрытия в README.md
update-coverage:
	@./scripts/update_coverage.sh

# Создание .env файла из шаблона
init-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		cp env.example .env; \
		echo ".env file created! Please edit it with your configuration."; \
	else \
		echo ".env file already exists. Remove it first if you want to recreate."; \
	fi

# Демонстрация разных типов хранилища
demo-storage:
	@echo "Running storage types demonstration..."
	@./demo_storage.sh

# Очистка
clean:
	go clean
	rm -f $(BINARY_NAME)

# Docker команды
docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Запуск миграций
migrate:
	@echo "Running migrations..."
	@docker exec -it commentssystem_db_1 psql -U user -d postsdb -f /docker-entrypoint-initdb.d/001_init_schema.sql

# Генерация GraphQL кода
generate-gql:
	@echo "Generating GraphQL code with gqlgen..."
	@echo "Note: gqlgen must be installed: go install github.com/99designs/gqlgen@latest"
	@echo "Skipping generation as code was already generated during initial setup"

# Генерация кода ogen (требует установленный Go и ogen)
generate:
	go generate ./...

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Линтер
lint:
	golangci-lint run

# Форматирование кода
fmt:
	go fmt ./...

# Помощь
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run only unit tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  update-coverage- Update test coverage in README.md"
	@echo "  init-env       - Create .env file from template"
	@echo "  demo-storage   - Demonstrate different storage types"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Start containers"
	@echo "  docker-down    - Stop containers"
	@echo "  docker-logs    - View container logs"
	@echo "  migrate        - Run database migrations"
	@echo "  generate-gql   - Generate GraphQL code with gqlgen"
	@echo "  generate       - Generate code using go generate"
	@echo "  deps           - Install dependencies"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  help           - Show this help"
