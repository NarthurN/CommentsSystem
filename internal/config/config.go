package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AppName содержит имя приложения
const AppName = "CommentsSystem"

// Значения конфигурации по умолчанию
const (
	// Настройки сервера по умолчанию
	DefaultHTTPAddr        = ":8080"
	DefaultReadTimeout     = 15 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 30 * time.Second
	DefaultRequestTimeout  = 60 * time.Second

	// Настройки базы данных по умолчанию
	DefaultStorageType = "postgres"
	DefaultLogLevel    = "info"

	// Настройки бизнес-логики по умолчанию
	DefaultPostsPageLimit    = 10
	DefaultCommentsPageLimit = 10
	DefaultMaxTitleLength    = 255
	DefaultMaxContentLength  = 10000
	DefaultMaxCommentLength  = 2000

	// Настройки PubSub по умолчанию
	DefaultChannelBufferSize = 100
	DefaultKeepAlivePing     = 10 * time.Second

	// Настройки CORS по умолчанию
	DefaultAllowOrigin  = "*"
	DefaultAllowMethods = "GET, POST, OPTIONS"
	DefaultAllowHeaders = "Content-Type, Authorization"

	// Настройки GraphQL по умолчанию
	DefaultPlaygroundTitle = "GraphQL Playground"
	DefaultGraphQLEndpoint = "/graphql"
)

// Config представляет конфигурацию приложения
type Config struct {
	// Конфигурация HTTP сервера
	HTTPAddr        string        `json:"http_addr"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	RequestTimeout  time.Duration `json:"request_timeout"`

	// Конфигурация базы данных
	StorageType string `json:"storage_type"`
	DatabaseDSN string `json:"database_dsn"`

	// Конфигурация логирования
	LogLevel string `json:"log_level"`

	// Ограничения бизнес-логики
	PostsPageLimit    int `json:"posts_page_limit"`
	CommentsPageLimit int `json:"comments_page_limit"`
	MaxTitleLength    int `json:"max_title_length"`
	MaxContentLength  int `json:"max_content_length"`
	MaxCommentLength  int `json:"max_comment_length"`

	// Конфигурация PubSub
	ChannelBufferSize int           `json:"channel_buffer_size"`
	KeepAlivePing     time.Duration `json:"keep_alive_ping"`

	// Конфигурация CORS
	AllowOrigin  string `json:"allow_origin"`
	AllowMethods string `json:"allow_methods"`
	AllowHeaders string `json:"allow_headers"`

	// Конфигурация GraphQL
	PlaygroundTitle     string `json:"playground_title"`
	GraphQLEndpoint     string `json:"graphql_endpoint"`
	EnableIntrospection bool   `json:"enable_introspection"`
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		// HTTP сервер
		HTTPAddr:        getEnv("HTTP_ADDR", DefaultHTTPAddr),
		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", DefaultReadTimeout),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", DefaultWriteTimeout),
		IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT", DefaultIdleTimeout),
		ShutdownTimeout: getDurationEnv("HTTP_SHUTDOWN_TIMEOUT", DefaultShutdownTimeout),
		RequestTimeout:  getDurationEnv("HTTP_REQUEST_TIMEOUT", DefaultRequestTimeout),

		// База данных
		StorageType: getEnv("STORAGE_TYPE", DefaultStorageType),
		DatabaseDSN: getEnv("DB_DSN", ""),

		// Логирование
		LogLevel: getEnv("LOG_LEVEL", DefaultLogLevel),

		// Ограничения бизнес-логики
		PostsPageLimit:    getIntEnv("POSTS_PAGE_LIMIT", DefaultPostsPageLimit),
		CommentsPageLimit: getIntEnv("COMMENTS_PAGE_LIMIT", DefaultCommentsPageLimit),
		MaxTitleLength:    getIntEnv("MAX_TITLE_LENGTH", DefaultMaxTitleLength),
		MaxContentLength:  getIntEnv("MAX_CONTENT_LENGTH", DefaultMaxContentLength),
		MaxCommentLength:  getIntEnv("MAX_COMMENT_LENGTH", DefaultMaxCommentLength),

		// PubSub
		ChannelBufferSize: getIntEnv("PUBSUB_CHANNEL_BUFFER_SIZE", DefaultChannelBufferSize),
		KeepAlivePing:     getDurationEnv("PUBSUB_KEEP_ALIVE_PING", DefaultKeepAlivePing),

		// CORS
		AllowOrigin:  getEnv("CORS_ALLOW_ORIGIN", DefaultAllowOrigin),
		AllowMethods: getEnv("CORS_ALLOW_METHODS", DefaultAllowMethods),
		AllowHeaders: getEnv("CORS_ALLOW_HEADERS", DefaultAllowHeaders),

		// GraphQL
		PlaygroundTitle:     getEnv("GRAPHQL_PLAYGROUND_TITLE", DefaultPlaygroundTitle),
		GraphQLEndpoint:     getEnv("GRAPHQL_ENDPOINT", DefaultGraphQLEndpoint),
		EnableIntrospection: getBoolEnv("GRAPHQL_ENABLE_INTROSPECTION", true),
	}

	// Валидируем конфигурацию
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.StorageType == "postgres" && c.DatabaseDSN == "" {
		return fmt.Errorf("DB_DSN is required when STORAGE_TYPE is postgres")
	}

	if c.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR cannot be empty")
	}

	if c.ReadTimeout <= 0 || c.WriteTimeout <= 0 || c.IdleTimeout <= 0 {
		return fmt.Errorf("all timeout values must be positive")
	}

	if c.PostsPageLimit <= 0 || c.CommentsPageLimit <= 0 {
		return fmt.Errorf("page limits must be positive")
	}

	if c.MaxTitleLength <= 0 || c.MaxContentLength <= 0 || c.MaxCommentLength <= 0 {
		return fmt.Errorf("content length limits must be positive")
	}

	if c.ChannelBufferSize <= 0 {
		return fmt.Errorf("channel buffer size must be positive")
	}

	return nil
}

// GetDSNForTests возвращает DSN для тестов (может быть переопределено)
func (c *Config) GetDSNForTests() string {
	testDSN := getEnv("TEST_DB_DSN", c.DatabaseDSN)
	if testDSN == "" {
		return "postgres://user:password@localhost:5433/postsdb_test?sslmode=disable"
	}
	return testDSN
}

// Вспомогательные функции для парсинга переменных окружения

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv возвращает целочисленное значение переменной окружения или значение по умолчанию
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv возвращает значение времени из переменной окружения или значение по умолчанию
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getBoolEnv возвращает булево значение переменной окружения или значение по умолчанию
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
