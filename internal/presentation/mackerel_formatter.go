package presentation

import (
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// MackerelPlugin implements the Mackerel plugin interface
type MackerelPlugin struct {
	prefix string
}

// NewMackerelPlugin creates a new Mackerel plugin
func NewMackerelPlugin(prefix string) *MackerelPlugin {
	return &MackerelPlugin{
		prefix: prefix,
	}
}

// GraphDefinition returns graph definitions for Mackerel
func (p *MackerelPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"workers": {
			Label: "Puma Workers",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "workers", Label: "Workers"},
				{Name: "booted_workers", Label: "Booted Workers"},
				{Name: "old_workers", Label: "Old Workers"},
			},
		},
		"threads": {
			Label: "Puma Threads",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "running", Label: "Running"},
				{Name: "pool_capacity", Label: "Pool Capacity"},
				{Name: "max_threads", Label: "Max Threads"},
			},
		},
		"backlog": {
			Label: "Puma Backlog",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "backlog", Label: "Backlog"},
			},
		},
		"phase": {
			Label: "Puma Phase",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "phase", Label: "Phase"},
			},
		},
		"requests": {
			Label: "Puma Requests",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "requests_count", Label: "Requests Count", Diff: true},
			},
		},
		"uptime": {
			Label: "Puma Uptime",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "uptime", Label: "Uptime"},
			},
		},
		"memory": {
			Label: "Memory Usage",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "memory.alloc", Label: "Allocated"},
				{Name: "memory.sys", Label: "System"},
				{Name: "memory.heap_inuse", Label: "Heap In Use"},
			},
		},
		"gc": {
			Label: "Garbage Collection",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "gc.num_gc", Label: "GC Count", Diff: true},
				{Name: "ruby.gc.count", Label: "Ruby GC Count", Diff: true},
			},
		},
		"ruby_heap": {
			Label: "Ruby Heap",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{Name: "ruby.gc.heap_used", Label: "Heap Used"},
				{Name: "ruby.gc.heap_length", Label: "Heap Length"},
			},
		},
		"thread_utilization": {
			Label: "Thread Utilization",
			Unit:  mp.UnitPercentage,
			Metrics: []mp.Metrics{
				{Name: "thread_utilization", Label: "Utilization %"},
			},
		},
	}
}

// FormatMetrics formats metrics for Mackerel output
func (p *MackerelPlugin) FormatMetrics(collection *domain.MetricCollection) map[string]float64 {
	result := make(map[string]float64)

	for _, metric := range collection.All() {
		key := p.buildMetricKey(metric.Name)
		result[key] = metric.Value
	}

	return result
}

// buildMetricKey builds the full metric key with prefix
func (p *MackerelPlugin) buildMetricKey(name string) string {
	return name
}