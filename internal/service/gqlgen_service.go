package service

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/internal/service/generated"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
	"github.com/gorilla/websocket"
)

// GQLGenService представляет сервис с использованием gqlgen.
// Инкапсулирует GraphQL сервер, резолверы и конфигурацию.
//
// Основные возможности:
// - GraphQL API с поддержкой запросов, мутаций и подписок
// - WebSocket транспорт для real-time подписок
// - Настраиваемые CORS политики
// - Health check для мониторинга
type GQLGenService struct {
	storage  repository.Storage // Интерфейс для работы с данными
	pubsub   *pubsub.PubSub     // Система pub/sub для подписок
	resolver *Resolver          // GraphQL резолверы
	server   *handler.Server    // GraphQL сервер
	config   *config.Config     // Конфигурация приложения
}

// NewGQLGenService создает новый экземпляр сервиса с gqlgen и конфигурацией по умолчанию.
// Использует стандартные настройки для WebSocket и CORS.
func NewGQLGenService(storage repository.Storage, ps *pubsub.PubSub) *GQLGenService {
	// Создаем временную конфигурацию для обратной совместимости
	cfg := &config.Config{
		KeepAlivePing:       config.DefaultKeepAlivePing,
		AllowOrigin:         config.DefaultAllowOrigin,
		PlaygroundTitle:     config.DefaultPlaygroundTitle,
		GraphQLEndpoint:     config.DefaultGraphQLEndpoint,
		EnableIntrospection: true,
	}

	return NewGQLGenServiceWithConfig(storage, ps, cfg)
}

// NewGQLGenServiceWithConfig создает новый экземпляр сервиса с gqlgen и пользовательской конфигурацией.
// Позволяет полностью настроить поведение GraphQL сервера.
//
// Параметры:
//   - storage: интерфейс для работы с данными
//   - ps: система pub/sub для real-time подписок
//   - cfg: конфигурация приложения
//
// Настраивает:
//   - HTTP POST/GET транспорты
//   - WebSocket транспорт с настраиваемыми параметрами
//   - CORS политики на основе конфигурации
//   - GraphQL интроспекцию (опционально)
func NewGQLGenServiceWithConfig(storage repository.Storage, ps *pubsub.PubSub, cfg *config.Config) *GQLGenService {
	resolver := NewResolver(storage, ps)

	// Создаем GraphQL сервер с сгенерированной схемой
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Настраиваем HTTP транспорты
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})

	// Настраиваем WebSocket транспорт для подписок
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// В продакшене следует реализовать более строгую проверку
				// основанную на cfg.AllowOrigin
				return cfg.AllowOrigin == "*" || checkOriginAllowed(r, cfg.AllowOrigin)
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: cfg.KeepAlivePing,
	})

	// Добавляем расширения на основе конфигурации
	if cfg.EnableIntrospection {
		srv.Use(extension.Introspection{})
	}

	return &GQLGenService{
		storage:  storage,
		pubsub:   ps,
		resolver: resolver,
		server:   srv,
		config:   cfg,
	}
}

// GetHandler возвращает HTTP обработчик для GraphQL эндпоинта.
// Используется для регистрации маршрута в HTTP роутере.
func (s *GQLGenService) GetHandler() http.Handler {
	return s.server
}

// GetPlaygroundHandler возвращает обработчик для GraphQL Playground.
// Предоставляет интерактивный интерфейс для тестирования GraphQL запросов.
func (s *GQLGenService) GetPlaygroundHandler() http.Handler {
	return playground.Handler(s.config.PlaygroundTitle, s.config.GraphQLEndpoint)
}

// HealthCheck проверяет состояние сервиса и его зависимостей.
// Возвращает ошибку если какой-либо компонент недоступен.
func (s *GQLGenService) HealthCheck(ctx context.Context) error {
	// Проверяем состояние хранилища данных
	if err := s.storage.HealthCheck(ctx); err != nil {
		return err
	}

	// Дополнительные проверки можно добавить здесь:
	// - Проверка pub/sub системы
	// - Проверка внешних API
	// - Проверка метрик производительности

	return nil
}

// GetConfig возвращает текущую конфигурацию сервиса.
// Используется API обработчиком для получения настроек.
func (s *GQLGenService) GetConfig() *config.Config {
	return s.config
}

// GetSubscribersCount возвращает количество активных подписчиков.
// Используется для мониторинга нагрузки real-time функций.
func (s *GQLGenService) GetSubscribersCount() int {
	// Поскольку удалили GetTotalSubscribers из PubSub, используем заглушку
	// В реальной реализации можно подсчитывать подписчиков другим способом
	return 0
}

// checkOriginAllowed проверяет, разрешен ли origin для WebSocket соединений.
// В продакшене следует реализовать более сложную логику проверки.
func checkOriginAllowed(r *http.Request, allowedOrigin string) bool {
	origin := r.Header.Get("Origin")
	return origin == allowedOrigin
}
