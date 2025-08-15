package domain

// MetricDefinitions holds all metric definitions for Puma
var MetricDefinitions = map[string]MetricDefinition{
	// Worker metrics
	"workers": {
		Name:  "workers",
		Label: "Workers",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"booted_workers": {
		Name:  "booted_workers",
		Label: "Booted Workers",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"old_workers": {
		Name:  "old_workers",
		Label: "Old Workers",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"phase": {
		Name:  "phase",
		Label: "Phase",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},

	// Thread pool metrics
	"backlog": {
		Name:  "backlog",
		Label: "Backlog",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"running": {
		Name:  "running",
		Label: "Running Threads",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"pool_capacity": {
		Name:  "pool_capacity",
		Label: "Pool Capacity",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"max_threads": {
		Name:  "max_threads",
		Label: "Max Threads",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},

	// Puma 6.x specific metrics
	"requests_count": {
		Name:  "requests_count",
		Label: "Requests Count",
		Type:  MetricTypeCounter,
		Unit:  "integer",
	},
	"uptime": {
		Name:  "uptime",
		Label: "Uptime",
		Type:  MetricTypeGauge,
		Unit:  "seconds",
	},

	// Memory metrics
	"memory.alloc": {
		Name:  "memory.alloc",
		Label: "Memory Allocated",
		Type:  MetricTypeGauge,
		Unit:  "megabytes",
	},
	"memory.sys": {
		Name:  "memory.sys",
		Label: "Memory System",
		Type:  MetricTypeGauge,
		Unit:  "megabytes",
	},
	"memory.heap_inuse": {
		Name:  "memory.heap_inuse",
		Label: "Heap In Use",
		Type:  MetricTypeGauge,
		Unit:  "megabytes",
	},

	// GC metrics
	"gc.num_gc": {
		Name:  "gc.num_gc",
		Label: "GC Count",
		Type:  MetricTypeCounter,
		Unit:  "integer",
	},
	"ruby.gc.count": {
		Name:  "ruby.gc.count",
		Label: "Ruby GC Count",
		Type:  MetricTypeCounter,
		Unit:  "integer",
	},
	"ruby.gc.heap_used": {
		Name:  "ruby.gc.heap_used",
		Label: "Ruby Heap Used",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.heap_length": {
		Name:  "ruby.gc.heap_length",
		Label: "Ruby Heap Length",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},

	// Thread utilization metrics
	"thread_utilization": {
		Name:  "thread_utilization",
		Label: "Thread Utilization",
		Type:  MetricTypeGauge,
		Unit:  "percentage",
	},

	// Go runtime metrics
	"go.goroutines": {
		Name:  "go.goroutines",
		Label: "Goroutines",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},

	// Detailed Ruby GC metrics (from lib/gc.go)
	"ruby.gc.minor_count": {
		Name:  "ruby.gc.minor_count",
		Label: "Ruby Minor GC Count",
		Type:  MetricTypeCounter,
		Unit:  "integer",
	},
	"ruby.gc.major_count": {
		Name:  "ruby.gc.major_count",
		Label: "Ruby Major GC Count",
		Type:  MetricTypeCounter,
		Unit:  "integer",
	},
	"ruby.gc.heap_available_slots": {
		Name:  "ruby.gc.heap_available_slots",
		Label: "Ruby Heap Available Slots",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.heap_live_slots": {
		Name:  "ruby.gc.heap_live_slots",
		Label: "Ruby Heap Live Slots",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.heap_free_slots": {
		Name:  "ruby.gc.heap_free_slots",
		Label: "Ruby Heap Free Slots",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.heap_final_slots": {
		Name:  "ruby.gc.heap_final_slots",
		Label: "Ruby Heap Final Slots",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.heap_marked_slots": {
		Name:  "ruby.gc.heap_marked_slots",
		Label: "Ruby Heap Marked Slots",
		Type:  MetricTypeGauge,
		Unit:  "slots",
	},
	"ruby.gc.old_objects": {
		Name:  "ruby.gc.old_objects",
		Label: "Ruby Old Objects",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"ruby.gc.old_objects_limit": {
		Name:  "ruby.gc.old_objects_limit",
		Label: "Ruby Old Objects Limit",
		Type:  MetricTypeGauge,
		Unit:  "integer",
	},
	"ruby.gc.oldmalloc_bytes": {
		Name:  "ruby.gc.oldmalloc_bytes",
		Label: "Ruby Old Malloc Bytes",
		Type:  MetricTypeGauge,
		Unit:  "bytes",
	},
	"ruby.gc.oldmalloc_limit": {
		Name:  "ruby.gc.oldmalloc_limit",
		Label: "Ruby Old Malloc Limit",
		Type:  MetricTypeGauge,
		Unit:  "bytes",
	},
}

// MetricDefinition defines a metric's properties
type MetricDefinition struct {
	Name  string
	Label string
	Type  MetricType
	Unit  string
}