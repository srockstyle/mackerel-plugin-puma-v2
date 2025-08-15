package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// PumaClient interfaces with Puma control server
type PumaClient interface {
	GetStats(ctx context.Context) (*PumaStats, error)
	GetGCStats(ctx context.Context) (map[string]interface{}, error)
}

// PumaStats represents Puma statistics
type PumaStats struct {
	Workers       int                    `json:"workers"`
	Phase         int                    `json:"phase"`
	BootedWorkers int                    `json:"booted_workers"`
	OldWorkers    int                    `json:"old_workers"`
	WorkerStatus  []WorkerStatus         `json:"worker_status"`
	// Single mode fields
	Backlog      *int                   `json:"backlog,omitempty"`
	Running      *int                   `json:"running,omitempty"`
	PoolCapacity *int                   `json:"pool_capacity,omitempty"`
	MaxThreads   *int                   `json:"max_threads,omitempty"`
	// Puma 6.x fields
	RequestsCount *int64                `json:"requests_count,omitempty"`
	Uptime        *int                  `json:"uptime,omitempty"`
}

// WorkerStatus represents individual worker status
type WorkerStatus struct {
	PID          int        `json:"pid"`
	Index        int        `json:"index"`
	Phase        int        `json:"phase"`
	Booted       bool       `json:"booted"`
	LastCheckin  string     `json:"last_checkin"`
	LastStatus   LastStatus `json:"last_status"`
}

// LastStatus represents worker's last status
type LastStatus struct {
	Backlog      int `json:"backlog"`
	Running      int `json:"running"`
	PoolCapacity int `json:"pool_capacity"`
	MaxThreads   int `json:"max_threads"`
}

// DefaultPumaClient is the default implementation of PumaClient
type DefaultPumaClient struct {
	client        *UnixSocketClient
	retryCount    int
	retryInterval time.Duration
}

// NewPumaClient creates a new Puma client
func NewPumaClient(socketPath string, timeout time.Duration) PumaClient {
	return &DefaultPumaClient{
		client:        NewUnixSocketClient(socketPath, timeout),
		retryCount:    3,
		retryInterval: 1 * time.Second,
	}
}

// GetStats retrieves Puma statistics with retry logic
func (c *DefaultPumaClient) GetStats(ctx context.Context) (*PumaStats, error) {
	var lastErr error

	for i := range c.retryCount + 1 {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryInterval):
			}
		}

		stats, err := c.fetchStats()
		if err == nil {
			return stats, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.retryCount+1, lastErr)
}

// fetchStats performs the actual fetch operation
func (c *DefaultPumaClient) fetchStats() (*PumaStats, error) {
	body, err := c.client.Get("/stats")
	if err != nil {
		return nil, err
	}

	var stats PumaStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("parsing stats JSON: %w", err)
	}

	return &stats, nil
}

// GetGCStats retrieves GC statistics
func (c *DefaultPumaClient) GetGCStats(ctx context.Context) (map[string]interface{}, error) {
	body, err := c.client.Get("/gc-stats")
	if err != nil {
		return nil, err
	}

	var gcStats map[string]interface{}
	if err := json.Unmarshal(body, &gcStats); err != nil {
		return nil, fmt.Errorf("parsing gc-stats JSON: %w", err)
	}

	return gcStats, nil
}