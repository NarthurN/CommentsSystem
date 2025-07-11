package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NarthurN/CommentsSystem/internal/api"
	"github.com/NarthurN/CommentsSystem/internal/config"
	"github.com/NarthurN/CommentsSystem/internal/repository"
	"github.com/NarthurN/CommentsSystem/internal/service"
	"github.com/NarthurN/CommentsSystem/pkg/pubsub"
	"github.com/joho/godotenv"
)

// Константы приложения
const (
	// ExitCodeSuccess - код успешного завершения программы
	ExitCodeSuccess = 0
	// ExitCodeError - код завершения программы с ошибкой
	ExitCodeError = 1
)

// main - точка входа в приложение CommentsSystem.
// Инициализирует все компоненты, запускает HTTP сервер и обрабатывает graceful shutdown.
func main() {
	// Загружаем переменные окружения из .env файла (если он существует)
	// Игнорируем ошибку, так как .env файл опционален
	_ = godotenv.Load()

	// Загружаем конфигурацию приложения из переменных окружения
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Создаем контекст для координации graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем слой хранения данных на основе конфигурации
	storage, err := initializeStorage(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer func() {
		if closeErr := storage.Close(); closeErr != nil {
			log.Printf("Error closing storage: %v", closeErr)
		}
	}()

	// Инициализируем pub/sub систему для real-time подписок
	ps := pubsub.NewWithConfig(cfg.ChannelBufferSize)

	// Создаем GraphQL сервис с использованием gqlgen
	gqlgenService := service.NewGQLGenServiceWithConfig(storage, ps, cfg)

	// Создаем HTTP обработчик с конфигурацией
	handler := api.NewGQLGenHandlerWithConfig(gqlgenService, cfg)
	router := handler.SetupRoutes()

	// Создаем и настраиваем HTTP сервер
	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("Starting %s server on %s", config.AppName, cfg.HTTPAddr)
		log.Printf("GraphQL Playground available at http://localhost%s/", cfg.HTTPAddr)
		log.Printf("GraphQL endpoint at http://localhost%s%s", cfg.HTTPAddr, cfg.GraphQLEndpoint)
		log.Printf("Health check endpoint at http://localhost%s/health", cfg.HTTPAddr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Ожидаем сигнал прерывания для graceful shutdown
	waitForShutdownSignal()

	log.Println("Shutting down server...")

	// Создаем контекст для shutdown с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	// Пытаемся выполнить graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		os.Exit(ExitCodeError)
	}

	log.Println("Server stopped gracefully")
	os.Exit(ExitCodeSuccess)
}

// initializeStorage создает и инициализирует слой хранения данных на основе конфигурации.
// Поддерживает различные типы хранилищ и возвращает правильно настроенный интерфейс Storage.
func initializeStorage(ctx context.Context, cfg *config.Config) (repository.Storage, error) {
	log.Printf("Initializing storage type: %s", cfg.StorageType)

	switch cfg.StorageType {
	case "postgres":
		log.Printf("Connecting to PostgreSQL database...")
		storage, err := repository.NewPostgresStorage(ctx, cfg.DatabaseDSN)
		if err != nil {
			return nil, err
		}
		log.Printf("PostgreSQL storage initialized successfully")
		return storage, nil
	case "memory":
		log.Printf("Initializing in-memory storage...")
		storage := repository.NewMemoryStorage()
		log.Printf("In-memory storage initialized successfully")
		return storage, nil
	default:
		return nil, fmt.Errorf("%w: supported types are 'postgres' and 'memory', got '%s'",
			repository.ErrUnsupportedStorageType, cfg.StorageType)
	}
}

// waitForShutdownSignal блокирует выполнение до получения сигнала прерывания.
// Прослушивает сигналы SIGINT (Ctrl+C) и SIGTERM для graceful shutdown.
func waitForShutdownSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
