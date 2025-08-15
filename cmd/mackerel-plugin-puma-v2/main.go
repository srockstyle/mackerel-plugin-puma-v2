package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/application"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/presentation"
)

// PumaPlugin represents the Puma plugin
type PumaPlugin struct {
	Socket   string
	Prefix   string
	collector *application.MetricsCollector
	formatter *presentation.MackerelPlugin
}

// MetricKeyPrefix returns the metric key prefix
func (p *PumaPlugin) MetricKeyPrefix() string {
	return p.Prefix
}

// GraphDefinition returns graph definitions
func (p *PumaPlugin) GraphDefinition() map[string]mp.Graphs {
	return p.formatter.GraphDefinition()
}

// FetchMetrics fetches metrics from Puma
func (p *PumaPlugin) FetchMetrics() (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection, err := p.collector.Collect(ctx)
	if err != nil {
		return nil, err
	}

	return p.formatter.FormatMetrics(collection), nil
}

func main() {
	// Command line options
	optSocket := flag.String("socket", "", "Path to Puma control socket")
	// optScheme := flag.String("scheme", "http", "Scheme (not used with socket)")
	// optHost := flag.String("host", "localhost", "Hostname (not used with socket)")
	// optPort := flag.String("port", "9293", "Port (not used with socket)")
	optPrefix := flag.String("metric-key-prefix", "puma", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	// Setup logger
	logger := log.New(os.Stderr, "[mackerel-plugin-puma] ", log.LstdFlags)

	// Create config
	config := application.DefaultConfig()
	config.MetricPrefix = *optPrefix

	// Socket takes precedence
	if *optSocket != "" {
		config.SocketPath = *optSocket
	} else {
		// Check environment variable
		if envSocket := os.Getenv("PUMA_SOCKET"); envSocket != "" {
			config.SocketPath = envSocket
		}
	}

	// Default to Unix socket
	if config.SocketPath == "" {
		config.SocketPath = "/tmp/puma.sock"
	}

	// Validate config
	if err := config.Validate(); err != nil {
		logger.Fatalf("Invalid configuration: %v", err)
	}

	// Create plugin
	collector := application.NewMetricsCollector(config, logger)
	formatter := presentation.NewMackerelPlugin(config.MetricPrefix)

	plugin := &PumaPlugin{
		Socket:    config.SocketPath,
		Prefix:    config.MetricPrefix,
		collector: collector,
		formatter: formatter,
	}

	// Run plugin
	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile
	helper.Run()
}