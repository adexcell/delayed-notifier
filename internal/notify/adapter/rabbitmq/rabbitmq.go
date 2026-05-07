package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds the configuration parameters for the RabbitMQ client.
type Config struct {
	User     string `default:"" envconfig:"RABBIT_USER"`
	Password string `default:"" envconfig:"RABBIT_PASSWORD"`
	Host     string `default:"" envconfig:"RABBIT_HOST"`
	Port     int    `default:"" envconfig:"RABBIT_PORT"`

	DelayedExchange string `default:"notify.delayexchange" envconfig:"RABBIT_EXCHANGE"`
	MainQueue       string `default:"notify.queue"         envconfig:"RABBIT_QUEUE_NAME"`
	RoutingKey      string `default:"notify.delay"         envconfig:"RABBIT_ROUTING_KEY"`

	DLXExchange   string `default:"notify.dlx"  envconfig:"RABBIT_DLX_EXCHANGE"`
	DLQQueue      string `default:"notify.dlq"  envconfig:"RABBIT_DLQ_QUEUE"`
	DLQRoutingKey string `default:"notify.dead" envconfig:"RABBIT_DLQ_ROUTING_KEY"`
}

// Client handles RabbitMQ connections, publishing, and consumption with automatic recovery.
type Client struct {
	addr      string
	wg        *sync.WaitGroup
	cfg       Config
	mu        sync.Mutex
	conn      *amqp.Connection
	pubMu     sync.Mutex
	pubCh     *amqp.Channel
	subMu     sync.Mutex
	subCh     *amqp.Channel
	closeOnce sync.Once
}

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	prefetchCount  = 10
)

var (
	errNotConnected = errors.New("not connected")
	errMsgNacked    = errors.New("message was nack'ed by broker")
)

// New creates a new RabbitMQ client instance and starts connection management.
func New(ctx context.Context, c Config) (*Client, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
	)

	client := Client{
		addr:   addr,
		mu:     sync.Mutex{},
		wg:     new(sync.WaitGroup),
		cfg:    c,
	}

	err := client.start(ctx)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	return &client, nil
}

func (c *Client) start(ctx context.Context) error {
	c.Go(ctx, c.connect)
	c.Go(ctx, c.managePubChannel)
	c.Go(ctx, c.manageConsumeChannel)

	return nil
}

// Close gracefully shuts down the RabbitMQ client and closes all connections.
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		conn := c.getConn()
		c.clearAndCloseConn(conn)
		c.wg.Wait()
	})
}
