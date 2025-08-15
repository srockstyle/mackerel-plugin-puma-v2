package domain_test

import (
	"testing"
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
)

func TestMetric_Validate(t *testing.T) {
	tests := []struct {
		name    string
		metric  domain.Metric
		wantErr bool
	}{
		{
			name: "valid gauge metric",
			metric: domain.Metric{
				Name:  "test_metric",
				Value: 42.0,
				Type:  domain.MetricTypeGauge,
			},
			wantErr: false,
		},
		{
			name: "valid counter metric",
			metric: domain.Metric{
				Name:  "test_counter",
				Value: 100.0,
				Type:  domain.MetricTypeCounter,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			metric: domain.Metric{
				Name:  "",
				Value: 42.0,
				Type:  domain.MetricTypeGauge,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			metric: domain.Metric{
				Name:  "test_metric",
				Value: 42.0,
				Type:  "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metric.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Metric.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricCollection(t *testing.T) {
	t.Run("add and retrieve metrics", func(t *testing.T) {
		collection := domain.NewMetricCollection()

		metric1 := domain.Metric{
			Name:      "test1",
			Value:     10.0,
			Type:      domain.MetricTypeGauge,
			Timestamp: time.Now(),
		}

		metric2 := domain.Metric{
			Name:      "test2",
			Value:     20.0,
			Type:      domain.MetricTypeCounter,
			Timestamp: time.Now(),
		}

		// Add metrics
		if err := collection.Add(metric1); err != nil {
			t.Fatalf("Failed to add metric1: %v", err)
		}

		if err := collection.Add(metric2); err != nil {
			t.Fatalf("Failed to add metric2: %v", err)
		}

		// Check all metrics
		all := collection.All()
		if len(all) != 2 {
			t.Errorf("Expected 2 metrics, got %d", len(all))
		}
	})

	t.Run("filter metrics", func(t *testing.T) {
		collection := domain.NewMetricCollection()

		for i := range 5 {
			metric := domain.Metric{
				Name:      "test",
				Value:     float64(i),
				Type:      domain.MetricTypeGauge,
				Timestamp: time.Now(),
			}
			collection.Add(metric)
		}

		// Filter metrics with value > 2
		filtered := collection.Filter(func(m domain.Metric) bool {
			return m.Value > 2
		})

		if len(filtered) != 2 {
			t.Errorf("Expected 2 filtered metrics, got %d", len(filtered))
		}
	})

	t.Run("iterator", func(t *testing.T) {
		collection := domain.NewMetricCollection()

		for i := range 3 {
			metric := domain.Metric{
				Name:      "test",
				Value:     float64(i),
				Type:      domain.MetricTypeGauge,
				Timestamp: time.Now(),
			}
			collection.Add(metric)
		}

		count := 0
		for i, m := range collection.Iter() {
			if i != count {
				t.Errorf("Expected index %d, got %d", count, i)
			}
			if m.Value != float64(count) {
				t.Errorf("Expected value %f, got %f", float64(count), m.Value)
			}
			count++
		}

		if count != 3 {
			t.Errorf("Expected to iterate 3 times, got %d", count)
		}
	})

	t.Run("clear", func(t *testing.T) {
		collection := domain.NewMetricCollection()

		// Add a metric
		collection.Add(domain.Metric{
			Name:  "test",
			Value: 42.0,
			Type:  domain.MetricTypeGauge,
		})

		// Clear collection
		collection.Clear()

		// Check that it's empty
		if len(collection.All()) != 0 {
			t.Error("Collection should be empty after Clear()")
		}
	})
}