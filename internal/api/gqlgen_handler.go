package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Общие константы HTTP статусов
const (
	StatusOK    = "ok"
	StatusError = "error"
)

// GQLGenHandler представляет HTTP обработчики с gqlgen.
// Управляет маршрутизацией, middleware и обработкой HTTP запросов.
//
// Основные возможности:
// - GraphQL API endpoint
// - GraphQL Playground интерфейс
// - Health check мониторинг
// - Настраиваемые CORS политики
// - Управление timeout запросов
type GQLGenHandler struct {
	service *service.GQLGenService // GraphQL сервис
	config  *config.Config         // Конфигурация приложения
}

// NewGQLGenHandler создает новый экземпляр GQLGenHandler с конфигурацией по умолчанию.
// Использует стандартные настройки для CORS и timeout'ов.
func NewGQLGenHandler(svc *service.GQLGenService) *GQLGenHandler {
	// Получаем конфигурацию из сервиса или создаем базовую
	cfg := svc.GetConfig()
	if cfg == nil {
		cfg = &config.Config{
			RequestTimeout: config.DefaultRequestTimeout,
			AllowOrigin:    config.DefaultAllowOrigin,
			AllowMethods:   config.DefaultAllowMethods,
			AllowHeaders:   config.DefaultAllowHeaders,
		}
	}

	return NewGQLGenHandlerWithConfig(svc, cfg)
}

// NewGQLGenHandlerWithConfig создает новый экземпляр GQLGenHandler с пользовательской конфигурацией.
// Позволяет полностью настроить поведение HTTP обработчика.
//
// Параметры:
//   - svc: GraphQL сервис
//   - cfg: конфигурация приложения
//
// Настраивает:
//   - Request timeout на основе конфигурации
//   - CORS политики из конфигурации
//   - Маршруты GraphQL и Playground
func NewGQLGenHandlerWithConfig(svc *service.GQLGenService, cfg *config.Config) *GQLGenHandler {
	return &GQLGenHandler{
		service: svc,
		config:  cfg,
	}
}

// SetupRoutes настраивает маршруты и middleware для gqlgen.
// Создает полный HTTP роутер с необходимыми обработчиками.
//
// Настраивает:
//   - Логирование запросов
//   - Восстановление после паник
//   - Timeout для запросов
//   - CORS политики
//   - GraphQL endpoints
//   - Health check endpoint
func (h *GQLGenHandler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Основные middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(h.config.RequestTimeout))

	// CORS middleware с настраиваемыми политиками
	r.Use(h.corsMiddleware())

	// GraphQL эндпоинт
	r.Handle(h.config.GraphQLEndpoint, h.service.GetHandler())

	// GraphQL Playground (обычно на корневом пути)
	r.Handle("/", h.service.GetPlaygroundHandler())

	// Health check endpoint
	r.Get("/health", h.HandleHealthCheck)

	return r
}

// corsMiddleware создает middleware для обработки CORS запросов.
// Использует конфигурацию для определения разрешенных origins, методов и заголовков.
func (h *GQLGenHandler) corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Устанавливаем CORS заголовки
			origin := r.Header.Get("Origin")
			if h.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if h.config.AllowOrigin == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", h.config.AllowMethods)
			w.Header().Set("Access-Control-Allow-Headers", h.config.AllowHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Обрабатываем preflight запросы
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed проверяет, разрешен ли указанный origin.
// Поддерживает как одиночные origins, так и списки через запятую.
func (h *GQLGenHandler) isOriginAllowed(origin string) bool {
	if h.config.AllowOrigin == "*" {
		return true
	}

	allowedOrigins := strings.Split(h.config.AllowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}

	return false
}

// HandleHealthCheck обрабатывает проверку состояния сервиса.
// Возвращает детальную информацию о состоянии всех компонентов.
//
// HTTP 200: сервис работает нормально
// HTTP 500: обнаружены проблемы
func (h *GQLGenHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Проверяем состояние сервиса
	err := h.service.HealthCheck(ctx)

	response := map[string]interface{}{
		"status":    StatusOK,
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   config.AppName,
		"version":   "1.0.0", // TODO: получать из build-time переменных
	}

	if err != nil {
		response["status"] = StatusError
		response["error"] = err.Error()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Добавляем дополнительную информацию при успешной проверке
		response["subscribers_count"] = h.service.GetSubscribersCount()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
