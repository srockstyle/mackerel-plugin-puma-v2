package application

import (
	"context"
	"encoding/json"
	"runtime"
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure/parsers"
)

// ExtendedMetricsCollector collects extended metrics including system stats
type ExtendedMetricsCollector struct {
	baseCollector *MetricsCollector
}

// NewExtendedMetricsCollector creates a new extended metrics collector
func NewExtendedMetricsCollector(base *MetricsCollector) *ExtendedMetricsCollector {
	return &ExtendedMetricsCollector{
		baseCollector: base,
	}
}

// CollectWithSystemMetrics collects both Puma and system metrics
func (c *ExtendedMetricsCollector) CollectWithSystemMetrics(ctx context.Context) (*domain.MetricCollection, error) {
	// Get base metrics
	collection, err := c.baseCollector.Collect(ctx)
	if err != nil {
		return nil, err
	}

	// Add system metrics
	c.addMemoryMetrics(collection)
	c.addGoroutineMetrics(collection)
	c.addUptimeMetrics(collection)

	// Try to get GC stats if available
	c.tryAddGCMetrics(ctx, collection)

	return collection, nil
}

// addMemoryMetrics adds memory usage metrics
func (c *ExtendedMetricsCollector) addMemoryMetrics(collection *domain.MetricCollection) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	timestamp := time.Now()

	// Total allocated memory
	collection.Add(domain.Metric{
		Name:      "memory.alloc",
		Value:     float64(m.Alloc) / 1024 / 1024, // Convert to MB
		Type:      domain.MetricTypeGauge,
		Unit:      "megabytes",
		Timestamp: timestamp,
	})

	// Total system memory
	collection.Add(domain.Metric{
		Name:      "memory.sys",
		Value:     float64(m.Sys) / 1024 / 1024, // Convert to MB
		Type:      domain.MetricTypeGauge,
		Unit:      "megabytes",
		Timestamp: timestamp,
	})

	// Heap memory in use
	collection.Add(domain.Metric{
		Name:      "memory.heap_inuse",
		Value:     float64(m.HeapInuse) / 1024 / 1024, // Convert to MB
		Type:      domain.MetricTypeGauge,
		Unit:      "megabytes",
		Timestamp: timestamp,
	})

	// Number of GC cycles
	collection.Add(domain.Metric{
		Name:      "gc.num_gc",
		Value:     float64(m.NumGC),
		Type:      domain.MetricTypeCounter,
		Unit:      "count",
		Timestamp: timestamp,
	})
}

// addGoroutineMetrics adds goroutine metrics
func (c *ExtendedMetricsCollector) addGoroutineMetrics(collection *domain.MetricCollection) {
	timestamp := time.Now()

	collection.Add(domain.Metric{
		Name:      "go.goroutines",
		Value:     float64(runtime.NumGoroutine()),
		Type:      domain.MetricTypeGauge,
		Unit:      "count",
		Timestamp: timestamp,
	})
}

// addUptimeMetrics adds plugin uptime metrics
func (c *ExtendedMetricsCollector) addUptimeMetrics(collection *domain.MetricCollection) {
	// This is a placeholder - in a real implementation, you'd track
	// when the plugin started and calculate uptime
	timestamp := time.Now()

	collection.Add(domain.Metric{
		Name:      "plugin.uptime",
		Value:     0, // Would be calculated from start time
		Type:      domain.MetricTypeGauge,
		Unit:      "seconds",
		Timestamp: timestamp,
	})
}

// tryAddGCMetrics attempts to add GC metrics from Puma
func (c *ExtendedMetricsCollector) tryAddGCMetrics(ctx context.Context, collection *domain.MetricCollection) {
	gcStats, err := c.baseCollector.client.GetGCStats(ctx)
	if err != nil {
		// GC stats might not be available, especially in newer Puma versions
		c.baseCollector.logger.Printf("GC stats not available: %v", err)
		return
	}

	// Convert to JSON bytes for the parser
	gcData, err := json.Marshal(gcStats)
	if err != nil {
		c.baseCollector.logger.Printf("Failed to marshal GC stats: %v", err)
		return
	}

	// Use the detailed GC parser
	gcParser := &parsers.GCParser{}
	gcMetrics, err := gcParser.ParseGCStats(gcData)
	if err != nil {
		c.baseCollector.logger.Printf("Failed to parse GC stats: %v", err)
		return
	}

	// Add all parsed GC metrics to the collection
	timestamp := time.Now()
	for _, metric := range gcMetrics.All() {
		metric.Timestamp = timestamp
		collection.Add(metric)
	}
}

// getFloat64FromInterface safely extracts float64 from interface{}
func getFloat64FromInterface(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}