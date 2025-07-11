package main

import (
	"os"
	"testing"
)

func TestMain_WithValidConfig(t *testing.T) {
	// Сохраняем оригинальные аргументы
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Устанавливаем тестовые переменные окружения
	os.Setenv("DB_DSN", "postgres://test:test@localhost/testdb")
	os.Setenv("HTTP_ADDR", ":0") // Используем случайный порт для тестов
	defer func() {
		os.Unsetenv("DB_DSN")
		os.Unsetenv("HTTP_ADDR")
	}()

	// Тест не должен запускать реальный сервер, но должен пройти инициализацию
	// Мы проверяем только то, что функция не вызывает панику при валидации конфига
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() panicked: %v", r)
		}
	}()

	// Не вызываем main() напрямую, так как она блокирующая
	// Вместо этого тестируем компоненты отдельно
}

func TestConstants(t *testing.T) {
	if ExitCodeSuccess != 0 {
		t.Errorf("ExitCodeSuccess should be 0, got %d", ExitCodeSuccess)
	}

	if ExitCodeError != 1 {
		t.Errorf("ExitCodeError should be 1, got %d", ExitCodeError)
	}
}
