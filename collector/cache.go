package collector

import (
	"sync"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

var (
	cliTimeout          = kingpin.Flag("cli.timeout", "Timeout for a Slurm CLI command.").Envar("CLI_TIMEOUT").Default("10s").Duration()
	defaultPollInterval = kingpin.Flag("collector.cache-ttl", "How long shared Slurm CLI results are cached.").Default("30s").Duration()
)

// throttledCache serializes refreshes. A caller that arrives while a refresh is
// in flight waits for it and receives the same result rather than starting a
// second Slurm command.
type throttledCache[T any] struct {
	mu        sync.Mutex
	fetch     func() (T, error)
	ttl       func() time.Duration
	value     T
	err       error
	fetchedAt time.Time
}

func newThrottledCache[T any](fetch func() (T, error), ttl func() time.Duration) *throttledCache[T] {
	return &throttledCache[T]{fetch: fetch, ttl: ttl}
}

func (c *throttledCache[T]) Get() (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.fetchedAt.IsZero() && time.Since(c.fetchedAt) < c.ttl() {
		return c.value, c.err
	}
	c.value, c.err = c.fetch()
	c.fetchedAt = time.Now()
	return c.value, c.err
}
