# ЭТАП 1: Сборка приложения
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник Go
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./main.go

# ЭТАП 2: Легковесный финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный файл из первого этапа
COPY --from=builder /payment-service .

# Открываем порт для сервиса
EXPOSE 8080

# Команда для запуска
CMD ["./payment-service"]