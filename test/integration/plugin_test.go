// +build integration

package integration_test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/application"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

func TestPumaPlugin_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a mock Puma stats server
	mux := http.NewServeMux()
	
	// Mock /stats endpoint
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := infrastructure.PumaStats{
			Workers:       4,
			Phase:         0,
			BootedWorkers: 4,
			OldWorkers:    0,
			RequestsCount: int64Ptr(12345),
			Uptime:        intPtr(3600),
			WorkerStatus: []infrastructure.WorkerStatus{
				{
					PID:    1234,
					Index:  0,
					Phase:  0,
					Booted: true,
					LastStatus: infrastructure.LastStatus{
						Backlog:      0,
						Running:      5,
						PoolCapacity: 16,
						MaxThreads:   16,
					},
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Mock /gc-stats endpoint
	mux.HandleFunc("/gc-stats", func(w http.ResponseWriter, r *http.Request) {
		gcStats := map[string]interface{}{
			"count":       100,
			"heap_used":   1000,
			"heap_length": 2000,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(gcStats)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// Skip HTTP endpoint test for now as it requires socket implementation
	t.Run("HTTP endpoint", func(t *testing.T) {
		t.Skip("HTTP endpoint test requires socket implementation")
	})

	// Test with Unix socket
	t.Run("Unix socket", func(t *testing.T) {
		// Create temporary socket
		tmpDir, err := os.MkdirTemp("", "puma-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		socketPath := tmpDir + "/puma.sock"
		
		// Create Unix socket listener
		listener, err := net.Listen("unix", socketPath)
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()

		// Start mock server on Unix socket
		go http.Serve(listener, mux)

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		config := application.DefaultConfig()
		config.SocketPath = socketPath
		
		logger := log.New(os.Stderr, "[test] ", log.LstdFlags)
		collector := application.NewMetricsCollector(config, logger)
		ctx := context.Background()
		
		collection, err := collector.Collect(ctx)
		if err != nil {
			t.Fatalf("Failed to collect metrics via Unix socket: %v", err)
		}

		// Verify metrics were collected
		metrics := collection.All()
		if len(metrics) == 0 {
			t.Error("No metrics collected via Unix socket")
		}
	})
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}