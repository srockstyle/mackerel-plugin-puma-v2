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
}

// MetricDefinition defines a metric's properties
type MetricDefinition struct {
	Name  string
	Label string
	Type  MetricType
	Unit  string
}