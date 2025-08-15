package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure/parsers"
)

// MetricsCollector collects metrics from Puma
type MetricsCollector struct {
	client        infrastructure.PumaClient
	parser        *parsers.V6Parser
	retryCount    int
	retryInterval time.Duration
	logger        *log.Logger
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config *Config, logger *log.Logger) *MetricsCollector {
	client := infrastructure.NewPumaClient(config.SocketPath, config.Timeout)
	
	return &MetricsCollector{
		client:        client,
		parser:        parsers.NewV6Parser(),
		retryCount:    config.RetryCount,
		retryInterval: config.RetryInterval,
		logger:        logger,
	}
}

// Collect collects metrics from Puma
func (c *MetricsCollector) Collect(ctx context.Context) (*domain.MetricCollection, error) {
	var lastErr error

	for i := 0; i <= c.retryCount; i++ {
		if i > 0 {
			c.logger.Printf("Retry attempt %d/%d after %v", i, c.retryCount, c.retryInterval)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryInterval):
			}
		}

		stats, err := c.collectWithTimeout(ctx)
		if err == nil {
			return stats, nil
		}

		lastErr = err
		c.logger.Printf("Collection failed: %v", err)
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.retryCount+1, lastErr)
}

// collectWithTimeout performs collection with timeout
func (c *MetricsCollector) collectWithTimeout(ctx context.Context) (*domain.MetricCollection, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	statsChan := make(chan *domain.MetricCollection, 1)
	errChan := make(chan error, 1)

	go func() {
		stats, err := c.fetchAndParse(ctx)
		if err != nil {
			errChan <- err
			return
		}
		statsChan <- stats
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("collection timeout: %w", ctx.Err())
	case err := <-errChan:
		return nil, err
	case stats := <-statsChan:
		return stats, nil
	}
}

// fetchAndParse fetches stats from Puma and parses them
func (c *MetricsCollector) fetchAndParse(ctx context.Context) (*domain.MetricCollection, error) {
	stats, err := c.client.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	collection, err := c.parser.Parse(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return collection, nil
}