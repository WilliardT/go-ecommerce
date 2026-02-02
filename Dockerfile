# Многоступенчатая сборка для оптимизации размера образа
# Используем latest для совместимости с Go 1.25+
FROM golang:alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/main .

# Копируем .env файл (опционально, лучше использовать environment в docker-compose)
# COPY --from=builder /app/.env .

# Открываем порт
EXPOSE 8000

# Запускаем приложение
CMD ["./main"]
