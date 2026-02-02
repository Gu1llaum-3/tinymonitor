# First Configuration

This guide walks you through setting up your first TinyMonitor configuration.

## Configuration File Location

TinyMonitor searches for configuration files in this order:

1. **CLI flag**: `-c /path/to/config.toml`
2. **Current directory**: `./config.toml`
3. **User config**: `~/.config/tinymonitor/config.toml`
4. **System config**: `/etc/tinymonitor/config.toml`

The installation script creates a default configuration at:

- **Linux**: `/etc/tinymonitor/config.toml`
- **macOS**: `~/.config/tinymonitor/config.toml`

## Minimal Configuration

Here's a minimal configuration to get started:

```toml
# Check every 5 seconds
refresh = 5

# Wait 60 seconds between repeated alerts
cooldown = 60

[cpu]
enabled = true
warning = 70    # Alert at 70% usage
critical = 90   # Critical at 90% usage

[memory]
enabled = true
warning = 80
critical = 95

[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/my-unique-topic"
```

## Adding More Metrics

Enable additional metrics based on your needs:

```toml
[filesystem]
enabled = true
warning = 85
critical = 95

[load]
enabled = true
auto = true           # Calculates thresholds based on CPU count
warning_ratio = 0.7
critical_ratio = 0.9

[reboot]
enabled = true        # Debian/Ubuntu only
```

See [Metrics](../metrics/index.md) for detailed documentation on each metric.

## Adding Alert Providers

You can enable multiple alert providers simultaneously:

```toml
# Push notifications via ntfy.sh
[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/my-topic"

# Email alerts for critical issues only
[alerts.smtp]
enabled = true
host = "smtp.gmail.com"
port = 587
user = "your-email@gmail.com"
password = "your-app-password"
from_addr = "your-email@gmail.com"
to_addrs = ["admin@example.com"]
use_tls = true

  [alerts.smtp.rules]
  default = ["CRITICAL"]  # Only send emails for critical alerts
```

See [Alerts](../alerts/index.md) for all available providers.

## Testing Your Configuration

Before running TinyMonitor, validate and test your configuration:

```bash
# Validate syntax and required fields
tinymonitor validate -c /etc/tinymonitor/config.toml

# View parsed configuration summary
tinymonitor info -c /etc/tinymonitor/config.toml

# Send a test alert to all providers
tinymonitor test-alert -c /etc/tinymonitor/config.toml

# Test a specific provider
tinymonitor test-alert --provider ntfy
```

## Running TinyMonitor

Once configured, you can run TinyMonitor:

```bash
# Foreground (for testing)
tinymonitor -c /etc/tinymonitor/config.toml

# As a systemd service (Linux)
sudo tinymonitor service install

# Check service status
tinymonitor service status
```

## Next Steps

- [Configuration Reference](../configuration.md) - All available options
- [Running as a Service](../guides/systemd.md) - Systemd setup
- [Troubleshooting](../guides/troubleshooting.md) - Common issues
