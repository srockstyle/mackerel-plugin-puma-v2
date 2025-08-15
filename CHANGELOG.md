# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2024-08-16

### Added
- Complete rewrite with Clean Architecture design
- Full support for Puma 6.x including new metrics
- Automatic version detection for Puma 4.x, 5.x, and 6.x
- Extended metrics collection with `-extended` flag
  - Memory usage metrics (alloc, sys, heap)
  - GC statistics (Go and Ruby)
  - Thread utilization percentage
  - Goroutine count
- Unix socket support as primary connection method
- Retry mechanism with configurable attempts
- Timeout configuration for all operations
- Environment variable support (PUMA_SOCKET)
- Comprehensive test suite with 80%+ coverage
- Go 1.24 support with modern language features
  - Range over integers
  - Slices package utilization
  - Iterator pattern support

### Changed
- Module name to `mackerel-plugin-puma-v2`
- Primary connection method from HTTP to Unix socket
- Default socket path to `/tmp/puma.sock`
- Improved error handling and logging
- Modular architecture for better maintainability

### Fixed
- Unix socket communication buffer overflow issues
- Resource leaks in error conditions
- HTTP response parsing reliability
- Error handling in retry logic

### Security
- Unix socket provides better security than HTTP endpoint
- No authentication token required for Unix socket

### Removed
- HTTP/HTTPS as primary connection method (Unix socket recommended)
- Authentication token support (planned for future release)
- Direct HTTP configuration options