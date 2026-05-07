# Этап сборки
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Копируем только модули для кэширования
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем код
COPY . . 

# Статическая сборка минимального бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

# Финальный минимальный образ
FROM alpine:3.22

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/main .
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/web ./web


# Запуск приложения
CMD ["./main"]
