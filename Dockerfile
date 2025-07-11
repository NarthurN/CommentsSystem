# --- Сборка ---
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git

WORKDIR /app

# Копируем зависимости и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/app

# --- Финальный образ ---
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из стадии сборки
COPY --from=builder /app/server .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Указываем порт, на котором работает приложение
EXPOSE 8080

# Команда для запуска
CMD ["./server"]
