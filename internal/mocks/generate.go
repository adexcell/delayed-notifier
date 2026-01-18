package mocks

//go:generate mockgen -destination=mock_postgres.go -package=mocks github.com/adexcell/delayed-notifier/internal/domain NotifyPostgres
//go:generate mockgen -destination=mock_redis.go -package=mocks github.com/adexcell/delayed-notifier/internal/domain NotifyRedis
//go:generate mockgen -destination=mock_queue.go -package=mocks github.com/adexcell/delayed-notifier/internal/domain QueueProvider
//go:generate mockgen -destination=mock_usecase.go -package=mocks github.com/adexcell/delayed-notifier/internal/domain NotifyUsecase
//go:generate mockgen -destination=mock_sender.go -package=mocks github.com/adexcell/delayed-notifier/internal/domain Sender
