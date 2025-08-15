package parsers

import (
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

// V5Parser parses Puma v5.x stats
type V5Parser struct{}

// NewV5Parser creates a new V5 parser
func NewV5Parser() *V5Parser {
	return &V5Parser{}
}

// Parse converts PumaStats to MetricCollection for v5.x
func (p *V5Parser) Parse(stats *infrastructure.PumaStats) (*domain.MetricCollection, error) {
	collection := domain.NewMetricCollection()
	timestamp := time.Now()

	// Worker metrics
	collection.Add(domain.Metric{
		Name:      "workers",
		Value:     float64(stats.Workers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	collection.Add(domain.Metric{
		Name:      "booted_workers",
		Value:     float64(stats.BootedWorkers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	collection.Add(domain.Metric{
		Name:      "old_workers",
		Value:     float64(stats.OldWorkers),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})

	collection.Add(domain.Metric{
		Name:      "phase",
		Value:     float64(stats.Phase),
		Type:      domain.MetricTypeGauge,
		Unit:      "phase",
		Timestamp: timestamp,
	})

	// Process worker status
	var totalBacklog, totalRunning, totalPoolCapacity, totalMaxThreads int
	for _, worker := range stats.WorkerStatus {
		totalBacklog += worker.LastStatus.Backlog
		totalRunning += worker.LastStatus.Running
		totalPoolCapacity += worker.LastStatus.PoolCapacity
		totalMaxThreads += worker.LastStatus.MaxThreads
	}

	// Aggregate worker metrics
	if len(stats.WorkerStatus) > 0 {
		collection.Add(domain.Metric{
			Name:      "backlog",
			Value:     float64(totalBacklog),
			Type:      domain.MetricTypeGauge,
			Unit:      "requests",
			Timestamp: timestamp,
		})

		collection.Add(domain.Metric{
			Name:      "running",
			Value:     float64(totalRunning),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})

		collection.Add(domain.Metric{
			Name:      "pool_capacity",
			Value:     float64(totalPoolCapacity),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})

		collection.Add(domain.Metric{
			Name:      "max_threads",
			Value:     float64(totalMaxThreads),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	// Single mode metrics
	if stats.Backlog != nil {
		collection.Add(domain.Metric{
			Name:      "backlog",
			Value:     float64(*stats.Backlog),
			Type:      domain.MetricTypeGauge,
			Unit:      "requests",
			Timestamp: timestamp,
		})
	}

	if stats.Running != nil {
		collection.Add(domain.Metric{
			Name:      "running",
			Value:     float64(*stats.Running),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	if stats.PoolCapacity != nil {
		collection.Add(domain.Metric{
			Name:      "pool_capacity",
			Value:     float64(*stats.PoolCapacity),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	if stats.MaxThreads != nil {
		collection.Add(domain.Metric{
			Name:      "max_threads",
			Value:     float64(*stats.MaxThreads),
			Type:      domain.MetricTypeGauge,
			Unit:      "threads",
			Timestamp: timestamp,
		})
	}

	return collection, nil
}