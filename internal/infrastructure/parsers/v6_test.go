package parsers_test

import (
	"testing"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure/parsers"
)

func TestV6Parser_Parse(t *testing.T) {
	parser := parsers.NewV6Parser()

	t.Run("cluster mode metrics", func(t *testing.T) {
		stats := &infrastructure.PumaStats{
			Workers:       4,
			Phase:         0,
			BootedWorkers: 4,
			OldWorkers:    0,
			RequestsCount: int64Ptr(12345),
			Uptime:        intPtr(3600),
			WorkerStatus: []infrastructure.WorkerStatus{
				{
					PID:   1234,
					Index: 0,
					Phase: 0,
					Booted: true,
					LastStatus: infrastructure.LastStatus{
						Backlog:      0,
						Running:      5,
						PoolCapacity: 16,
						MaxThreads:   16,
					},
				},
				{
					PID:   1235,
					Index: 1,
					Phase: 0,
					Booted: true,
					LastStatus: infrastructure.LastStatus{
						Backlog:      2,
						Running:      3,
						PoolCapacity: 16,
						MaxThreads:   16,
					},
				},
			},
		}

		collection, err := parser.Parse(stats)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Check basic metrics
		checkMetric(t, collection, "workers", 4.0)
		checkMetric(t, collection, "booted_workers", 4.0)
		checkMetric(t, collection, "old_workers", 0.0)
		checkMetric(t, collection, "phase", 0.0)

		// Check Puma 6.x specific metrics
		checkMetric(t, collection, "requests_count", 12345.0)
		checkMetric(t, collection, "uptime", 3600.0)

		// Check aggregated worker metrics
		checkMetric(t, collection, "backlog", 2.0)
		checkMetric(t, collection, "running", 8.0)
		checkMetric(t, collection, "pool_capacity", 32.0)
		checkMetric(t, collection, "max_threads", 32.0)

		// Check thread utilization
		checkMetric(t, collection, "thread_utilization", 25.0) // 8/32 * 100
	})

	t.Run("single mode metrics", func(t *testing.T) {
		stats := &infrastructure.PumaStats{
			Workers:       1,
			Phase:         0,
			BootedWorkers: 1,
			OldWorkers:    0,
			Backlog:       intPtr(5),
			Running:       intPtr(10),
			PoolCapacity:  intPtr(20),
			MaxThreads:    intPtr(20),
		}

		collection, err := parser.Parse(stats)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Check single mode metrics
		checkMetric(t, collection, "backlog", 5.0)
		checkMetric(t, collection, "running", 10.0)
		checkMetric(t, collection, "pool_capacity", 20.0)
		checkMetric(t, collection, "max_threads", 20.0)
	})

	t.Run("no thread utilization when pool capacity is zero", func(t *testing.T) {
		stats := &infrastructure.PumaStats{
			Workers: 1,
			WorkerStatus: []infrastructure.WorkerStatus{
				{
					LastStatus: infrastructure.LastStatus{
						Running:      5,
						PoolCapacity: 0,
					},
				},
			},
		}

		collection, err := parser.Parse(stats)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Thread utilization should not be calculated
		metric := findMetric(collection, "thread_utilization")
		if metric != nil {
			t.Error("thread_utilization should not be calculated when pool capacity is 0")
		}
	})
}

func checkMetric(t *testing.T, collection *domain.MetricCollection, name string, expectedValue float64) {
	t.Helper()
	metric := findMetric(collection, name)
	if metric == nil {
		t.Errorf("Metric %s not found", name)
		return
	}
	if metric.Value != expectedValue {
		t.Errorf("Metric %s: got value %f, want %f", name, metric.Value, expectedValue)
	}
}

func findMetric(collection *domain.MetricCollection, name string) *domain.Metric {
	for _, m := range collection.All() {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}