package infrastructure_test

import (
	"context"
	"testing"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

// MockPumaClient is a mock implementation of PumaClient
type MockPumaClient struct {
	stats   *infrastructure.PumaStats
	gcStats map[string]interface{}
	err     error
}

func (m *MockPumaClient) GetStats(ctx context.Context) (*infrastructure.PumaStats, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stats, nil
}

func (m *MockPumaClient) GetGCStats(ctx context.Context) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.gcStats, nil
}

func TestVersionDetector_DetectVersion(t *testing.T) {
	tests := []struct {
		name    string
		stats   *infrastructure.PumaStats
		want    string
		wantErr bool
	}{
		{
			name: "Puma 6.x with requests_count",
			stats: &infrastructure.PumaStats{
				Workers:       4,
				RequestsCount: int64Ptr(12345),
			},
			want:    "6.x",
			wantErr: false,
		},
		{
			name: "Puma 5.x with detailed worker status",
			stats: &infrastructure.PumaStats{
				Workers: 4,
				WorkerStatus: []infrastructure.WorkerStatus{
					{
						PID:   1234,
						Index: 0,
						LastStatus: infrastructure.LastStatus{
							MaxThreads: 16,
						},
					},
				},
			},
			want:    "5.x",
			wantErr: false,
		},
		{
			name: "Puma 4.x default",
			stats: &infrastructure.PumaStats{
				Workers: 4,
			},
			want:    "4.x",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockPumaClient{
				stats: tt.stats,
			}
			detector := infrastructure.NewVersionDetector(mockClient)

			got, err := detector.DetectVersion(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPuma6OrLater(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"6.x", true},
		{"6.0.0", true},
		{"6.6.1", true},
		{"7.0.0", true},
		{"5.x", false},
		{"5.6.1", false},
		{"4.x", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := infrastructure.IsPuma6OrLater(tt.version); got != tt.want {
				t.Errorf("IsPuma6OrLater(%s) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
