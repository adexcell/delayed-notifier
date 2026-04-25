package rabbit

import "github.com/wb-go/wbf/rabbitmq"

type Rabbit struct {
	client *rabbitmq.RabbitClient
}

func New(client *rabbitmq.RabbitClient) *Rabbit {
	return &Rabbit{client: client}
}
