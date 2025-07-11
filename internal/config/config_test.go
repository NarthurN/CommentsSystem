package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadFromEnv(t *testing.T) {
	// Сохраняем оригинальные переменные окружения для восстановления
	originalEnv := make(map[string]string)
	envVars := []string{
		"HTTP_ADDR", "DB_DSN", "STORAGE_TYPE", "LOG_LEVEL",
		"HTTP_READ_TIMEOUT", "HTTP_WRITE_TIMEOUT", "HTTP_IDLE_TIMEOUT",
		"POSTS_PAGE_LIMIT", "MAX_TITLE_LENGTH", "CORS_ALLOW_ORIGIN",
	}

	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val != "" {
			originalEnv[envVar] = val
		}
	}

	// Очищаем переменные окружения после теста
	defer func() {
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}
		for envVar, val := range originalEnv {
			os.Setenv(envVar, val)
		}
	}()

	t.Run("загрузка с минимальными настройками", func(t *testing.T) {
		// Очищаем все переменные
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		// Устанавливаем только обязательные переменные
		os.Setenv("DB_DSN", "postgres://test:test@localhost/testdb")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("LoadFromEnv() error = %v", err)
		}

		// Проверяем значения по умолчанию
		if cfg.HTTPAddr != DefaultHTTPAddr {
			t.Errorf("HTTPAddr = %v, expected %v", cfg.HTTPAddr, DefaultHTTPAddr)
		}
		if cfg.StorageType != DefaultStorageType {
			t.Errorf("StorageType = %v, expected %v", cfg.StorageType, DefaultStorageType)
		}
		if cfg.PostsPageLimit != DefaultPostsPageLimit {
			t.Errorf("PostsPageLimit = %v, expected %v", cfg.PostsPageLimit, DefaultPostsPageLimit)
		}
	})

	t.Run("загрузка с пользовательскими настройками", func(t *testing.T) {
		// Очищаем все переменные
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		// Устанавливаем пользовательские значения
		os.Setenv("HTTP_ADDR", ":9000")
		os.Setenv("DB_DSN", "postgres://custom:pass@localhost/customdb")
		os.Setenv("STORAGE_TYPE", "postgres")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("HTTP_READ_TIMEOUT", "30s")
		os.Setenv("POSTS_PAGE_LIMIT", "20")
		os.Setenv("MAX_TITLE_LENGTH", "300")
		os.Setenv("CORS_ALLOW_ORIGIN", "https://example.com")

		cfg, err := LoadFromEnv()
		if err != nil {
			t.Fatalf("LoadFromEnv() error = %v", err)
		}

		// Проверяем пользовательские значения
		if cfg.HTTPAddr != ":9000" {
			t.Errorf("HTTPAddr = %v, expected :9000", cfg.HTTPAddr)
		}
		if cfg.DatabaseDSN != "postgres://custom:pass@localhost/customdb" {
			t.Errorf("DatabaseDSN = %v, expected custom DSN", cfg.DatabaseDSN)
		}
		if cfg.LogLevel != "debug" {
			t.Errorf("LogLevel = %v, expected debug", cfg.LogLevel)
		}
		if cfg.ReadTimeout != 30*time.Second {
			t.Errorf("ReadTimeout = %v, expected 30s", cfg.ReadTimeout)
		}
		if cfg.PostsPageLimit != 20 {
			t.Errorf("PostsPageLimit = %v, expected 20", cfg.PostsPageLimit)
		}
		if cfg.MaxTitleLength != 300 {
			t.Errorf("MaxTitleLength = %v, expected 300", cfg.MaxTitleLength)
		}
		if cfg.AllowOrigin != "https://example.com" {
			t.Errorf("AllowOrigin = %v, expected https://example.com", cfg.AllowOrigin)
		}
	})

	t.Run("ошибка при отсутствии DSN для postgres", func(t *testing.T) {
		// Очищаем все переменные
		for _, envVar := range envVars {
			os.Unsetenv(envVar)
		}

		os.Setenv("STORAGE_TYPE", "postgres")
		// Не устанавливаем DB_DSN

		_, err := LoadFromEnv()
		if err == nil {
			t.Fatal("Expected error when DB_DSN is missing for postgres storage")
		}
	})
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "валидная конфигурация",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: false,
		},
		{
			name: "отсутствует DSN для postgres",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: true,
		},
		{
			name: "пустой HTTP адрес",
			config: &Config{
				HTTPAddr:          "",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: true,
		},
		{
			name: "отрицательные таймауты",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       -1 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: true,
		},
		{
			name: "отрицательные лимиты страниц",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    -1,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: true,
		},
		{
			name: "отрицательные лимиты контента",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    -1,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: 100,
			},
			wantErr: true,
		},
		{
			name: "отрицательный размер буфера канала",
			config: &Config{
				HTTPAddr:          ":8080",
				StorageType:       "postgres",
				DatabaseDSN:       "postgres://user:pass@localhost/db",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				PostsPageLimit:    10,
				CommentsPageLimit: 10,
				MaxTitleLength:    255,
				MaxContentLength:  10000,
				MaxCommentLength:  2000,
				ChannelBufferSize: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDSNForTests(t *testing.T) {
	// Сохраняем оригинальную переменную
	originalTestDSN := os.Getenv("TEST_DB_DSN")
	defer func() {
		if originalTestDSN != "" {
			os.Setenv("TEST_DB_DSN", originalTestDSN)
		} else {
			os.Unsetenv("TEST_DB_DSN")
		}
	}()

	t.Run("использует TEST_DB_DSN если установлен", func(t *testing.T) {
		testDSN := "postgres://test:test@localhost/testdb"
		os.Setenv("TEST_DB_DSN", testDSN)

		cfg := &Config{DatabaseDSN: "postgres://prod:prod@localhost/proddb"}
		result := cfg.GetDSNForTests()

		if result != testDSN {
			t.Errorf("GetDSNForTests() = %v, expected %v", result, testDSN)
		}
	})

	t.Run("использует основной DSN если TEST_DB_DSN не установлен", func(t *testing.T) {
		os.Unsetenv("TEST_DB_DSN")

		prodDSN := "postgres://prod:prod@localhost/proddb"
		cfg := &Config{DatabaseDSN: prodDSN}
		result := cfg.GetDSNForTests()

		if result != prodDSN {
			t.Errorf("GetDSNForTests() = %v, expected %v", result, prodDSN)
		}
	})

	t.Run("возвращает дефолтный DSN если оба пусты", func(t *testing.T) {
		os.Unsetenv("TEST_DB_DSN")

		cfg := &Config{DatabaseDSN: ""}
		result := cfg.GetDSNForTests()

		expected := "postgres://user:password@localhost:5433/postsdb_test?sslmode=disable"
		if result != expected {
			t.Errorf("GetDSNForTests() = %v, expected %v", result, expected)
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getEnv", func(t *testing.T) {
		key := "TEST_GET_ENV_KEY"
		value := "test_value"
		defaultValue := "default_value"

		// Тест с установленной переменной
		os.Setenv(key, value)
		result := getEnv(key, defaultValue)
		if result != value {
			t.Errorf("getEnv() = %v, expected %v", result, value)
		}

		// Тест с неустановленной переменной
		os.Unsetenv(key)
		result = getEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getEnv() = %v, expected %v", result, defaultValue)
		}

		// Очистка
		os.Unsetenv(key)
	})

	t.Run("getIntEnv", func(t *testing.T) {
		key := "TEST_GET_INT_ENV_KEY"
		defaultValue := 42

		// Тест с валидным числом
		os.Setenv(key, "123")
		result := getIntEnv(key, defaultValue)
		if result != 123 {
			t.Errorf("getIntEnv() = %v, expected 123", result)
		}

		// Тест с невалидным числом
		os.Setenv(key, "not_a_number")
		result = getIntEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getIntEnv() = %v, expected %v", result, defaultValue)
		}

		// Тест с неустановленной переменной
		os.Unsetenv(key)
		result = getIntEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getIntEnv() = %v, expected %v", result, defaultValue)
		}

		// Очистка
		os.Unsetenv(key)
	})

	t.Run("getDurationEnv", func(t *testing.T) {
		key := "TEST_GET_DURATION_ENV_KEY"
		defaultValue := 30 * time.Second

		// Тест с валидной длительностью
		os.Setenv(key, "1m30s")
		result := getDurationEnv(key, defaultValue)
		expected := 90 * time.Second
		if result != expected {
			t.Errorf("getDurationEnv() = %v, expected %v", result, expected)
		}

		// Тест с невалидной длительностью
		os.Setenv(key, "not_a_duration")
		result = getDurationEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getDurationEnv() = %v, expected %v", result, defaultValue)
		}

		// Тест с неустановленной переменной
		os.Unsetenv(key)
		result = getDurationEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getDurationEnv() = %v, expected %v", result, defaultValue)
		}

		// Очистка
		os.Unsetenv(key)
	})

	t.Run("getBoolEnv", func(t *testing.T) {
		key := "TEST_GET_BOOL_ENV_KEY"
		defaultValue := false

		// Тест с валидным булевым значением
		os.Setenv(key, "true")
		result := getBoolEnv(key, defaultValue)
		if result != true {
			t.Errorf("getBoolEnv() = %v, expected true", result)
		}

		os.Setenv(key, "false")
		result = getBoolEnv(key, defaultValue)
		if result != false {
			t.Errorf("getBoolEnv() = %v, expected false", result)
		}

		// Тест с невалидным булевым значением
		os.Setenv(key, "not_a_bool")
		result = getBoolEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getBoolEnv() = %v, expected %v", result, defaultValue)
		}

		// Тест с неустановленной переменной
		os.Unsetenv(key)
		result = getBoolEnv(key, defaultValue)
		if result != defaultValue {
			t.Errorf("getBoolEnv() = %v, expected %v", result, defaultValue)
		}

		// Очистка
		os.Unsetenv(key)
	})
}
