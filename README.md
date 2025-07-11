# Comments System

Система комментариев - это современный backend для управления постами и комментариями, построенный на Go с использованием GraphQL API.

## Особенности

- **GraphQL API** с полной поддержкой запросов, мутаций и подписок
- **WebSocket подписки** для получения комментариев в реальном времени
- **Иерархические комментарии** с поддержкой вложенных ответов
- **PostgreSQL** для надежного хранения данных
- **Clean Architecture** для maintainable кода
- **Docker** поддержка для простого развертывания
- **Graceful shutdown** с корректным завершением соединений
- **Comprehensive тестирование** с unit и integration тестами

## Технологический стек

- **Backend**: Go 1.21+
- **API**: GraphQL (graphql-go/graphql)
- **Database**: PostgreSQL 15+
- **WebSockets**: gorilla/websocket
- **HTTP Router**: chi/v5
- **Database Driver**: pgx/v5
- **Testing**: Go testing + testify
- **Containerization**: Docker & Docker Compose

## Структура проекта

```
.
├── api/                    # OpenAPI спецификация
├── cmd/app/               # Точка входа приложения
├── internal/
│   ├── api/               # HTTP обработчики (включая OpenAPI)
│   ├── config/            # Конфигурация
│   ├── model/             # Доменные модели (Domain Layer)
│   ├── converter/         # Конвертеры API ↔ Domain
│   ├── repository/        # Слой доступа к данным
│   │   ├── model/         # Модели данных для БД
│   │   └── converter/     # Конвертеры Domain ↔ Storage
│   └── service/           # Бизнес-логика и GraphQL резолверы
├── pkg/pubsub/            # In-memory Pub/Sub для подписок
├── migrations/            # SQL миграции
├── docker-compose.yml     # Конфигурация Docker
├── Dockerfile            # Образ приложения
├── CLEAN_ARCHITECTURE.md  # Детальное описание архитектуры
└── TechnicalDocument.md   # Технический документ
```

## Clean Architecture

Проект реализует принципы чистой архитектуры с четким разделением слоев:

### Слои архитектуры

1. **Domain Layer** (`internal/model/`): Чистые бизнес-модели без зависимостей
2. **Repository Layer** (`internal/repository/`): Доступ к данным с собственными моделями
3. **Service Layer** (`internal/service/`): Бизнес-логика и GraphQL резолверы
4. **API Layer** (`internal/api/`): HTTP обработчики и внешние интерфейсы
5. **Converter Layer** (`internal/converter/`): Изоляция между слоями

### Конвертеры

- **API ↔ Domain**: `internal/converter/api_converter.go`
- **Domain ↔ Storage**: `internal/repository/converter/converter.go`

### Модели

- **Domain модели**: `internal/model/entities.go`
- **Storage модели**: `internal/repository/model/storage_models.go`

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.22+ (для локальной разработки)
- Make (опционально)

### Запуск с Docker Compose

1. Клонируйте репозиторий:
```bash
git clone https://github.com/NarthurN/CommentsSystem.git
cd CommentsSystem
```

2. Запустите приложение:
```bash
docker compose up -d
```

3. Приложение будет доступно по адресу: http://localhost:8080

### Локальная разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите PostgreSQL:
```bash
docker-compose up -d db
```

3. Настройте переменные окружения:
```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Отредактируйте .env файл под ваше окружение
# Как минимум, установите DB_DSN
```

4. Запустите приложение:
```bash
go run cmd/app/main.go
```

## API Reference

### GraphQL Endpoint

- **URL**: `http://localhost:8080/graphql`
- **Method**: POST
- **Content-Type**: application/json

### GraphQL Playground

Доступен по адресу: `http://localhost:8080/`

### WebSocket Subscriptions

- **URL**: `ws://localhost:8080/subscriptions`
- **Protocol**: WebSocket

### Health Check

- **URL**: `http://localhost:8080/health`
- **Method**: GET

### Примеры GraphQL запросов

#### Получение списка постов

```graphql
query GetPosts {
  posts(limit: 10, offset: 0) {
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

#### Создание поста

```graphql
mutation CreatePost {
  createPost(title: "Новый пост", content: "Содержание поста") {
    id
    title
    content
    commentsEnabled
    createdAt
  }
}
```

#### Создание комментария

```graphql
mutation CreateComment {
  createComment(
    postId: "post-uuid-here"
    content: "Текст комментария"
    parentId: null
  ) {
    id
    content
    parentId
    createdAt
  }
}
```

#### Подписка на новые комментарии

```graphql
subscription CommentAdded {
  commentAdded(postId: "post-uuid-here") {
    id
    content
    parentId
    createdAt
  }
}
```

## Генерация кода

### OpenAPI генерация
```bash
# Генерация клиентского кода из OpenAPI спецификации
make generate-openapi
```

### Обновление зависимостей
```bash
# Установка зависимостей
make deps
```

## Тестирование

### Запуск всех тестов
```bash
make test
```

### Только unit-тесты
```bash
make test-unit
# или
go test -short ./internal/...
```

### Тесты с покрытием
```bash
make test-coverage
```

### Интеграционные тесты
```bash
# Требуют запущенной PostgreSQL
go test ./internal/repository/
```

## Команды Make

```bash
make build         # Сборка приложения
make run           # Запуск приложения
make test          # Запуск тестов
make test-unit     # Unit-тесты (без БД)
make test-coverage # Тесты с покрытием
make docker-build  # Сборка Docker образов
make docker-up     # Запуск контейнеров
make docker-down   # Остановка контейнеров
make generate-openapi # Генерация OpenAPI кода
make deps          # Установка зависимостей
make help          # Список всех команд
```

## Переменные окружения

### Основные настройки

| Переменная | Описание | Значение по умолчанию | Обязательная |
|------------|----------|----------------------|--------------|
| HTTP_ADDR | Адрес HTTP сервера | :8080 | Нет |
| STORAGE_TYPE | Тип хранилища | postgres | Нет |
| DB_DSN | Строка подключения к PostgreSQL | - | Да |
| LOG_LEVEL | Уровень логирования | info | Нет |

### HTTP сервер

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| HTTP_READ_TIMEOUT | Timeout чтения запросов | 15s |
| HTTP_WRITE_TIMEOUT | Timeout записи ответов | 15s |
| HTTP_IDLE_TIMEOUT | Timeout idle соединений | 60s |
| HTTP_SHUTDOWN_TIMEOUT | Timeout graceful shutdown | 30s |
| HTTP_REQUEST_TIMEOUT | Timeout обработки запросов | 60s |

### Бизнес-лимиты

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| POSTS_PAGE_LIMIT | Лимит постов на страницу | 10 |
| COMMENTS_PAGE_LIMIT | Лимит комментариев на страницу | 10 |
| MAX_TITLE_LENGTH | Максимальная длина заголовка | 255 |
| MAX_CONTENT_LENGTH | Максимальная длина контента | 10000 |
| MAX_COMMENT_LENGTH | Максимальная длина комментария | 2000 |

### PubSub

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| PUBSUB_CHANNEL_BUFFER_SIZE | Размер буфера каналов | 100 |
| PUBSUB_KEEP_ALIVE_PING | Интервал keep-alive пинга | 10s |

### CORS

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| CORS_ALLOW_ORIGIN | Разрешенные origins | * |
| CORS_ALLOW_METHODS | Разрешенные HTTP методы | GET, POST, OPTIONS |
| CORS_ALLOW_HEADERS | Разрешенные заголовки | Content-Type, Authorization |

### GraphQL

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| GRAPHQL_PLAYGROUND_TITLE | Заголовок Playground | GraphQL Playground |
| GRAPHQL_ENDPOINT | Путь к GraphQL API | /graphql |
| GRAPHQL_ENABLE_INTROSPECTION | Включить интроспекцию | true |

Подробное описание всех переменных можно найти в файле [.env.example](.env.example).

## Архитектурные особенности

### Принципы Clean Architecture

1. **Dependency Inversion**: Сервисный слой определяет интерфейсы
2. **Single Responsibility**: Каждый слой имеет одну ответственность
3. **Interface Segregation**: Мелкие, специфичные интерфейсы
4. **Open/Closed**: Легко расширяемая архитектура

### Изоляция слоев через конвертеры

#### API ↔ Domain (`internal/converter/api_converter.go`)
```go
type GraphQLConverter struct{}
func (c *GraphQLConverter) PostToGraphQL(post *model.Post) map[string]interface{}

type OpenAPIConverter struct{}
func (c *OpenAPIConverter) PostToOpenAPI(post *model.Post) *PostResponse

type ValidationConverter struct{}
func (c *ValidationConverter) ValidateAndConvertPostInput(title, content string) (*model.Post, error)
```

#### Domain ↔ Storage (`internal/repository/converter/converter.go`)
```go
type PostConverter struct{}
func (c *PostConverter) ToStorage(post *model.Post) *storage_model.PostDB
func (c *PostConverter) FromStorage(postDB *storage_model.PostDB) *model.Post
```

### OpenAPI интеграция

Система включает полную интеграцию с OpenAPI 3.0:

- **Автоматическая генерация типов** из OpenAPI спецификации
- **Swagger UI** для интерактивной документации
- **Валидация запросов** на основе OpenAPI схемы
- **Health check эндпоинты** для мониторинга

### Особенности реализации

- **Иерархические комментарии** с рекурсивными SQL запросами
- **Real-time подписки** через WebSocket
- **Пагинация** для постов и комментариев
- **Валидация данных** на всех уровнях:
  - Комментарии: макс. 2000 символов
  - Заголовки постов: 1-255 символов
  - Содержимое постов: 1-10000 символов
- **Возможность отключения комментариев** для поста
- **Graceful shutdown** с корректным завершением соединений
- **Пул соединений** к базе данных
- **Health check** эндпоинты для мониторинга

### Модели данных

#### Domain модели (`internal/model/entities.go`)
```go
type Post struct {
    ID              uuid.UUID `json:"id"`
    Title           string    `json:"title"`
    Content         string    `json:"content"`
    CommentsEnabled bool      `json:"commentsEnabled"`
    CreatedAt       time.Time `json:"createdAt"`
}

type Comment struct {
    ID        uuid.UUID  `json:"id"`
    PostID    uuid.UUID  `json:"postId"`
    ParentID  *uuid.UUID `json:"parentId,omitempty"`
    Content   string     `json:"content"`
    CreatedAt time.Time  `json:"createdAt"`
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

## Документация

- **[CLEAN_ARCHITECTURE.md](CLEAN_ARCHITECTURE.md)** - Детальное описание архитектуры
- **[TechnicalDocument.md](TechnicalDocument.md)** - Технический документ
- **[Swagger UI](http://localhost:8080/swagger/)** - Интерактивная API документация

## Преимущества архитектуры

1. **Тестируемость** - Легкое создание мок-объектов для каждого слоя
2. **Изоляция** - Изменения в одном слое не влияют на другие
3. **Расширяемость** - Простое добавление новых storage провайдеров или API форматов
4. **Поддерживаемость** - Четкое разделение ответственности
5. **Переиспользование** - Конвертеры могут использоваться в разных контекстах

## Диагностика

### Health Check
```bash
curl http://localhost:8080/api/health
```

### Логи
```bash
# Просмотр логов приложения
docker compose logs app

# Просмотр логов базы данных
docker compose logs db
```

### Проверка статуса
```bash
# Статус контейнеров
docker compose ps

# Проверка GraphQL
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "query { posts { id title } }"}'
```

## 🧪 Покрытие тестами

![Coverage](https://img.shields.io/badge/coverage-7.1%25-red)

**Общее покрытие:** 7.1%

📊 **Детализированное покрытие по модулям:**

```
github.com/NarthurN/CommentsSystem/cmd/app/main.go:28:					main											0.0%
github.com/NarthurN/CommentsSystem/cmd/app/main.go:102:					initializeStorage									0.0%
github.com/NarthurN/CommentsSystem/cmd/app/main.go:117:					waitForShutdownSignal									0.0%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:37:			NewGQLGenHandler									75.0%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:63:			NewGQLGenHandlerWithConfig								100.0%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:80:			SetupRoutes										100.0%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:105:			corsMiddleware										14.3%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:133:			isOriginAllowed										100.0%
github.com/NarthurN/CommentsSystem/internal/api/gqlgen_handler.go:153:			HandleHealthCheck									61.5%
github.com/NarthurN/CommentsSystem/internal/config/config.go:88:			LoadFromEnv										100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:136:			Validate										100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:165:			GetDSNForTests										100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:176:			getEnv											100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:184:			getIntEnv										100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:194:			getDurationEnv										100.0%
github.com/NarthurN/CommentsSystem/internal/config/config.go:204:			getBoolEnv										100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:38:		NewGraphQLConverter									100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:53:		PostToGraphQL										80.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:81:		CommentToGraphQL									77.8%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:117:		PostsToGraphQL										77.8%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:143:		CommentsToGraphQL									0.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:174:		NewValidationConverter									100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:190:		ValidateAndConvertCreatePost								100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:220:		ValidateAndConvertCreateComment								100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:265:		ValidatePaginationParams								100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:294:		validatePostInput									100.0%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:314:		validateCommentInput									85.7%
github.com/NarthurN/CommentsSystem/internal/converter/api_converter.go:330:		parseUUID										83.3%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:89:			IsValidTitle										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:95:			IsValidContent										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:100:			CanAddComments										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:106:			IsValid											100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:112:			Prepare											100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:125:			IsValidComment										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:130:			IsRootComment										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:135:			HasValidPost										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:140:			IsValid											100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:146:			Prepare											100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:158:			HasChildren										100.0%
github.com/NarthurN/CommentsSystem/internal/model/entities.go:163:			GetChildrenCount									100.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:16:	NewPostConverter									100.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:21:	ToRepositoryModel									66.7%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:36:	ToDomainModel										66.7%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:51:	ToDomainModels										83.3%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:65:	CreateNewPost										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:79:	NewCommentConverter									100.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:84:	ToRepositoryModel									66.7%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:99:	ToDomainModel										66.7%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:114:	ToDomainModels										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:128:	CreateNewComment									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:144:	NewTreeConverter									100.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:151:	BuildCommentTree									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:167:	buildTree										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/converter/converter.go:187:	ToPostWithComments									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:48:	TableName										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:53:	TableName										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:60:	GetSelectColumns									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:65:	GetSelectColumns									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:70:	GetInsertColumns									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:75:	GetInsertColumns									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:80:	GetUpdateColumns									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:87:	Validate										100.0%
github.com/NarthurN/CommentsSystem/internal/repository/model/storage_models.go:101:	Validate										100.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:25:			NewPostgresStorage									66.7%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:45:			Close											0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:51:			HealthCheck										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:58:			CreatePost										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:106:			GetPost											0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:130:			GetPosts										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:177:			UpdatePost										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:209:			DeletePost										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:225:			TogglePostComments									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:247:			CreateComment										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:295:			GetComment										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:319:			GetCommentsByPostID									0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:366:			GetCommentTree										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:418:			DeleteComment										0.0%
github.com/NarthurN/CommentsSystem/internal/repository/postgres.go:436:			GetPostWithComments									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:25:		NewExecutableSchema									100.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:119:		Schema											0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:126:		Complexity										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:294:		Exec											0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:377:		processDeferredGroup									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:396:		introspectSchema									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:403:		introspectType										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:451:		field_Comment_children_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:466:		field_Comment_children_argsLimit							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:484:		field_Comment_children_argsOffset							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:502:		field_Mutation_createComment_args							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:522:		field_Mutation_createComment_argsPostID							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:540:		field_Mutation_createComment_argsParentID						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:558:		field_Mutation_createComment_argsContent						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:576:		field_Mutation_createPost_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:591:		field_Mutation_createPost_argsTitle							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:609:		field_Mutation_createPost_argsContent							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:627:		field_Mutation_toggleComments_args							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:642:		field_Mutation_toggleComments_argsPostID						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:660:		field_Mutation_toggleComments_argsEnable						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:678:		field_Post_comments_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:693:		field_Post_comments_argsLimit								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:711:		field_Post_comments_argsOffset								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:729:		field_Query___type_args									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:739:		field_Query___type_argsName								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:757:		field_Query_post_args									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:767:		field_Query_post_argsID									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:785:		field_Query_posts_args									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:800:		field_Query_posts_argsLimit								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:818:		field_Query_posts_argsOffset								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:836:		field_Subscription_commentAdded_args							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:846:		field_Subscription_commentAdded_argsPostID						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:864:		field___Directive_args_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:874:		field___Directive_args_argsIncludeDeprecated						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:892:		field___Field_args_args									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:902:		field___Field_args_argsIncludeDeprecated						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:920:		field___Type_enumValues_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:930:		field___Type_enumValues_argsIncludeDeprecated						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:948:		field___Type_fields_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:958:		field___Type_fields_argsIncludeDeprecated						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:984:		_Comment_id										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1015:		fieldContext_Comment_id									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1028:		_Comment_content									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1059:		fieldContext_Comment_content								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1072:		_Comment_parentId									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1100:		fieldContext_Comment_parentId								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1113:		_Comment_createdAt									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1144:		fieldContext_Comment_createdAt								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1157:		_Comment_children									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1188:		fieldContext_Comment_children								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1224:		_Mutation_createPost									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1255:		fieldContext_Mutation_createPost							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1293:		_Mutation_createComment									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1324:		fieldContext_Mutation_createComment							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1360:		_Mutation_toggleComments								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1391:		fieldContext_Mutation_toggleComments							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1429:		_Post_id										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1460:		fieldContext_Post_id									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1473:		_Post_title										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1504:		fieldContext_Post_title									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1517:		_Post_content										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1548:		fieldContext_Post_content								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1561:		_Post_commentsEnabled									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1592:		fieldContext_Post_commentsEnabled							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1605:		_Post_createdAt										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1636:		fieldContext_Post_createdAt								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1649:		_Post_comments										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1680:		fieldContext_Post_comments								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1716:		_Query_posts										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1747:		fieldContext_Query_posts								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1785:		_Query_post										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1813:		fieldContext_Query_post									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1851:		_Query___type										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1879:		fieldContext_Query___type								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1927:		_Query___schema										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1955:		fieldContext_Query___schema								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:1982:		_Subscription_commentAdded								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2027:		fieldContext_Subscription_commentAdded							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2063:		___Directive_name									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2094:		fieldContext___Directive_name								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2107:		___Directive_description								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2135:		fieldContext___Directive_description							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2148:		___Directive_isRepeatable								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2179:		fieldContext___Directive_isRepeatable							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2192:		___Directive_locations									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2223:		fieldContext___Directive_locations							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2236:		___Directive_args									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2267:		fieldContext___Directive_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2305:		___EnumValue_name									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2336:		fieldContext___EnumValue_name								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2349:		___EnumValue_description								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2377:		fieldContext___EnumValue_description							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2390:		___EnumValue_isDeprecated								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2421:		fieldContext___EnumValue_isDeprecated							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2434:		___EnumValue_deprecationReason								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2462:		fieldContext___EnumValue_deprecationReason						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2475:		___Field_name										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2506:		fieldContext___Field_name								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2519:		___Field_description									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2547:		fieldContext___Field_description							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2560:		___Field_args										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2591:		fieldContext___Field_args								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2629:		___Field_type										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2660:		fieldContext___Field_type								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2697:		___Field_isDeprecated									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2728:		fieldContext___Field_isDeprecated							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2741:		___Field_deprecationReason								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2769:		fieldContext___Field_deprecationReason							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2782:		___InputValue_name									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2813:		fieldContext___InputValue_name								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2826:		___InputValue_description								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2854:		fieldContext___InputValue_description							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2867:		___InputValue_type									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2898:		fieldContext___InputValue_type								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2935:		___InputValue_defaultValue								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2963:		fieldContext___InputValue_defaultValue							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:2976:		___InputValue_isDeprecated								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3007:		fieldContext___InputValue_isDeprecated							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3020:		___InputValue_deprecationReason								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3048:		fieldContext___InputValue_deprecationReason						0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3061:		___Schema_description									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3089:		fieldContext___Schema_description							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3102:		___Schema_types										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3133:		fieldContext___Schema_types								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3170:		___Schema_queryType									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3201:		fieldContext___Schema_queryType								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3238:		___Schema_mutationType									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3266:		fieldContext___Schema_mutationType							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3303:		___Schema_subscriptionType								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3331:		fieldContext___Schema_subscriptionType							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3368:		___Schema_directives									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3399:		fieldContext___Schema_directives							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3424:		___Type_kind										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3455:		fieldContext___Type_kind								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3468:		___Type_name										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3496:		fieldContext___Type_name								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3509:		___Type_description									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3537:		fieldContext___Type_description								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3550:		___Type_specifiedByURL									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3578:		fieldContext___Type_specifiedByURL							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3591:		___Type_fields										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3619:		fieldContext___Type_fields								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3657:		___Type_interfaces									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3685:		fieldContext___Type_interfaces								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3722:		___Type_possibleTypes									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3750:		fieldContext___Type_possibleTypes							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3787:		___Type_enumValues									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3815:		fieldContext___Type_enumValues								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3849:		___Type_inputFields									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3877:		fieldContext___Type_inputFields								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3904:		___Type_ofType										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3932:		fieldContext___Type_ofType								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3969:		___Type_isOneOf										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:3997:		fieldContext___Type_isOneOf								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4024:		_Comment										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4204:		_Mutation										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4267:		_Post											0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4424:		_Query											0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4515:		_Subscription										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4535:		___Directive										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4591:		___EnumValue										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4639:		___Field										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4697:		___InputValue										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4752:		___Schema										0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4807:		___Type											0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4868:		unmarshalNBoolean2bool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4873:		marshalNBoolean2bool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4884:		marshalNComment2githubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐComment		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4888:		marshalNComment2ᚕᚖgithubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐCommentᚄ		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4932:		marshalNComment2ᚖgithubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐComment		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4942:		unmarshalNID2string									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4947:		marshalNID2string									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4958:		marshalNPost2githubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐPost			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:4962:		marshalNPost2ᚕᚖgithubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐPostᚄ			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5006:		marshalNPost2ᚖgithubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐPost			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5016:		unmarshalNString2string									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5021:		marshalNString2string									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5032:		marshalN__Directive2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirective		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5036:		marshalN__Directive2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirectiveᚄ	0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5080:		unmarshalN__DirectiveLocation2string							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5085:		marshalN__DirectiveLocation2string							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5096:		unmarshalN__DirectiveLocation2ᚕstringᚄ							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5111:		marshalN__DirectiveLocation2ᚕstringᚄ							0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5155:		marshalN__EnumValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValue		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5159:		marshalN__Field2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐField			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5163:		marshalN__InputValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValue	0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5167:		marshalN__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ	0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5211:		marshalN__Type2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5215:		marshalN__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5259:		marshalN__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5269:		unmarshalN__TypeKind2string								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5274:		marshalN__TypeKind2string								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5285:		unmarshalOBoolean2bool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5290:		marshalOBoolean2bool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5297:		unmarshalOBoolean2ᚖbool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5305:		marshalOBoolean2ᚖbool									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5315:		unmarshalOID2ᚖstring									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5323:		marshalOID2ᚖstring									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5333:		unmarshalOInt2ᚖint									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5341:		marshalOInt2ᚖint									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5351:		marshalOPost2ᚖgithubᚗcomᚋNarthurNᚋCommentsSystemᚋinternalᚋmodelᚐPost			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5358:		unmarshalOString2ᚖstring								0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5366:		marshalOString2ᚖstring									0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5376:		marshalO__EnumValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValueᚄ	0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5423:		marshalO__Field2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐFieldᚄ		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5470:		marshalO__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ	0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5517:		marshalO__Schema2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐSchema		0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5524:		marshalO__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ			0.0%
github.com/NarthurN/CommentsSystem/internal/service/generated/exec.go:5571:		marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType			0.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:36:		NewGQLGenService									100.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:62:		NewGQLGenServiceWithConfig								88.9%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:104:		GetHandler										100.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:110:		GetPlaygroundHandler									100.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:116:		HealthCheck										66.7%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:132:		GetConfig										100.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:138:		GetSubscribersCount									100.0%
github.com/NarthurN/CommentsSystem/internal/service/gqlgen_service.go:146:		checkOriginAllowed									0.0%
github.com/NarthurN/CommentsSystem/internal/service/resolver.go:18:			NewResolver										100.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:18:		ID											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:23:		ParentID										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:32:		CreatedAt										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:37:		Children										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:43:		CreatePost										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:58:		CreateComment										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:112:		ToggleComments										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:136:		ID											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:141:		CreatedAt										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:146:		Comments										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:162:		Posts											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:182:		Post											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:197:		CommentAdded										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:250:		Comment											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:253:		Mutation										0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:256:		Post											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:259:		Query											0.0%
github.com/NarthurN/CommentsSystem/internal/service/schema.resolvers.go:262:		Subscription										0.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:42:				New											100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:49:				NewWithConfig										100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:69:				Subscribe										100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:94:				Unsubscribe										100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:124:				Publish											100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:158:				GetSubscribersCount									100.0%
github.com/NarthurN/CommentsSystem/pkg/pubsub/pubsub.go:171:				Close											100.0%
```

*Отчет автоматически обновлен 2025-07-11 18:25:11*
