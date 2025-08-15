# Migration Guide: mackerel-plugin-puma to mackerel-plugin-puma-v2

## Overview

This document explains the migration from the original mackerel-plugin-puma to mackerel-plugin-puma-v2.

## Why v2?

The original mackerel-plugin-puma hasn't been updated since 2017 and has the following issues:

- No support for Puma 6.x features
- Limited error handling
- Outdated dependencies
- Buffer overflow issues with Unix socket communication
- No active maintenance

mackerel-plugin-puma-v2 is a complete rewrite that addresses these issues with:

- Full Puma 6.x compatibility with automatic version detection
- Clean Architecture design
- Comprehensive error handling with retry logic
- Modern Go 1.24 features
- Extended metrics collection
- Active maintenance

## Major Changes

### Primary Connection Method
- **Old**: HTTP/HTTPS connection
- **New**: Unix domain socket (recommended for security and performance)

### Architecture
- **Old**: Single file implementation
- **New**: Clean Architecture with modular design

### Go Version
- **Old**: Go 1.11+
- **New**: Go 1.24+ (uses modern language features)

## Installation Changes

### Old (mackerel-plugin-puma)
```bash
mkr plugin install mackerel-plugin-puma
```

### New (mackerel-plugin-puma-v2)
```bash
mkr plugin install mackerel-plugin-puma-v2
# or
go install github.com/srockstyle/mackerel-plugin-puma-v2@latest
```

## Configuration Changes

### Binary Name
- Old: `mackerel-plugin-puma`
- New: `mackerel-plugin-puma-v2`

### Command Line Options

#### Removed Options
- `-host string` - No longer supported
- `-port string` - No longer supported
- `-token string` - Authentication not required for Unix socket
- `-scheme string` - Unix socket doesn't use URL scheme
- `-with-gc` - Now included in `-extended` flag

#### New/Changed Options
- `-socket string` - Path to Puma control socket (default "/tmp/puma.sock")
- `-extended` - Collect extended metrics including GC stats
- Environment variable support: `PUMA_SOCKET`

### Configuration Examples

#### Old Configuration (HTTP)
```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma -host=localhost -port=9293 -token=secret"
```

#### New Configuration (Unix Socket)
```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock"
```

#### With Extended Metrics
```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -extended"
```

## Puma Configuration Changes

### Enable Unix Socket Control
Update your `config/puma.rb`:

```ruby
# Old (HTTP)
activate_control_app 'tcp://127.0.0.1:9293', { auth_token: 'secret' }

# New (Unix Socket) - Recommended
activate_control_app 'unix:///tmp/puma.sock'
```

## New Metrics in v2

### Core Metrics (All Versions)
- Same as v1 (workers, threads, backlog, etc.)

### Puma 6.x Specific
- `puma.requests_count` - Total requests processed
- `puma.uptime` - Server uptime in seconds
- `puma.thread_utilization` - Thread utilization percentage

### Extended Metrics (-extended flag)

#### Memory Metrics
- `puma.memory.alloc` - Allocated memory (MB)
- `puma.memory.sys` - System memory (MB)
- `puma.memory.heap_inuse` - Heap memory in use (MB)

#### Detailed Ruby GC Metrics
- `puma.ruby.gc.minor_count` - Minor GC count
- `puma.ruby.gc.major_count` - Major GC count
- `puma.ruby.gc.heap_available_slots` - Available heap slots
- `puma.ruby.gc.heap_live_slots` - Live heap slots
- `puma.ruby.gc.heap_free_slots` - Free heap slots
- `puma.ruby.gc.old_objects` - Old generation objects
- `puma.ruby.gc.oldmalloc_bytes` - Old malloc bytes

#### Go Runtime Metrics
- `puma.go.goroutines` - Number of goroutines

## Migration Steps

1. **Update Puma Configuration**
   ```ruby
   # Add to config/puma.rb
   activate_control_app 'unix:///tmp/puma.sock'
   ```

2. **Test Connection**
   ```bash
   # Test with Unix socket
   mackerel-plugin-puma-v2 -socket=/tmp/puma.sock
   ```

3. **Update Mackerel Agent Configuration**
   ```toml
   [plugin.metrics.puma]
   command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock"
   ```

4. **Enable Extended Metrics (Optional)**
   ```toml
   [plugin.metrics.puma]
   command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -extended"
   ```

5. **Remove Old Plugin**
   ```bash
   rm /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma
   ```

## Breaking Changes

### Connection Method
- HTTP/HTTPS connection support removed in favor of Unix socket
- Authentication token not supported (Unix socket provides file system security)

### Metric Names
- All metric names remain backward compatible

### GC Stats
- `-with-gc` flag replaced by `-extended` flag
- GC stats now include more detailed Ruby GC metrics

## Troubleshooting

### Permission Denied
Ensure mackerel-agent user has read permission for the socket:
```bash
ls -l /tmp/puma.sock
# srw-rw-rw- 1 puma puma 0 Aug 16 10:00 /tmp/puma.sock
```

### Socket Not Found
Verify Puma is running and control server is enabled:
```bash
ps aux | grep puma
# Check if socket file exists
ls -la /tmp/puma.sock
```

### No Metrics Collected
Enable debug mode:
```bash
MACKEREL_PLUGIN_DEBUG=1 mackerel-plugin-puma-v2 -socket=/tmp/puma.sock
```

## Rollback Plan

If you need to rollback:

1. Re-enable HTTP control in Puma configuration
2. Reinstall old plugin
3. Update mackerel-agent configuration to use old plugin

## Support

For questions or issues:
- GitHub Issues: https://github.com/srockstyle/mackerel-plugin-puma-v2/issues
- Documentation: https://github.com/srockstyle/mackerel-plugin-puma-v2/blob/main/README.md