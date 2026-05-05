package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// sleep returns true if the timer expires before the client context is canceled.
func (c *Client) sleep(ctx context.Context, d time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(d):
		return true
	}
}

// getConn return current conn.
func (c *Client) getConn() *amqp.Connection {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn
}

// Go calls f in a new goroutine and adds that task to the [WaitGroup].
// When f returns, the task is removed from the WaitGroup.
//
// The function f must not panic.
func (c *Client) Go(ctx context.Context, f func(ctx context.Context)) {
	c.wg.Add(1)
	go func() {
		defer func() {
			if x := recover(); x != nil {
				// f panicked, which will be fatal because
				// this is a new goroutine.
				//
				// Calling Done will unblock Wait in the main goroutine,
				// allowing it to race with the fatal panic and
				// possibly even exit the process (os.Exit(0))
				// before the panic completes.
				//
				// This is almost certainly undesirable,
				// so instead avoid calling Done and simply panic.
				panic(x)
			}

			// f completed normally, or abruptly using goexit.
			// Either way, decrement the semaphore.
			c.wg.Done()
		}()
		f(ctx)
	}()
}
