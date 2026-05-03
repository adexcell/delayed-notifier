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

// Client is the base struct for handling connection recovery, consumption and
// publishing. Note that this struct has an internal mutex to safeguard against
// data races. As you develop and iterate over this example, you may need to add
// further locks, or safeguards, to keep your application safe from data races.
type Client struct {
	addr      string
	wg        sync.WaitGroup
	cfg       Config
	mu        sync.Mutex
	conn      *amqp.Connection
	pubMu     sync.Mutex
	pubCh     *amqp.Channel
	subMu     sync.Mutex
	subCh     *amqp.Channel
	ctx       context.Context
	cancel    context.CancelFunc
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

// New creates a new broker state instance, and automatically
// attempts to connect to the server.
func New(ctx context.Context, c Config) (*Client, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
	)

	ctx, cancel := context.WithCancel(ctx)

	client := Client{
		addr:   addr,
		mu:     sync.Mutex{},
		wg:     sync.WaitGroup{},
		cfg:    c,
		ctx:    ctx,
		cancel: cancel,
	}

	err := client.start()
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	return &client, nil
}

func (c *Client) start() error {
	c.wg.Go(c.connect)
	c.wg.Go(c.managePubChannel)
	c.wg.Go(c.manageConsumeChannel)

	return nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		c.cancel()

		conn := c.getConn()

		c.clearAndCloseConn(conn)
		c.wg.Wait()
	})
}
