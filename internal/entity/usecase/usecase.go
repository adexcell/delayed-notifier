package usecase

type Postgres interface{}

type Redis interface{}

type RabbitMQ interface{}

type UseCase struct {
	postgres Postgres
	redis Redis
	rabbit RabbitMQ
}

func New(postgres Postgres, redis Redis, rabbit RabbitMQ) *UseCase {
	return &UseCase{
		postgres: postgres,
		redis: redis,
		rabbit: rabbit,
	}
}
