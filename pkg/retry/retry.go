package retry

import (
	"context"
	"time"
)

// Config holds the configuration parameters for retry logic.
type Config struct {
	Attempts int           `default:"5"   envconfig:"RETRY_ATTEMPTS"`
	Delay    time.Duration `default:"1s"  envconfig:"RETRY_DELAY"`
	Backoff  float64       `default:"2" envconfig:"RETRY_BACKOFF"`
}

// DoWithContext executes the provided function with retries based on the configuration and context.
func DoWithContext(ctx context.Context, c Config, fn func() error) error {
	var err error

	delay := c.Delay

	for range c.Attempts {
		err = fn()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * c.Backoff)
	}

	return err
}
