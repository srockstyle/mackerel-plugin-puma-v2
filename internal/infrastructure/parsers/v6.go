package parsers

import (
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

// V6Parser parses Puma v6.x stats
type V6Parser struct{}

// NewV6Parser creates a new V6 parser
func NewV6Parser() *V6Parser {
	return &V6Parser{}
}

// Parse converts PumaStats to MetricCollection
func (p *V6Parser) Parse(stats *infrastructure.PumaStats) (*domain.MetricCollection, error) {
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

	// New in Puma 6.x: Request count
	if stats.RequestsCount != nil {
		collection.Add(domain.Metric{
			Name:      "requests_count",
			Value:     float64(*stats.RequestsCount),
			Type:      domain.MetricTypeCounter,
			Unit:      "requests",
			Timestamp: timestamp,
		})
	}

	// New in Puma 6.x: Uptime
	if stats.Uptime != nil {
		collection.Add(domain.Metric{
			Name:      "uptime",
			Value:     float64(*stats.Uptime),
			Type:      domain.MetricTypeGauge,
			Unit:      "seconds",
			Timestamp: timestamp,
		})
	}

	// Process worker status
	var totalBacklog, totalRunning, totalPoolCapacity, totalMaxThreads int
	for _, worker := range stats.WorkerStatus {
		totalBacklog += worker.LastStatus.Backlog
		totalRunning += worker.LastStatus.Running
		totalPoolCapacity += worker.LastStatus.PoolCapacity
		totalMaxThreads += worker.LastStatus.MaxThreads
	}

	// Calculate thread utilization
	if totalPoolCapacity > 0 {
		utilization := (float64(totalRunning) / float64(totalPoolCapacity)) * 100
		collection.Add(domain.Metric{
			Name:      "thread_utilization",
			Value:     utilization,
			Type:      domain.MetricTypeGauge,
			Unit:      "percentage",
			Timestamp: timestamp,
		})
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