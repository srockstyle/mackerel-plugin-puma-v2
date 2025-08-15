# mackerel-plugin-puma-v2
[![Go Report Card](https://goreportcard.com/badge/github.com/srockstyle/mackerel-plugin-puma-v2)](https://goreportcard.com/report/github.com/srockstyle/mackerel-plugin-puma-v2)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A modern Mackerel plugin for monitoring Puma web server, built from scratch for Puma 6.x compatibility. This is a complete rewrite of the original mackerel-plugin-puma with enhanced features and better performance.

## Features

- **Puma 6.x Compatible**: Fully supports Puma 6.6.1 and newer versions
- **Comprehensive Metrics**: Workers, threads, request backlogs, memory usage, and more
- **Multiple Connection Methods**: HTTP, HTTPS, and Unix domain socket support
- **Flexible Configuration**: Single mode and cluster mode support
- **Extended Metrics**: Request count, uptime, thread utilization (Puma 6.x)
- **Backward Compatible**: Works with Puma 4.x and 5.x
- **Modern Go Features**: Built with Go 1.24, utilizing range over integers, slices package, and iterators

## Requirements

- Puma 4.0 or later (optimized for 6.6.1+)
- Puma control server enabled
- Go 1.24 or later (for building from source)

## Installation

### Via mkr

You can install this plugin via `mkr plugin install`:

```console
$ mkr plugin install mackerel-plugin-puma-v2
```

### From Source

```console
$ go install github.com/srockstyle/mackerel-plugin-puma-v2@latest
```

### Manual Download

Download the latest release from the [releases page](https://github.com/srockstyle/mackerel-plugin-puma-v2/releases).

## Usage

```
Usage of mackerel-plugin-puma-v2:
  -socket string
        Path to Puma control socket (default: /tmp/puma.sock)
  -metric-key-prefix string
        Metric key prefix (default "puma")
  -tempfile string
        Temp file name for storing state
  -extended
        Collect extended metrics (memory, GC, thread utilization, etc)
```

## Configuration

### Basic Unix Socket Connection (Recommended)

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock"
```

### Custom Socket Path

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/var/run/puma/pumactl.sock"
```

### With Extended Metrics

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -extended"
```

### Custom Metric Prefix

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -metric-key-prefix=myapp_puma"
```

### Environment Variables

You can also use environment variables:

```bash
export PUMA_SOCKET=/var/run/puma/pumactl.sock
/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

## Metrics

### Core Metrics

#### Worker Metrics
- `puma.workers` - Number of worker processes
- `puma.booted_workers` - Number of booted workers
- `puma.old_workers` - Number of old workers (during phased restart)
- `puma.phase` - Current phase number (during phased restart)

#### Thread Metrics
- `puma.backlog` - Request backlog
- `puma.running` - Running threads
- `puma.pool_capacity` - Thread pool capacity
- `puma.max_threads` - Maximum threads configured
- `puma.thread_utilization` - Thread utilization percentage (Puma 6.x)

#### Request Metrics (Puma 6.x)
- `puma.requests_count` - Total number of requests processed (counter)
- `puma.uptime` - Server uptime in seconds

### Extended Metrics (with -extended flag)

#### Memory Metrics
- `puma.memory.alloc` - Allocated memory (MB)
- `puma.memory.sys` - System memory (MB)
- `puma.memory.heap_inuse` - Heap memory in use (MB)

#### GC Metrics
- `puma.gc.num_gc` - Number of GC runs (counter)
- `puma.ruby.gc.count` - Ruby GC count (if available)
- `puma.ruby.gc.heap_used` - Ruby heap slots used
- `puma.ruby.gc.heap_length` - Ruby heap slots total

#### Go Runtime Metrics
- `puma.go.goroutines` - Number of goroutines

## Puma Configuration

To enable the control server in your Puma configuration:

### config/puma.rb

```ruby
# Recommended: Unix socket (more secure and efficient)
activate_control_app 'unix:///tmp/puma.sock'

# Alternative: With custom path
activate_control_app 'unix:///var/run/puma/pumactl.sock'

# Note: Authentication token is not currently supported in v2
# This feature is planned for a future release
```

## Version Compatibility

| Puma Version | Plugin Support | Notes |
|--------------|----------------|-------|
| 6.x          | ✅ Full        | All features including extended metrics |
| 5.x          | ✅ Full        | All features except some 6.x specific metrics |
| 4.x          | ✅ Basic       | Basic metrics only |
| 3.x          | ⚠️ Limited     | May work but not tested |

## Troubleshooting

### Connection Refused

Ensure Puma control server is enabled:
```bash
$ curl http://127.0.0.1:9293/stats?token=your-secret-token
```

### Authentication Error

Check your token configuration matches between Puma and the plugin.

### No Metrics

1. Verify Puma is running: `ps aux | grep puma`
2. Check control server is accessible
3. Enable debug logging: `MACKEREL_PLUGIN_DEBUG=1 mackerel-plugin-puma-v2 ...`

### Unix Socket Permission

Ensure the mackerel-agent user has read permission for the socket file:
```bash
$ ls -l /path/to/pumactl.sock
```

## Development

### Building

```bash
$ make build
```

### Testing

```bash
$ make test
```

### Running locally

```bash
$ go run main.go -host=localhost -port=9293
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

[srockstyle](https://github.com/srockstyle)

## Why v2?

The original mackerel-plugin-puma hasn't been updated since 2017 and lacks support for modern Puma versions. This v2 implementation:

- Complete rewrite with clean architecture
- Full Puma 6.x support with new metrics
- Better error handling and retry logic
- Comprehensive test coverage
- Active maintenance and updates

## Acknowledgments

- Original mackerel-plugin-puma by [rmanzoku](https://github.com/rmanzoku)
- Mackerel team for the plugin framework
- Puma team for the excellent web server
