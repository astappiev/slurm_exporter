package collector

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottledCacheCachesAndRefreshes(t *testing.T) {
	var calls atomic.Int32
	cache := newThrottledCache(func() (int, error) { return int(calls.Add(1)), nil }, func() time.Duration { return 20 * time.Millisecond })

	first, err := cache.Get()
	assert.NoError(t, err)
	second, err := cache.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, first)
	assert.Equal(t, 1, second)
	assert.Equal(t, int32(1), calls.Load())

	time.Sleep(25 * time.Millisecond)
	third, err := cache.Get()
	assert.NoError(t, err)
	assert.Equal(t, 2, third)
	assert.Equal(t, int32(2), calls.Load())
}

func TestThrottledCacheCoalescesConcurrentGets(t *testing.T) {
	var calls atomic.Int32
	cache := newThrottledCache(func() (int, error) {
		calls.Add(1)
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	}, func() time.Duration { return time.Minute })

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() { _, err := cache.Get(); assert.NoError(t, err) })
	}
	wg.Wait()
	assert.Equal(t, int32(1), calls.Load())
}
