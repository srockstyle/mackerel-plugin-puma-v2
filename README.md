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

## Requirements

- Puma 4.0 or later (optimized for 6.6.1+)
- Puma control server enabled
- Go 1.20 or later (for building from source)

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
  -host string
        The bind host for the control server (default "127.0.0.1")
  -port string
        The bind port for the control server (default "9293")
  -sock string
        Unix domain socket path (overrides host/port)
  -scheme string
        URL scheme (http or https) (default "http")
  -token string
        Authentication token for the control server
  -single
        Monitor Puma in single mode (non-clustered)
  -with-gc
        Include GC statistics (may not be available in Puma 6.x)
  -metric-key-prefix string
        Metric key prefix (default "puma")
  -tempfile string
        Temp file path for storing state
  -timeout duration
        Request timeout (default 10s)
  -retry int
        Number of retry attempts (default 3)
```

## Configuration

### Basic HTTP Connection

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2"
```

### With Authentication Token

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -token=your-secret-token"
```

### Unix Domain Socket

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -sock=/path/to/pumactl.sock -token=your-secret-token"
```

### Single Mode (Non-clustered Puma)

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -single -token=your-secret-token"
```

### HTTPS Connection

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -scheme=https -host=puma.example.com -port=9293"
```

### Advanced Configuration

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -token=secret -timeout=30s -retry=5 -metric-key-prefix=myapp_puma"
```

## Metrics

### Worker Metrics (Cluster Mode)

- `puma.workers` - Number of worker processes
- `puma.booted_workers` - Number of booted workers
- `puma.old_workers` - Number of old workers (during phased restart)
- `puma.phase` - Current phase number (during phased restart)

### Thread Metrics

#### Cluster Mode
- `puma.worker_backlog.worker{N}.backlog` - Request backlog per worker
- `puma.worker_running.worker{N}.running` - Running threads per worker
- `puma.worker_pool_capacity.worker{N}.pool_capacity` - Thread pool capacity per worker
- `puma.worker_thread_utilization.worker{N}.utilization` - Thread utilization percentage per worker (Puma 6.x)

#### Single Mode
- `puma.backlog` - Request backlog
- `puma.running` - Running threads
- `puma.pool_capacity` - Thread pool capacity
- `puma.max_threads` - Maximum threads configured

### Extended Metrics (Puma 6.x)

- `puma.requests_count` - Total number of requests processed
- `puma.uptime` - Server uptime in seconds

### GC Metrics (Optional, may not be available in Puma 6.x)

When using `-with-gc` flag:
- `puma.gc.count.total` - Total GC count
- `puma.gc.count.minor` - Minor GC count
- `puma.gc.count.major` - Major GC count
- `puma.gc.heap_slot.*` - Heap slot statistics
- `puma.gc.old_objects.*` - Old object statistics

## Puma Configuration

To enable the control server in your Puma configuration:

### config/puma.rb

```ruby
# Enable control server
activate_control_app 'tcp://127.0.0.1:9293', { auth_token: 'your-secret-token' }

# Or with Unix socket
activate_control_app 'unix:///path/to/pumactl.sock', { auth_token: 'your-secret-token' }

# For development without auth
activate_control_app 'tcp://127.0.0.1:9293', { no_token: true }
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
