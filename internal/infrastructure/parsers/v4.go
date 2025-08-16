package parsers

import (
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

// V4Parser parses Puma v4.x stats
type V4Parser struct{}

// NewV4Parser creates a new V4 parser
func NewV4Parser() *V4Parser {
	return &V4Parser{}
}

// Parse converts PumaStats to MetricCollection for v4.x
func (p *V4Parser) Parse(stats *infrastructure.PumaStats) (*domain.MetricCollection, error) {
	collection := domain.NewMetricCollection()
	timestamp := time.Now()

	// Basic worker metrics
	_ = collection.Add(domain.Metric{
		Name:      "workers",
		Value:     float64(stats.Workers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	_ = collection.Add(domain.Metric{
		Name:      "booted_workers",
		Value:     float64(stats.BootedWorkers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	_ = collection.Add(domain.Metric{
		Name:      "old_workers",
		Value:     float64(stats.OldWorkers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	// Phase (if available in v4)
	_ = collection.Add(domain.Metric{
		Name:      "phase",
		Value:     float64(stats.Phase),
		Type:      domain.MetricTypeGauge,
		Unit:      "phase",
		Timestamp: timestamp,
	})

	// For v4, worker status might be simpler
	var totalBacklog, totalRunning int
	for _, worker := range stats.WorkerStatus {
		totalBacklog += worker.LastStatus.Backlog
		totalRunning += worker.LastStatus.Running
	}

	if len(stats.WorkerStatus) > 0 {
		_ = collection.Add(domain.Metric{
			Name:      "backlog",
			Value:     float64(totalBacklog),
			Type:      domain.MetricTypeGauge,
			Unit:      "requests",
			Timestamp: timestamp,
		})

		_ = collection.Add(domain.Metric{
			Name:      "running",
			Value:     float64(totalRunning),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	// Single mode metrics
	if stats.Backlog != nil {
		_ = collection.Add(domain.Metric{
			Name:      "backlog",
			Value:     float64(*stats.Backlog),
			Type:      domain.MetricTypeGauge,
			Unit:      "requests",
			Timestamp: timestamp,
		})
	}

	if stats.Running != nil {
		_ = collection.Add(domain.Metric{
			Name:      "running",
			Value:     float64(*stats.Running),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	return collection, nil
}