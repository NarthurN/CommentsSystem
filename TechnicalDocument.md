# Технический документ: Система постов и комментариев

## 📋 Содержание

1. [Обзор проекта](#обзор-проекта)
2. [Архитектура](#архитектура)
3. [Структура проекта](#структура-проекта)
4. [API эндпоинты](#api-эндпоинты)
5. [База данных](#база-данных)
6. [Конфигурация](#конфигурация)
7. [Развертывание](#развертывание)
8. [Тестирование](#тестирование)

---

## 🎯 Обзор проекта

### Цель
Разработка бэкенд-системы для управления постами и комментариями с использованием GraphQL API, поддерживающей иерархическую структуру комментариев и real-time обновления.

### Ключевые возможности
- ✅ Создание и чтение постов
- ✅ Иерархические комментарии (вложенность)
- ✅ Real-time подписки через WebSocket
- ✅ Пагинация для постов и комментариев
- ✅ Валидация данных
- ✅ Возможность отключения комментариев для поста
- ✅ GraphQL Playground для интерактивного тестирования
- ✅ Health check эндпоинт
- ✅ Clean Architecture с четким разделением слоев

### Технологический стек
| Компонент | Технология | Версия |
|-----------|------------|--------|
| **Язык программирования** | Go | 1.22+ |
| **Web Framework** | Chi Router | v5 |
| **API** | GraphQL | graphql-go/graphql |
| **База данных** | PostgreSQL | 15 |
| **Драйвер БД** | pgx/v5 | v5.5.5 |
| **WebSockets** | gorilla/websocket | - |
| **Контейнеризация** | Docker + Docker Compose | - |
| **Архитектура** | Clean Architecture | - |
| **Документация** | GraphQL Playground | - |

---

## 🏗️ Архитектура

### Clean Architecture

Проект следует принципам чистой архитектуры с четким разделением слоев и использованием конвертеров для изоляции между слоями:

```
┌─────────────────────────────────────────────────────────────┐
│                    API Layer (HTTP/WebSocket)               │
│                          ↕ converters                       │
├─────────────────────────────────────────────────────────────┤
│                  Service Layer (Business Logic)             │
│                          ↕ interfaces                       │
├─────────────────────────────────────────────────────────────┤
│                Repository Layer (Data Access)               │
│                          ↕ converters                       │
├─────────────────────────────────────────────────────────────┤
│                    Domain Layer (Models)                    │
└─────────────────────────────────────────────────────────────┘
```

### Слои архитектуры

#### 1. **Domain Layer** (`internal/model/`)
- **Назначение**: Чистые бизнес-модели без внешних зависимостей
- **Компоненты**:
  - `entities.go` - Доменные сущности (Post, Comment, PostWithComments, CommentTree)
  - Бизнес-правила и валидация на уровне модели
  - Методы доменной логики (IsValidTitle, IsValidContent, IsValidComment)

#### 2. **Repository Layer** (`internal/repository/`)
- **Назначение**: Абстракция доступа к данным с собственными моделями
- **Компоненты**:
  - `storage.go` - Интерфейсы репозиториев
  - `postgres.go` - PostgreSQL реализация
  - `model/storage_models.go` - Модели данных для базы
  - `converter/converter.go` - Конвертеры между доменными и БД моделями

#### 3. **Service Layer** (`internal/service/`)
- **Назначение**: Бизнес-логика и GraphQL резолверы
- **Компоненты**:
  - `service.go` - GraphQL схема и сервис
  - `resolvers.go` - Резолверы для Query, Mutation, Subscription
  - Валидация бизнес-правил
  - Интеграция с Pub/Sub
  - Health check логика

#### 4. **API Layer** (`internal/api/`)
- **Назначение**: Обработка HTTP-запросов и WebSocket соединений
- **Компоненты**:
  - `handler.go` - GraphQL HTTP хендлер и WebSocket
  - CORS middleware
  - GraphQL Playground
  - Health check endpoint

#### 5. **Converter Layer** (`internal/converter/`)
- **Назначение**: Изоляция между слоями API и Domain
- **Компоненты**:
  - `api_converter.go` - Конвертеры для API (GraphQL)
  - GraphQLConverter - Конвертация в GraphQL представление
  - ValidationConverter - Валидация и парсинг входных данных

---

## 📁 Структура проекта

```
CommentsSystem/
├── 📁 cmd/                          # Точки входа приложения
│   └── 📁 app/
│       └── 📄 main.go               # Главный файл приложения
│
├── 📁 internal/                     # Внутренняя логика (не экспортируется)
│   ├── 📁 api/                      # HTTP обработчики
│   │   └── 📄 handler.go            # GraphQL и WebSocket хендлеры
│   │
│   ├── 📁 config/                   # Конфигурация
│   │   └── 📄 config.go             # Загрузка переменных окружения
│   │
│   ├── 📁 model/                    # Доменные модели (Domain Layer)
│   │   └── 📄 entities.go           # Чистые доменные сущности
│   │
│   ├── 📁 converter/                # Конвертеры API ↔ Domain
│   │   └── 📄 api_converter.go      # GraphQL конвертеры
│   │
│   ├── 📁 repository/               # Слой доступа к данным
│   │   ├── 📄 storage.go            # Интерфейсы репозиториев
│   │   ├── 📄 postgres.go           # PostgreSQL реализация
│   │   ├── 📄 postgres_test.go      # Интеграционные тесты
│   │   ├── 📁 model/                # Модели данных для БД
│   │   │   └── 📄 storage_models.go # PostDB, CommentDB и др.
│   │   └── 📁 converter/            # Конвертеры Domain ↔ Storage
│   │       └── 📄 converter.go      # PostConverter, CommentConverter
│   │
│   └── 📁 service/                  # Бизнес-логика
│       ├── 📄 service.go            # GraphQL схема и сервис
│       ├── 📄 resolvers.go          # GraphQL резолверы
│       └── 📄 service_test.go       # Unit тесты
│
├── 📁 pkg/                          # Публичные пакеты
│   └── 📁 pubsub/                   # Pub/Sub система
│       └── 📄 pubsub.go             # In-memory Pub/Sub
│
├── 📁 migrations/                   # Миграции базы данных
│   └── 📄 001_init_schema.sql      # Инициализация схемы
│
├── 📄 go.mod                        # Go модуль и зависимости
├── 📄 go.sum                        # Хеши зависимостей
├── 📄 Dockerfile                    # Многостадийная сборка
├── 📄 docker-compose.yml           # Оркестрация контейнеров
├── 📄 Makefile                     # Команды для разработки
├── 📄 .gitignore                   # Исключения Git
├── 📄 README.md                    # Документация проекта
├── 📄 CLEAN_ARCHITECTURE.md        # Детальное описание архитектуры
└── 📄 TechnicalDocument.md         # Технический документ
```

### Подробное описание файлов

#### **Архитектурные слои**

| Слой | Файлы | Назначение |
|------|-------|------------|
| **Domain Layer** | `internal/model/entities.go` | Чистые бизнес-модели без зависимостей |
| **Repository Layer** | `internal/repository/*.go` | Доступ к данным с собственными моделями |
| **Service Layer** | `internal/service/*.go` | Бизнес-логика и GraphQL резолверы |
| **API Layer** | `internal/api/*.go` | HTTP обработчики и внешние интерфейсы |
| **Converter Layer** | `internal/converter/*.go` | Изоляция между слоями через конвертеры |

#### **Основные компоненты**

| Компонент | Файл | Назначение |
|-----------|------|------------|
| **Точка входа** | `cmd/app/main.go` | Инициализация всех компонентов, запуск сервера |
| **Конфигурация** | `internal/config/config.go` | Загрузка переменных окружения |
| **Доменные модели** | `internal/model/entities.go` | Post, Comment, PostWithComments, CommentTree |
| **API конвертеры** | `internal/converter/api_converter.go` | GraphQL конвертеры |
| **Интерфейсы хранилища** | `internal/repository/storage.go` | Storage, PostRepository, CommentRepository |
| **PostgreSQL реализация** | `internal/repository/postgres.go` | Конкретная реализация для PostgreSQL |
| **Модели БД** | `internal/repository/model/storage_models.go` | PostDB, CommentDB для базы данных |
| **Конвертеры БД** | `internal/repository/converter/converter.go` | Domain ↔ Storage конвертация |
| **GraphQL схема** | `internal/service/service.go` | Определение типов и резолверов |
| **Резолверы** | `internal/service/resolvers.go` | Бизнес-логика GraphQL |
| **HTTP хендлеры** | `internal/api/handler.go` | Обработка запросов |
| **Pub/Sub** | `pkg/pubsub/pubsub.go` | Система подписок |

---

## 🔌 API эндпоинты

### HTTP эндпоинты

| Метод | Путь | Назначение | Описание |
|-------|------|------------|----------|
| `GET` | `/` | GraphQL Playground | Веб-интерфейс для тестирования GraphQL |
| `POST` | `/graphql` | GraphQL API | Основной эндпоинт для GraphQL запросов |
| `WebSocket` | `/subscriptions` | WebSocket | Подключение для real-time подписок |
| `GET` | `/health` | Health Check | Проверка состояния сервиса |

### GraphQL API

#### **Query (Запросы)**

##### `posts(limit: Int, offset: Int): [Post!]!`
Получение списка постов с пагинацией.

**Параметры:**
- `limit` (Int, по умолчанию: 10) - количество постов
- `offset` (Int, по умолчанию: 0) - смещение

**Пример запроса:**
```graphql
query {
  posts(limit: 5, offset: 0) {
    id
    title
    content
    commentsEnabled
    createdAt
    comments {
      id
      content
      parentId
      createdAt
    }
  }
}
```

##### `post(id: ID!): Post`
Получение конкретного поста по ID.

**Параметры:**
- `id` (ID!) - уникальный идентификатор поста

**Пример запроса:**
```graphql
query {
  post(id: "35d67a04-2829-4380-8f82-bcfdf8e5ca16") {
    id
    title
    content
    commentsEnabled
    createdAt
    comments {
      id
      content
      parentId
      createdAt
      children {
        id
        content
        createdAt
      }
    }
  }
}
```

#### **Mutation (Мутации)**

##### `createPost(title: String!, content: String!): Post!`
Создание нового поста.

**Параметры:**
- `title` (String!) - заголовок поста
- `content` (String!) - содержимое поста

**Валидация:**
- Заголовок: от 1 до 255 символов
- Содержимое: от 1 до 10000 символов

**Пример запроса:**
```graphql
mutation {
  createPost(
    title: "Мой первый пост"
    content: "Содержимое поста"
  ) {
    id
    title
    content
    commentsEnabled
    createdAt
  }
}
```

##### `createComment(postId: ID!, parentId: ID, content: String!): Comment!`
Создание комментария к посту.

**Параметры:**
- `postId` (ID!) - ID поста
- `parentId` (ID, опционально) - ID родительского комментария
- `content` (String!) - текст комментария (макс. 2000 символов)

**Валидация:**
- Длина комментария не более 2000 символов
- Комментарии должны быть разрешены для поста
- Пост должен существовать

**Пример запроса:**
```graphql
mutation {
  createComment(
    postId: "35d67a04-2829-4380-8f82-bcfdf8e5ca16"
    content: "Отличный пост!"
  ) {
    id
    content
    parentId
    createdAt
  }
}
```

##### `toggleComments(postId: ID!, enable: Boolean!): Post!`
Включение/отключение комментариев для поста.

**Параметры:**
- `postId` (ID!) - ID поста
- `enable` (Boolean!) - true для включения, false для отключения

**Пример запроса:**
```graphql
mutation {
  toggleComments(
    postId: "35d67a04-2829-4380-8f82-bcfdf8e5ca16"
    enable: false
  ) {
    id
    title
    commentsEnabled
  }
}
```

#### **Subscription (Подписки)**

##### `commentAdded(postId: ID!): Comment!`
Подписка на новые комментарии к посту.

**Параметры:**
- `postId` (ID!) - ID поста для отслеживания

**Пример подписки:**
```graphql
subscription {
  commentAdded(postId: "35d67a04-2829-4380-8f82-bcfdf8e5ca16") {
    id
    content
    parentId
    createdAt
  }
}
```

### WebSocket протокол

Для подписок используется WebSocket соединение:

1. **Подключение**: `ws://localhost:8080/subscriptions`
2. **Формат сообщения для подписки**:
```json
{
  "type": "subscribe",
  "query": "subscription { commentAdded(postId: \"...\") { id content createdAt } }",
  "variables": {
    "postId": "35d67a04-2829-4380-8f82-bcfdf8e5ca16"
  }
}
```

---

## 🔧 GraphQL Schema

### Типы данных

Система использует строго типизированную GraphQL схему:

#### **Post Type**
```graphql
type Post {
  id: ID!
  title: String!
  content: String!
  commentsEnabled: Boolean!
  createdAt: String!
  comments(limit: Int = 10, offset: Int = 0): [Comment!]!
}
```

#### **Comment Type**
```graphql
type Comment {
  id: ID!
  content: String!
  parentId: ID
  createdAt: String!
  children(limit: Int = 10, offset: Int = 0): [Comment!]!
}
```

#### **Operations**
- **Query**: Операции чтения (`posts`, `post`)
- **Mutation**: Операции изменения (`createPost`, `createComment`, `toggleComments`)
- **Subscription**: Real-time подписки (`commentAdded`)

### GraphQL Playground

#### **Доступ к Playground**
- **URL**: http://localhost:8080/
- **Функции**:
  - Интерактивное тестирование GraphQL
  - Автодополнение и валидация запросов
  - Просмотр схемы и документации
  - Тестирование подписок через WebSocket

### Health Check эндпоинт

#### **GET /health**
Проверка состояния сервиса и базы данных:

**Успешный ответ:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Ответ при ошибке:**
```json
{
  "status": "error",
  "error": "database connection failed",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 🗄️ База данных

### Схема базы данных

#### Таблица `posts`
| Поле | Тип | Описание | Ограничения |
|------|-----|----------|-------------|
| `id` | UUID | Первичный ключ | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `title` | VARCHAR(255) | Заголовок поста | NOT NULL |
| `content` | TEXT | Содержимое поста | NOT NULL |
| `comments_enabled` | BOOLEAN | Разрешены ли комментарии | NOT NULL, DEFAULT true |
| `created_at` | TIMESTAMPTZ | Время создания | NOT NULL, DEFAULT NOW() |

#### Таблица `comments`
| Поле | Тип | Описание | Ограничения |
|------|-----|----------|-------------|
| `id` | UUID | Первичный ключ | PRIMARY KEY, DEFAULT gen_random_uuid() |
| `post_id` | UUID | Ссылка на пост | FOREIGN KEY, NOT NULL |
| `parent_id` | UUID | Ссылка на родительский комментарий | FOREIGN KEY, NULL для корневых |
| `content` | VARCHAR(2000) | Текст комментария | NOT NULL, макс. 2000 символов |
| `created_at` | TIMESTAMPTZ | Время создания | NOT NULL, DEFAULT NOW() |

### Индексы
```sql
-- Оптимизация запросов комментариев по посту
CREATE INDEX idx_comments_post_id ON comments(post_id);

-- Оптимизация иерархических запросов
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
```

### Рекурсивные запросы

Для получения иерархии комментариев используется CTE (Common Table Expression):

```sql
WITH RECURSIVE comment_tree AS (
    -- Базовый случай: корневые комментарии
    SELECT id, post_id, parent_id, content, created_at, 0 as level
    FROM comments
    WHERE post_id = $1 AND parent_id IS NULL

    UNION ALL

    -- Рекурсивная часть: дочерние комментарии
    SELECT c.id, c.post_id, c.parent_id, c.content, c.created_at, ct.level + 1
    FROM comments c
    INNER JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT id, post_id, parent_id, content, created_at
FROM comment_tree
ORDER BY created_at;
```

---

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию | Обязательная |
|------------|----------|----------------------|--------------|
| `HTTP_ADDR` | Адрес HTTP сервера | `:8080` | Нет |
| `STORAGE_TYPE` | Тип хранилища (`postgres`, `memory`) | `postgres` | Нет |
| `DB_DSN` | Строка подключения к PostgreSQL | - | Да (только для postgres) |
| `LOG_LEVEL` | Уровень логирования | `info` | Нет |

### Типы хранилища

#### PostgreSQL (`STORAGE_TYPE=postgres`)
- **Назначение**: Production-ready хранилище с persistent данными
- **Файлы**: `internal/repository/postgres.go`, `internal/repository/postgres_test.go`
- **Особенности**:
  - Рекурсивные CTE запросы для иерархии комментариев
  - ACID транзакции и foreign key constraints
  - Connection pooling с pgxpool
  - Оптимизированные индексы для быстрого поиска

#### In-Memory (`STORAGE_TYPE=memory`)
- **Назначение**: Разработка, тестирование, демонстрация
- **Файлы**: `internal/repository/memory.go`, `internal/repository/memory_test.go`
- **Особенности**:
  - Thread-safe доступ через sync.RWMutex
  - Полная реализация интерфейса Storage
  - Рекурсивная иерархия комментариев в памяти
  - Каскадное удаление данных
  - Мгновенный запуск без внешних зависимостей

### Примеры конфигурации

#### Локальная разработка
```bash
export HTTP_ADDR=":8080"
export STORAGE_TYPE="postgres"
export DB_DSN="postgres://user:password@localhost:5433/postsdb?sslmode=disable"
export LOG_LEVEL="debug"
```

#### Docker Compose
```yaml
environment:
  - HTTP_ADDR=:8080
  - STORAGE_TYPE=postgres
  - DB_DSN=postgres://user:password@db:5432/postsdb?sslmode=disable
  - LOG_LEVEL=info
```

---

## 🚀 Развертывание

### Docker Compose (рекомендуется)

1. **Клонирование репозитория:**
```bash
git clone https://github.com/NarthurN/CommentsSystem.git
cd CommentsSystem
```

2. **Запуск:**
```bash
docker compose up -d
```

3. **Проверка работы:**
```bash
# Проверка статуса
docker compose ps

# Проверка логов
docker compose logs app

# Health check
curl http://localhost:8080/health
```

### Локальная разработка

1. **Установка зависимостей:**
```bash
go mod download
```

2. **Запуск PostgreSQL:**
```bash
docker compose up -d db
```

3. **Настройка переменных окружения:**
```bash
export DB_DSN="postgres://user:password@localhost:5433/postsdb?sslmode=disable"
export STORAGE_TYPE="postgres"
export HTTP_ADDR=":8080"
export LOG_LEVEL="debug"
```

4. **Запуск приложения:**
```bash
go run cmd/app/main.go
```

### Команды Make

```bash
# Сборка и запуск
make build              # Сборка приложения
make run               # Запуск приложения
make docker-build      # Сборка Docker образов
make docker-up         # Запуск контейнеров
make docker-down       # Остановка контейнеров

# Тестирование
make test              # Запуск всех тестов
make test-unit         # Только unit-тесты
make test-coverage     # Тесты с покрытием

# Разработка
make deps             # Установка зависимостей
make clean            # Очистка артефактов

# Помощь
make help             # Список всех команд
```

---

## 🧪 Тестирование

### Структура тестов

| Тип тестов | Расположение | Назначение |
|------------|--------------|------------|
| **Unit тесты** | `internal/service/service_test.go` | Тестирование бизнес-логики |
| **Integration тесты** | `internal/repository/postgres_test.go` | Тестирование с реальной БД |
| **Mock тесты** | Встроены в unit тесты | Изоляция внешних зависимостей |

### Запуск тестов

```bash
# Все тесты
make test

# Только unit-тесты (без БД)
make test-unit
# или
go test -short ./internal/...

# Интеграционные тесты
go test ./internal/repository/

# Тесты с покрытием
make test-coverage
# или
go test -cover ./internal/...

# Тесты определенного пакета
go test ./internal/service/
```

### Примеры тестов

#### Unit тест с мок-объектами
```go
func TestCreateComment_Validation(t *testing.T) {
    mockStorage := NewMockStorage()
    ps := pubsub.New()
    svc, err := service.New(mockStorage, ps)

    // Тест валидации длины комментария
    // Тест проверки разрешений на комментарии
    // Тест существования поста
}
```

#### Интеграционный тест
```go
func TestPostgresStorage_CommentHierarchy(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Тест с реальной PostgreSQL
    // Проверка иерархии комментариев
    // Проверка рекурсивных запросов
}
```

---

## 📊 Мониторинг и диагностика

### Health Check

Система предоставляет эндпоинт для проверки состояния:

```bash
# Проверка статуса
curl http://localhost:8080/health

# Ответ при нормальной работе
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Логирование

Система использует структурированное логирование с различными уровнями:

```go
// Пример логирования в коде
log.Info("Starting server", "addr", config.HTTPAddr)
log.Error("Database connection failed", "error", err)
```

### Graceful Shutdown

Приложение корректно завершает работу при получении сигналов:

```go
// Обработка сигналов завершения
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

<-c
log.Info("Shutting down server...")
```

---

## 🔧 Архитектурные решения

### Принципы Clean Architecture

1. **Dependency Inversion**: Сервисный слой определяет интерфейсы, которые реализует репозиторный слой
2. **Single Responsibility**: Каждый слой имеет одну ответственность
3. **Interface Segregation**: Мелкие, специфичные интерфейсы
4. **Open/Closed**: Легко расширяемая архитектура

### Конвертеры между слоями

#### API ↔ Domain конвертеры (`internal/converter/`)
```go
type GraphQLConverter struct{}
func (c *GraphQLConverter) PostToGraphQL(post *model.Post) map[string]interface{}
func (c *GraphQLConverter) PostsToGraphQL(posts []*model.Post) []interface{}

type ValidationConverter struct{}
func (c *ValidationConverter) ValidateAndConvertPostInput(title, content string) (*model.Post, error)
```

#### Domain ↔ Storage конвертеры (`internal/repository/converter/`)
```go
type PostConverter struct{}
func (c *PostConverter) ToStorage(post *model.Post) *storage_model.PostDB
func (c *PostConverter) FromStorage(postDB *storage_model.PostDB) *model.Post

type CommentConverter struct{}
func (c *CommentConverter) ToStorage(comment *model.Comment) *storage_model.CommentDB
func (c *CommentConverter) FromStorage(commentDB *storage_model.CommentDB) *model.Comment
```

### Модели по слоям

#### Domain модели (`internal/model/entities.go`)
```go
type Post struct {
    ID              uuid.UUID `json:"id"`
    Title           string    `json:"title"`
    Content         string    `json:"content"`
    CommentsEnabled bool      `json:"commentsEnabled"`
    CreatedAt       time.Time `json:"createdAt"`
}
```

#### Storage модели (`internal/repository/model/storage_models.go`)
```go
type PostDB struct {
    ID              string    `db:"id"`
    Title           string    `db:"title"`
    Content         string    `db:"content"`
    CommentsEnabled bool      `db:"comments_enabled"`
    CreatedAt       time.Time `db:"created_at"`
}
```

### Преимущества архитектуры

1. **Тестируемость**: Легкое создание мок-объектов для каждого слоя
2. **Изоляция**: Изменения в одном слое не влияют на другие
3. **Расширяемость**: Простое добавление новых storage провайдеров или API форматов
4. **Поддерживаемость**: Четкое разделение ответственности
5. **Переиспользование**: Конвертеры могут использоваться в разных контекстах

### Real-time возможности

#### WebSocket подписки
- Автоматическое уведомление клиентов о новых комментариях
- Эффективная система Pub/Sub для масштабирования
- Graceful handling отключений клиентов

#### Pub/Sub система
```go
// Публикация события о новом комментарии
topicName := fmt.Sprintf("post:%s:comments", postID.String())
s.pubsub.Publish(topicName, createdComment)

// Подписка на события
subscriber := s.pubsub.Subscribe(topicName, subscriberID)
```

---

## 🚀 Производительность и оптимизации

### Оптимизации базы данных
- Индексы на часто запрашиваемые поля
- Эффективные рекурсивные запросы для иерархии комментариев
- Пагинация для больших наборов данных

### Кэширование
- Connection pooling для PostgreSQL
- In-memory Pub/Sub для real-time уведомлений

### Масштабирование
- Stateless дизайн сервиса
- Готовность к горизонтальному масштабированию
- Изоляция бизнес-логики от инфраструктуры
