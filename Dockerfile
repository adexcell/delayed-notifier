FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod tidy

RUN go build -o app ./cmd/main.go

CMD ["./app"]
