# Stage 1: Builder
FROM golang:1.25 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
# CGO_ENABLED=0 for static binary, -ldflags="-s -w" to strip debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

# Install common runtime dependencies (certs, timezone)
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /app/app .

# Copy configuration (required by config.go)
COPY --from=builder /app/config/config.yaml ./config/config.yaml

# Copy static files (frontend)
COPY --from=builder /app/static ./static

# Switch to non-root user
USER appuser

CMD ["./app"]
