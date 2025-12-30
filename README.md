```text
delayed-notifier/
├── cmd/
│   ├── app/
│   │    └── app.go          # wiring
│   └── main.go              # только запуск
├── internal/
│   ├── config/              # загрузка конфигов
│   ├── logger/              # инициализация zerolog/zap
│   ├── httpserver/          # обёртка над net/http + graceful shutdown
│   ├── transport/
│   │   └── http/            # роутер, handlers, middleware
│   ├── domain/              # сущности и интерфейсы (ports)
│   ├── usecase/             # бизнес-логика (application layer)
│   ├── repository/          # реализация портов (DB, external API)
│   └── observability/       # metrics (Prometheus), tracing (OTel)
├── pkg/                     # общие утилиты (если нужны вне сервиса)
├── config/                  # config
├── Dockerfile
├── docker-compose.yml
└── Makefile
```