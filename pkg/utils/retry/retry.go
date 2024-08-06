package retry

import (
	"context"
	"fmt"
	"time"
)

// Retrier represents a retry mechanism with constant backoff.
type Retrier struct {
	attempts int
	delay    time.Duration
}

// NewRetrier creates a new Retrier instance.
func NewRetrier(attempts int, delay time.Duration) *Retrier {
	return &Retrier{
		attempts: attempts,
		delay:    delay,
	}
}

// Retry executes the given function repeatedly until it succeeds or the maximum number of attempts is reached.
func (r *Retrier) Retry(ctx context.Context, fn func() error) error {
	for attempt := 1; attempt <= r.attempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Sleep for the specified delay before the next attempt.
		select {
		case <-time.After(r.delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("retry attempts exhausted after %d tries", r.attempts)
}
