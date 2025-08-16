# Release Notes for v2.0.0

## 🎉 mackerel-plugin-puma-v2 Initial Release

This is the first official release of mackerel-plugin-puma-v2, a complete rewrite of the original mackerel-plugin-puma for modern Puma versions.

### ✨ Key Features

- **Full Puma 6.x Support**: Native support for all Puma 6.6.1 features including new metrics
- **Unix Socket Connection**: Secure and efficient monitoring via Unix domain sockets
- **Automatic Version Detection**: Automatically detects Puma version (4.x, 5.x, 6.x)
- **Extended Metrics**: Memory usage, GC statistics, and thread utilization metrics
- **Clean Architecture**: Modular design for maintainability and extensibility
- **Go 1.24**: Built with modern Go features for better performance

### 📊 Available Metrics

#### Core Metrics
- Worker process metrics (workers, booted_workers, old_workers)
- Thread metrics (running, pool_capacity, max_threads)
- Request backlog
- Phase information

#### Puma 6.x Specific
- `requests_count`: Total requests processed
- `uptime`: Server uptime
- `thread_utilization`: Thread utilization percentage

#### Extended Metrics (with `-extended` flag)
- Memory usage (alloc, sys, heap_inuse)
- Ruby GC statistics (count, minor/major GC, heap slots)
- Go runtime metrics (goroutines, GC count)

### 🔧 Installation

```bash
# Using mkr
sudo mkr plugin install srockstyle/mackerel-plugin-puma-v2

# From source
go install github.com/srockstyle/mackerel-plugin-puma-v2@latest
```

### 📝 Configuration

Basic setup in `/etc/mackerel-agent/mackerel-agent.conf`:

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2"
```

### 📚 Documentation

- [Quick Start Guide (日本語)](docs/quick-start.md)
- [Full Integration Guide (日本語)](docs/mackerel-integration-guide.md)
- [Migration from v1](MIGRATION.md)

### 🙏 Acknowledgments

- Original mackerel-plugin-puma by @rmanzoku
- Mackerel team for the plugin framework
- Puma team for the excellent web server

### 📄 License

MIT License - see [LICENSE](LICENSE) file for details

---

**Full Changelog**: https://github.com/srockstyle/mackerel-plugin-puma-v2/commits/v2.0.0