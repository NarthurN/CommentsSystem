# 💬 Comments System

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![GraphQL](https://img.shields.io/badge/GraphQL-E10098?style=for-the-badge&logo=graphql)](https://graphql.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![Test Coverage](https://img.shields.io/badge/Coverage-10%25-yellow?style=for-the-badge)](#тестирование)

> Современная система комментариев с GraphQL API, real-time уведомлениями и иерархической структурой комментариев

## 📋 Описание проекта

**Comments System** - это высокопроизводительный backend-сервис для управления постами и комментариями, аналогичный системам комментирования на популярных платформах как **Хабр** или **Reddit**.

### 🎯 Назначение проекта

Проект решает следующие задачи:
- **Масштабируемая система комментирования** с поддержкой тысяч одновременных пользователей
- **Real-time взаимодействие** через WebSocket подписки
- **Производительная архитектура** с решением N+1 проблем и оптимизированными SQL запросами
- **Гибкое хранение данных** с поддержкой PostgreSQL и In-Memory режимов
- **Enterprise-grade решение** с rate limiting, централизованной обработкой ошибок и мониторингом

### 🚀 Зачем создавался проект

Система была разработана для демонстрации современных практик Go разработки:
- **Clean Architecture** с четким разделением ответственности
- **GraphQL-first подход** для гибкого API
- **Производительные решения** для высоких нагрузок
- **Comprehensive тестирование** с покрытием критических компонентов
- **Production-ready код** с Docker, мониторингом и graceful shutdown

## ✨ Ключевые возможности

### 🏗️ Архитектура
- **Clean Architecture** с dependency inversion
- **GraphQL API** с полной поддержкой запросов, мутаций и подписок
- **Два типа хранилища**: PostgreSQL и In-Memory
- **Microservices-ready** архитектура

### 🔄 Real-time функциональность
- **WebSocket подписки** для мгновенных уведомлений
- **Live комментарии** без обновления страницы
- **Pub/Sub система** для распределенной архитектуры

### 🌳 Комментарии
- **Иерархическая структура** без ограничения вложенности
- **Пагинация на всех уровнях** (корневые и дочерние комментарии)
- **Валидация контента** (до 2000 символов)
- **Управление доступом** (включение/отключение комментариев)

### ⚡ Производительность
- **Решение N+1 проблемы** через оптимизированные SQL запросы
- **Составные индексы** для быстрой пагинации
- **Rate limiting** с token bucket алгоритмом
- **Connection pooling** и graceful shutdown

### 🛡️ Безопасность и надежность
- **Централизованная обработка ошибок** с типизированными кодами
- **Rate limiting** для защиты от DDoS
- **Input validation** и sanitization
- **Health checks** для мониторинга

## 🛠️ Технологический стек

### Backend
- **Go 1.23+** - основной язык разработки
- **GraphQL** - API интерфейс с gqlgen
- **PostgreSQL 15+** - основная база данных
- **Chi Router** - HTTP маршрутизация
- **PGX v5** - PostgreSQL драйвер

### Infrastructure
- **Docker & Docker Compose** - контейнеризация
- **WebSockets** (gorilla/websocket) - real-time коммуникация
- **UUID** - уникальные идентификаторы
- **Graceful shutdown** - корректное завершение

### Разработка
- **Go Testing** - unit и integration тесты
- **Test Coverage** - мониторинг покрытия
- **Make** - автоматизация задач
- **Hot Reload** - быстрая разработка

## 📁 Структура проекта

```
CommentsSystem/
├── 📁 cmd/app/                 # 🚀 Точка входа приложения
├── 📁 internal/                # 🏢 Внутренняя бизнес-логика
│   ├── 📁 api/                 # 🌐 HTTP транспорт и обработчики
│   │   ├── handler.go          # REST API обработчики
│   │   ├── gqlgen_handler.go   # GraphQL обработчики
│   │   ├── errors.go           # ✅ Централизованная обработка ошибок
│   │   └── rate_limiter.go     # 🛡️ Rate limiting система
│   ├── 📁 config/              # ⚙️ Конфигурация приложения
│   ├── 📁 converter/           # 🔄 Конвертеры между слоями
│   ├── 📁 model/               # 🎯 Доменные модели (Domain Layer)
│   ├── 📁 repository/          # 💾 Слой доступа к данным
│   │   ├── 📁 model/           # 🗃️ Модели БД
│   │   ├── 📁 converter/       # 🔄 Domain ↔ Storage конвертеры
│   │   ├── storage.go          # 📋 Storage интерфейс
│   │   ├── postgres.go         # 🐘 PostgreSQL реализация
│   │   └── memory.go           # 🧠 In-Memory реализация
│   └── 📁 service/             # 🎮 Бизнес-логика и GraphQL
│       ├── schema.graphqls     # 📜 GraphQL схема
│       └── schema.resolvers.go # ✅ Оптимизированные резолверы
├── 📁 pkg/pubsub/              # 📡 Pub/Sub для real-time
├── 📁 migrations/              # 🗄️ SQL миграции и индексы
├── 📁 scripts/                 # 🔧 Автоматизация
├── 🐳 docker-compose.yml       # Контейнеризация
├── 🐳 Dockerfile              # Образ приложения
└── 📚 docs/                   # Документация
    └── TechnicalDocument.md   # Техническая документация
```

### 🏛️ Clean Architecture Layers

1. **🎯 Domain Layer** (`internal/model/`) - Чистые бизнес-модели без зависимостей
2. **💾 Repository Layer** (`internal/repository/`) - Доступ к данным с собственными моделями
3. **🎮 Service Layer** (`internal/service/`) - Бизнес-логика и GraphQL резолверы
4. **🌐 API Layer** (`internal/api/`) - HTTP обработчики и внешние интерфейсы
5. **🔄 Converter Layer** (`internal/converter/`) - Изоляция между слоями

## 🚀 Быстрый старт

### 📋 Системные требования

#### Минимальные требования
- **Docker** 20.10+ и **Docker Compose** 2.0+
- **Go** 1.23+ (для локальной разработки)
- **Make** (опционально, для удобства)

#### Рекомендуемые требования
- **CPU**: 2 cores
- **RAM**: 2GB
- **PostgreSQL**: 15+ (автоматически в Docker)

### 🐳 Развертывание с Docker (рекомендуется)

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/NarthurN/CommentsSystem.git
cd CommentsSystem
```

2. **Запустите приложение:**
```bash
# Запуск с PostgreSQL
docker compose up -d

# Или с In-Memory хранилищем
STORAGE_TYPE=memory docker compose up -d
```

3. **Проверьте работу:**
```bash
# Health check
curl http://localhost:8080/health

# GraphQL Playground
open http://localhost:8080
```

### 🔧 Локальная разработка

1. **Установите зависимости:**
```bash
go mod download
```

2. **Запустите PostgreSQL:**
```bash
docker compose up -d db
```

3. **Настройте окружение:**
```bash
# Создайте .env файл
make init-env

# Минимальная конфигурация
echo 'DB_DSN=postgres://user:password@localhost:5433/postsdb?sslmode=disable' > .env
echo 'STORAGE_TYPE=postgres' >> .env
```

4. **Запустите приложение:**
```bash
# Через Make
make run

# Или напрямую
go run cmd/app/main.go
```

### 🧪 Тестирование

```bash
# Все тесты
make test

# Unit тесты
make test-unit

# С покрытием
make test-coverage

# Обновить покрытие в README
make update-coverage
```

## 🌐 API Reference

### 🎯 GraphQL Endpoint

| Параметр | Значение |
|----------|----------|
| **URL** | `http://localhost:8080/graphql` |
| **Method** | `POST` |
| **Content-Type** | `application/json` |

### 🎮 GraphQL Playground

**URL**: `http://localhost:8080/`

Интерактивная среда для тестирования GraphQL запросов с автодополнением и документацией.

### 🔌 WebSocket Subscriptions

| Параметр | Значение |
|----------|----------|
| **URL** | `ws://localhost:8080/graphql` |
| **Protocol** | `WebSocket` |

### 💚 Health Check

```bash
curl http://localhost:8080/health
# Ответ: {"status": "healthy", "database": "connected"}
```

### 📝 Примеры API запросов

#### 📋 Получение постов с комментариями

```graphql
query GetPosts {
  posts(limit: 10, offset: 0) {
    id
    title
    content
    commentsEnabled
    createdAt
    comments(limit: 5, offset: 0) {
      id
      content
      parentId
      createdAt
      children(limit: 3, offset: 0) {
        id
        content
        createdAt
      }
    }
  }
}
```

#### ✍️ Создание поста

```graphql
mutation CreatePost {
  createPost(
    title: "Как работает GraphQL в Go"
    content: "Подробный разбор gqlgen и производительных решений..."
  ) {
    id
    title
    commentsEnabled
    createdAt
  }
}
```

#### 💬 Создание комментария

```graphql
mutation CreateComment {
  createComment(
    postId: "123e4567-e89b-12d3-a456-426614174000"
    content: "Отличная статья! Особенно понравилось про N+1 проблему."
    parentId: null
  ) {
    id
    content
    parentId
    createdAt
  }
}
```

#### 📡 Real-time подписка

```graphql
subscription NewComments {
  commentAdded(postId: "123e4567-e89b-12d3-a456-426614174000") {
    id
    content
    parentId
    createdAt
  }
}
```

#### 🔄 Управление комментариями

```graphql
mutation ToggleComments {
  toggleComments(
    postId: "123e4567-e89b-12d3-a456-426614174000"
    enable: false
  ) {
    id
    commentsEnabled
  }
}
```

## ⚙️ Конфигурация

### 🌍 Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `STORAGE_TYPE` | Тип хранилища: `postgres` или `memory` | `postgres` |
| `DB_DSN` | Строка подключения PostgreSQL | `postgres://user:password@localhost:5432/postsdb?sslmode=disable` |
| `HTTP_ADDR` | Адрес HTTP сервера | `:8080` |
| `GRAPHQL_INTROSPECTION` | Включить GraphQL introspection | `true` |
| `GRAPHQL_PLAYGROUND` | Включить GraphQL Playground | `true` |
| `CORS_ORIGINS` | Разрешенные CORS origins | `*` |
| `READ_TIMEOUT` | Таймаут чтения HTTP | `15s` |
| `WRITE_TIMEOUT` | Таймаут записи HTTP | `15s` |
| `SHUTDOWN_TIMEOUT` | Таймаут graceful shutdown | `30s` |

### 📊 Лимиты и производительность

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `MAX_PAGE_SIZE` | Максимальный размер страницы | `100` |
| `DEFAULT_PAGE_SIZE` | Размер страницы по умолчанию | `10` |
| `MAX_CONTENT_LENGTH` | Максимальная длина контента | `10000` |
| `MAX_TITLE_LENGTH` | Максимальная длина заголовка | `255` |
| `MAX_COMMENT_LENGTH` | Максимальная длина комментария | `2000` |

### 🐳 Docker конфигурация

```yaml
# .env для Docker
STORAGE_TYPE=postgres
DB_DSN=postgres://user:password@db:5432/postsdb?sslmode=disable
HTTP_ADDR=:8080
GRAPHQL_PLAYGROUND=true
```

## 🧪 Тестирование

### 📊 Покрытие тестами

Текущее покрытие: **10.0%**

| Пакет | Покрытие | Статус |
|-------|----------|--------|
| `internal/config` | 100% | ✅ |
| `internal/model` | 100% | ✅ |
| `pkg/pubsub` | 100% | ✅ |
| `internal/repository` | 34.9% | 🟡 |
| `internal/service` | 11.3% | 🔴 |
| `internal/api` | 17.8% | 🔴 |

### 🧪 Типы тестов

#### Unit тесты
```bash
# Быстрые модульные тесты
make test-unit

# С verbose выводом
go test -v -short ./...
```

#### Integration тесты
```bash
# Полные интеграционные тесты
make test

# Только тесты репозитория
go test -v ./internal/repository/
```

#### Coverage отчеты
```bash
# HTML отчет
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Консольный отчет
go tool cover -func=coverage.out
```

### 🎯 Тестовые сценарии

- ✅ **CRUD операции** для постов и комментариев
- ✅ **Иерархическая структура** комментариев
- ✅ **Пагинация** на всех уровнях
- ✅ **Валидация данных** и error handling
- ✅ **Concurrent access** и thread safety
- ✅ **Storage compatibility** между PostgreSQL и Memory
- ✅ **Pub/Sub функциональность** и WebSocket
- ✅ **Rate limiting** и security features

## 🔧 Команды разработки

### 🏗️ Сборка и запуск

```bash
# Сборка
make build

# Запуск
make run

# Очистка
make clean
```

### 🐳 Docker операции

```bash
# Сборка образов
make docker-build

# Запуск контейнеров
make docker-up

# Остановка
make docker-down

# Логи
make docker-logs
```

### 🛠️ Разработка

```bash
# Установка зависимостей
make deps

# Линтинг
make lint

# Форматирование
make fmt

# Генерация кода
make generate
```

### 📊 Мониторинг

```bash
# Health check
curl http://localhost:8080/health

# Metrics (если настроены)
curl http://localhost:8080/metrics

# Database status
docker compose exec db pg_isready -U user
```

## 🎯 Планы по доработке

### 🚀 Краткосрочные планы (1-2 месяца)

#### 📈 Повышение покрытия тестами
- **Цель**: Поднять общее покрытие с 10% до 85%+
- **Метод**: Добавить unit тесты для GraphQL резолверов и API обработчиков
- **Результат**: Повышение надежности и confidence в рефакторинге

#### 🔍 Observability и мониторинг
- **Добавить OpenTelemetry** для distributed tracing
- **Интегрировать Prometheus metrics** для мониторинга производительности
- **Настроить structured logging** с JSON форматом
- **Результат**: Production-ready observability для высоких нагрузок

#### 🛡️ Улучшение безопасности
- **Добавить JWT аутентификацию** для защиты API
- **Реализовать rate limiting на GraphQL operation level**
- **Добавить input sanitization** против XSS атак
- **Результат**: Enterprise-grade безопасность

### 🌟 Долгосрочные планы (3-6 месяцев)

#### 🗄️ Расширение хранилища
- **Добавить Redis кэширование** для популярных постов и комментариев
- **Реализовать distributed rate limiting** через Redis
- **Добавить поддержку MongoDB** как альтернативного NoSQL хранилища
- **Результат**: Flexible storage options и improved performance

#### 🚀 Производительность и масштабирование
- **Реализовать connection pooling** с различными стратегиями
- **Добавить GraphQL query complexity analysis** для предотвращения expensive queries
- **Внедрить database connection balancing** для read replicas
- **Результат**: Поддержка высоких нагрузок (10k+ concurrent users)

#### 🎮 Расширенная функциональность
- **Система реакций** (лайки, дизлайки) с real-time счетчиками
- **Модерация контента** с автоматическим спам-фильтром
- **Mentions система** (@username) с уведомлениями
- **Threading и Reply-to-specific-comment** для улучшенной навигации
- **Результат**: Feature parity с современными платформами

#### 📱 API расширения
- **REST API endpoints** для mobile applications
- **Webhooks система** для интеграции с внешними сервисами
- **GraphQL Federation** для microservices архитектуры
- **Результат**: Comprehensive API ecosystem

#### 🔒 Enterprise функции
- **Multi-tenancy поддержка** для SaaS решений
- **Advanced analytics** для content creators
- **Content versioning** и edit history
- **Backup и disaster recovery** процедуры
- **Результат**: Enterprise-ready solution

### 🎛️ Техническая модернизация

#### 🏗️ Архитектурные улучшения
- **Микросервисная декомпозиция** (Posts Service, Comments Service, Notifications Service)
- **Event Sourcing** для audit trail и восстановления состояния
- **CQRS pattern** для разделения read/write операций
- **Результат**: Scalable distributed architecture

#### ⚡ Производительность БД
- **Database sharding** по post_id для горизонтального масштабирования
- **Materialized views** для сложных аналитических запросов
- **Query optimization** с автоматическим query plan analysis
- **Результат**: Linear scalability для massive datasets

### 📅 Временная шкала

| Период | Приоритет | Фокус |
|--------|-----------|--------|
| **Месяц 1** | 🔴 Высокий | Тестирование + Observability |
| **Месяц 2** | 🔴 Высокий | Безопасность + Мониторинг |
| **Месяц 3** | 🟡 Средний | Redis + Performance |
| **Месяц 4** | 🟡 Средний | Расширенная функциональность |
| **Месяц 5-6** | 🟢 Низкий | Microservices + Enterprise |

## 📚 Дополнительная информация

### 📖 Документация

- 🏛️ **[Clean Architecture Guide](CLEAN_ARCHITECTURE.md)** - Подробное описание архитектурных принципов
- 📋 **[Technical Document](TechnicalDocument.md)** - Техническая спецификация и решения
- 🎯 **[Task Requirements](Task.md)** - Исходные требования проекта
- 🔧 **[API Documentation](http://localhost:8080/)** - GraphQL Playground с интерактивной документацией

### 🛠️ Инструменты разработки

#### 🔍 Статический анализ
```bash
# Линтинг с golangci-lint
golangci-lint run

# Проверка безопасности
gosec ./...

# Проверка зависимостей
go mod tidy && go mod verify
```

#### 🐛 Отладка
```bash
# Запуск с отладкой
dlv debug cmd/app/main.go

# Профилирование
go tool pprof http://localhost:8080/debug/pprof/profile
```

#### 📊 Мониторинг производительности
```bash
# Бенчмарки
go test -bench=. -benchmem ./...

# Memory профилирование
go tool pprof heap.prof
```

### 🤝 Участники разработки

| Роль | Участник | Вклад |
|------|----------|-------|
| **Lead Developer** | NarthurN | Architecture, Core Implementation |
| **Backend Engineer** | NarthurN | GraphQL, Database Design |
| **DevOps Engineer** | NarthurN | Docker, CI/CD Setup |

### 🎯 Используемые паттерны

- **🏛️ Clean Architecture** - Слоистая архитектура с dependency inversion
- **🎭 Repository Pattern** - Абстракция доступа к данным
- **🏭 Factory Pattern** - Создание storage implementations
- **🔄 Adapter Pattern** - Converter слой между доменами
- **📡 Observer Pattern** - Pub/Sub для real-time уведомлений
- **🛡️ Strategy Pattern** - Различные стратегии хранения
- **🎮 Command Pattern** - GraphQL мутации как команды

### 🔗 Полезные ссылки

#### 📚 Образовательные ресурсы
- [Effective Go](https://go.dev/doc/effective_go) - Best practices Go
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/) - GraphQL гайды
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Uncle Bob's Clean Architecture

#### 🛠️ Инструменты
- [gqlgen](https://github.com/99designs/gqlgen) - GraphQL code generation
- [Chi Router](https://github.com/go-chi/chi) - HTTP router
- [PGX](https://github.com/jackc/pgx) - PostgreSQL driver
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

#### 🌟 Похожие проекты
- [GraphQL Go](https://github.com/graphql-go/graphql) - GraphQL implementation
- [Gorm](https://github.com/go-gorm/gorm) - ORM library
- [Gin](https://github.com/gin-gonic/gin) - Web framework

---

<div align="center">

### 🎯 Готов к продакшену | 🚀 Enterprise-grade | 💎 Production-ready

**Создано с ❤️ используя Go и GraphQL**

</div>
