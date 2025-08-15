package application

import (
	"fmt"
	"time"
)

// Config holds application configuration
type Config struct {
	// Connection settings
	Host       string
	Port       string
	SocketPath string
	Scheme     string

	// Authentication
	Token string

	// Behavior settings
	SingleMode     bool
	WithGC         bool
	MetricPrefix   string

	// Performance settings
	Timeout       time.Duration
	RetryCount    int
	RetryInterval time.Duration
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Host:          "127.0.0.1",
		Port:          "9293",
		Scheme:        "http",
		MetricPrefix:  "puma",
		Timeout:       10 * time.Second,
		RetryCount:    3,
		RetryInterval: 1 * time.Second,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SocketPath == "" {
		if c.Host == "" {
			return fmt.Errorf("either socket path or host must be specified")
		}
		if c.Port == "" {
			return fmt.Errorf("port must be specified when using host")
		}
		if c.Scheme != "http" && c.Scheme != "https" {
			return fmt.Errorf("scheme must be http or https, got %s", c.Scheme)
		}
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.RetryCount < 0 {
		return fmt.Errorf("retry count must be non-negative")
	}

	return nil
}

// GetBaseURL returns the base URL for API requests
func (c *Config) GetBaseURL() string {
	if c.SocketPath != "" {
		return "http://localhost"
	}
	return fmt.Sprintf("%s://%s:%s", c.Scheme, c.Host, c.Port)
}