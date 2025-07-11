package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/NarthurN/CommentsSystem/internal/repository"
)

// APIError представляет структурированную ошибку API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse представляет стандартный формат ответа с ошибкой
type ErrorResponse struct {
	Error   APIError `json:"error"`
	Success bool     `json:"success"`
}

// Коды ошибок для клиентов
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeRateLimit        = "RATE_LIMIT_EXCEEDED"
	ErrCodeTooLarge         = "PAYLOAD_TOO_LARGE"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnavailable      = "SERVICE_UNAVAILABLE"
	ErrCodeDuplicate        = "DUPLICATE_ENTITY"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeCommentsDisabled = "COMMENTS_DISABLED"
)

// ErrorHandler обрабатывает ошибки и возвращает правильные HTTP коды
type ErrorHandler struct {
	logger *log.Logger
}

// NewErrorHandler создает новый обработчик ошибок
func NewErrorHandler(logger *log.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// HandleError обрабатывает ошибку и возвращает соответствующий HTTP статус и ответ
func (h *ErrorHandler) HandleError(ctx context.Context, err error) (int, ErrorResponse) {
	if err == nil {
		return http.StatusOK, ErrorResponse{}
	}

	// Логируем ошибку для мониторинга
	h.logError(ctx, err)

	// Определяем тип ошибки и возвращаем соответствующий код
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeNotFound,
				Message: "Requested resource not found",
				Details: "The entity you are looking for does not exist or has been removed",
			},
			Success: false,
		}

	case errors.Is(err, repository.ErrDuplicate):
		return http.StatusConflict, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeDuplicate,
				Message: "Resource already exists",
				Details: "An entity with the same identifier already exists",
			},
			Success: false,
		}

	case errors.Is(err, repository.ErrInvalidInput):
		return http.StatusBadRequest, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeInvalidInput,
				Message: "Invalid input provided",
				Details: err.Error(),
			},
			Success: false,
		}

	case isValidationError(err):
		return http.StatusBadRequest, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeValidation,
				Message: "Validation failed",
				Details: err.Error(),
			},
			Success: false,
		}

	case isCommentsDisabledError(err):
		return http.StatusForbidden, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeCommentsDisabled,
				Message: "Comments are disabled for this post",
				Details: "The post author has disabled comments for this post",
			},
			Success: false,
		}

	case isRateLimitError(err):
		return http.StatusTooManyRequests, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeRateLimit,
				Message: "Rate limit exceeded",
				Details: "Too many requests, please try again later",
			},
			Success: false,
		}

	case isPayloadTooLargeError(err):
		return http.StatusRequestEntityTooLarge, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeTooLarge,
				Message: "Payload too large",
				Details: err.Error(),
			},
			Success: false,
		}

	case errors.Is(err, repository.ErrConnectionFailed):
		return http.StatusServiceUnavailable, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeUnavailable,
				Message: "Service temporarily unavailable",
				Details: "Database connection failed, please try again later",
			},
			Success: false,
		}

	default:
		// Внутренняя ошибка сервера - не раскрываем детали
		return http.StatusInternalServerError, ErrorResponse{
			Error: APIError{
				Code:    ErrCodeInternal,
				Message: "Internal server error",
				Details: "An unexpected error occurred, please try again later",
			},
			Success: false,
		}
	}
}

// logError логирует ошибку с контекстом
func (h *ErrorHandler) logError(ctx context.Context, err error) {
	if h.logger == nil {
		return
	}

	// Добавляем контекстную информацию для отладки
	h.logger.Printf("ERROR: %v", err)
}

// Вспомогательные функции для определения типов ошибок

func isValidationError(err error) bool {
	errStr := err.Error()
	return contains(errStr, "validation") ||
		contains(errStr, "invalid") ||
		contains(errStr, "must be") ||
		contains(errStr, "required") ||
		contains(errStr, "exceed") ||
		contains(errStr, "between") ||
		contains(errStr, "negative")
}

func isCommentsDisabledError(err error) bool {
	return contains(err.Error(), "comments are disabled")
}

func isRateLimitError(err error) bool {
	return contains(err.Error(), "rate limit") ||
		contains(err.Error(), "too many requests")
}

func isPayloadTooLargeError(err error) bool {
	return contains(err.Error(), "too large") ||
		contains(err.Error(), "exceed") &&
			(contains(err.Error(), "characters") || contains(err.Error(), "symbols"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(containsAtIndex(s, substr, 0) ||
				containsAtIndex(s, substr, len(s)-len(substr)) ||
				containsInMiddle(s, substr))))
}

func containsAtIndex(s, substr string, index int) bool {
	if index < 0 || index+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[index+i] != substr[i] {
			return false
		}
	}
	return true
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr); i++ {
		if containsAtIndex(s, substr, i) {
			return true
		}
	}
	return false
}

// GraphQLErrorHandler специально для обработки GraphQL ошибок
type GraphQLErrorHandler struct {
	*ErrorHandler
}

// NewGraphQLErrorHandler создает обработчик для GraphQL ошибок
func NewGraphQLErrorHandler(logger *log.Logger) *GraphQLErrorHandler {
	return &GraphQLErrorHandler{
		ErrorHandler: NewErrorHandler(logger),
	}
}

// FormatGraphQLError форматирует ошибку для GraphQL ответа
func (h *GraphQLErrorHandler) FormatGraphQLError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Логируем ошибку
	h.logError(ctx, err)

	// Возвращаем пользователю дружественное сообщение
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fmt.Errorf("requested resource not found")

	case errors.Is(err, repository.ErrInvalidInput):
		return fmt.Errorf("invalid input: %s", err.Error())

	case isValidationError(err):
		return fmt.Errorf("validation error: %s", err.Error())

	case isCommentsDisabledError(err):
		return fmt.Errorf("comments are disabled for this post")

	case errors.Is(err, repository.ErrConnectionFailed):
		return fmt.Errorf("service temporarily unavailable")

	default:
		// Не раскрываем внутренние ошибки
		return fmt.Errorf("an error occurred while processing your request")
	}
}
