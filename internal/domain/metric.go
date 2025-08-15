package domain

import (
	"fmt"
	"time"
)

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

// Metric represents a single metric data point
type Metric struct {
	Name      string
	Value     float64
	Type      MetricType
	Unit      string
	Timestamp time.Time
	Labels    map[string]string
}

// Validate checks if the metric is valid
func (m *Metric) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}
	if m.Type != MetricTypeGauge && m.Type != MetricTypeCounter {
		return fmt.Errorf("invalid metric type: %s", m.Type)
	}
	return nil
}

// MetricCollection holds multiple metrics
type MetricCollection struct {
	metrics []Metric
}

// NewMetricCollection creates a new metric collection
func NewMetricCollection() *MetricCollection {
	return &MetricCollection{
		metrics: make([]Metric, 0),
	}
}

// Add adds a metric to the collection
func (mc *MetricCollection) Add(metric Metric) error {
	if err := metric.Validate(); err != nil {
		return fmt.Errorf("invalid metric: %w", err)
	}
	mc.metrics = append(mc.metrics, metric)
	return nil
}

// All returns all metrics in the collection
func (mc *MetricCollection) All() []Metric {
	return mc.metrics
}

// Filter returns metrics matching the given predicate
func (mc *MetricCollection) Filter(predicate func(Metric) bool) []Metric {
	filtered := make([]Metric, 0)
	for _, m := range mc.metrics {
		if predicate(m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}