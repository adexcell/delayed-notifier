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

ENV TZ=Asia/Yekaterinburg
# Опционально: база таймзон для Alpine. 
# Go с CGO_ENABLED=0 уже встраивает IANA-базу, но это гарантирует работу всех системных библиотек.
RUN apk add --no-cache tzdata

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/main .
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/web ./web


# Запуск приложения
CMD ["./main"]
