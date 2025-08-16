package infrastructure

import (
	"context"
	"regexp"
	"strings"
)

// VersionDetector detects Puma version from stats
type VersionDetector struct {
	client PumaClient
}

// NewVersionDetector creates a new version detector
func NewVersionDetector(client PumaClient) *VersionDetector {
	return &VersionDetector{
		client: client,
	}
}

// DetectVersion attempts to detect the Puma version
func (d *VersionDetector) DetectVersion(ctx context.Context) (string, error) {
	// Try to get stats first
	stats, err := d.client.GetStats(ctx)
	if err != nil {
		return "", err
	}

	// Check for version-specific fields
	if stats.RequestsCount != nil {
		// Puma 6.x has requests_count field
		return "6.x", nil
	}

	if len(stats.WorkerStatus) > 0 {
		// Check worker status for more detailed version info
		for _, worker := range stats.WorkerStatus {
			// Puma 5.x has more detailed worker status
			if worker.LastStatus.MaxThreads > 0 {
				return "5.x", nil
			}
		}
	}

	// Default to 4.x for older versions
	return "4.x", nil
}

// GetVersionFromGCStats tries to extract version from gc-stats endpoint
func (d *VersionDetector) GetVersionFromGCStats(ctx context.Context) (string, error) {
	gcStats, err := d.client.GetGCStats(ctx)
	if err != nil {
		// GC stats might not be available in newer versions
		return "", err
	}

	// Look for version string in GC stats
	if version, ok := gcStats["version"].(string); ok {
		return extractVersionNumber(version), nil
	}

	return "", nil
}

// extractVersionNumber extracts the version number from a version string
func extractVersionNumber(versionStr string) string {
	// Match patterns like "5.0.0", "6.6.1"
	re := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionStr)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try to match major.minor
	re = regexp.MustCompile(`(\d+\.\d+)`)
	matches = re.FindStringSubmatch(versionStr)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try to match major version only
	re = regexp.MustCompile(`(\d+)`)
	matches = re.FindStringSubmatch(versionStr)
	if len(matches) > 1 {
		return matches[1] + ".x"
	}

	return versionStr
}

// IsPuma6OrLater checks if the detected version is Puma 6.x or later
func IsPuma6OrLater(version string) bool {
	return strings.HasPrefix(version, "6.") ||
		strings.HasPrefix(version, "7.") ||
		strings.HasPrefix(version, "8.") ||
		strings.HasPrefix(version, "9.")
}
