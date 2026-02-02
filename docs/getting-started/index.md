# Quick Start

Get TinyMonitor running in under 5 minutes.

## 1. Install

```bash
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
```

## 2. Configure

Edit `/etc/tinymonitor/config.toml` (Linux) or `~/.config/tinymonitor/config.toml` (macOS):

```toml
refresh = 5
cooldown = 60

[cpu]
enabled = true
warning = 70
critical = 90

[memory]
enabled = true
warning = 80
critical = 95

[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/my-alerts"
```

## 3. Test

```bash
# Validate your configuration
tinymonitor validate

# Send a test alert
tinymonitor test-alert
```

## 4. Run

```bash
# Foreground (for testing)
tinymonitor

# As a service (Linux)
sudo tinymonitor service install
```

## Next Steps

- [Full Installation Guide](installation.md) - Manual installation, build from source
- [First Configuration](first-config.md) - Detailed configuration walkthrough
- [Configuration Reference](../configuration.md) - All available options
- [CLI Commands](../commands/index.md) - All available commands
