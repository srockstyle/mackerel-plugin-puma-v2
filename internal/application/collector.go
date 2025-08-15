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
	client          infrastructure.PumaClient
	parserFactory   *parsers.ParserFactory
	versionDetector *infrastructure.VersionDetector
	detectedVersion string
	retryCount      int
	retryInterval   time.Duration
	logger          *log.Logger
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config *Config, logger *log.Logger) *MetricsCollector {
	client := infrastructure.NewPumaClient(config.SocketPath, config.Timeout)
	
	return &MetricsCollector{
		client:          client,
		parserFactory:   parsers.NewParserFactory(),
		versionDetector: infrastructure.NewVersionDetector(client),
		detectedVersion: "",
		retryCount:      config.RetryCount,
		retryInterval:   config.RetryInterval,
		logger:          logger,
	}
}

// Collect collects metrics from Puma
func (c *MetricsCollector) Collect(ctx context.Context) (*domain.MetricCollection, error) {
	var lastErr error

	for i := range c.retryCount + 1 {
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
	// Detect version if not already detected
	if c.detectedVersion == "" {
		version, err := c.versionDetector.DetectVersion(ctx)
		if err != nil {
			c.logger.Printf("Failed to detect Puma version, assuming 6.x: %v", err)
			c.detectedVersion = "6.x"
		} else {
			c.detectedVersion = version
			c.logger.Printf("Detected Puma version: %s", version)
		}
	}

	stats, err := c.client.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	// Get appropriate parser for the detected version
	parser := c.parserFactory.GetParser(c.detectedVersion)
	collection, err := parser.Parse(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return collection, nil
}