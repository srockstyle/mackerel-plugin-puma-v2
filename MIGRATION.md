# Migration Guide: mackerel-plugin-puma to mackerel-plugin-puma-v2

## Overview

This document explains the migration from the original mackerel-plugin-puma to mackerel-plugin-puma-v2.

## Why v2?

The original mackerel-plugin-puma hasn't been updated since 2017 and has the following issues:

- No support for Puma 6.x features
- Limited error handling
- Outdated dependencies
- No active maintenance

mackerel-plugin-puma-v2 is a complete rewrite that addresses these issues with:

- Full Puma 6.x compatibility
- Modern Go architecture
- Comprehensive error handling
- Active maintenance

## Repository Setup

To use mackerel-plugin-puma-v2, you should:

1. **Fork or Clone this repository**
2. **Rename the directory** (optional but recommended):
   ```bash
   mv mackerel-plugin-puma mackerel-plugin-puma-v2
   cd mackerel-plugin-puma-v2
   ```

3. **Update the remote** (if forked):
   ```bash
   git remote set-url origin https://github.com/srockstyle/mackerel-plugin-puma-v2.git
   ```

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
The basic options remain the same, but v2 adds new options:

#### New Options in v2
- `-timeout duration` - Request timeout (default 10s)
- `-retry int` - Number of retry attempts (default 3)
- `-scheme string` - URL scheme support (http/https)

### Configuration Example

#### Old Configuration
```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma -token=secret"
```

#### New Configuration
```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -token=secret -timeout=30s -retry=5"
```

## New Metrics in v2

v2 adds several new metrics for Puma 6.x:

### Request Metrics
- `puma.requests_count` - Total requests processed
- `puma.uptime` - Server uptime in seconds

### Thread Utilization
- `puma.worker_thread_utilization.worker{N}.utilization` - Thread utilization percentage per worker

### Future Metrics (Planned)
- Memory usage per worker
- Request queue time
- Error rates

## Breaking Changes

### Metric Name Changes
None - v2 maintains backward compatibility with existing metric names.

### Removed Features
- GC stats collection (`-with-gc`) may not work with Puma 6.x as the endpoint might be removed

## Migration Steps

1. **Test in Development**
   ```bash
   # Test with your current Puma setup
   mackerel-plugin-puma-v2 -host=localhost -port=9293 -token=your-token
   ```

2. **Compare Metrics**
   - Run both plugins in parallel temporarily
   - Verify metrics are collected correctly

3. **Update Configuration**
   - Replace the old plugin command with v2
   - Add new options as needed

4. **Remove Old Plugin**
   ```bash
   # After successful migration
   rm /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma
   ```

## Rollback Plan

If you need to rollback:

1. Keep the old binary until migration is verified
2. Simply change the configuration back to use the old binary
3. Report any issues to: https://github.com/srockstyle/mackerel-plugin-puma-v2/issues

## Support

For questions or issues with v2:
- GitHub Issues: https://github.com/srockstyle/mackerel-plugin-puma-v2/issues
- Documentation: https://github.com/srockstyle/mackerel-plugin-puma-v2/blob/main/README.md