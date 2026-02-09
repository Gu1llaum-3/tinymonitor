# TinyMonitor

<p align="center">
  <img src="docs/assets/images/logo.png" alt="TinyMonitor Logo" width="200"/>
</p>

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Gu1llaum-3/tinymonitor?style=for-the-badge)](https://github.com/Gu1llaum-3/tinymonitor/releases)
[![License](https://img.shields.io/github/license/Gu1llaum-3/tinymonitor?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey?style=for-the-badge&logo=linux)](https://github.com/Gu1llaum-3/tinymonitor/releases)

**TinyMonitor** is a lightweight system monitoring agent written in Go, designed for simplicity and performance. It runs silently in the background, watching your system resources, and alerts you immediately when something goes wrong.

## Features

*   **Lightweight**: Single binary, minimal footprint (~9MB), low CPU/RAM usage.
*   **Zero Dependencies**: No runtime dependencies, just download and run.
*   **Multi-Channel Alerts**: Google Chat, Ntfy, Gotify, SMTP (Email), and Generic Webhooks.
*   **TOML Configuration**: Human-readable config with per-metric thresholds, durations, and alert routing rules.
*   **Config Validation**: Built-in `validate` command to check your configuration before deployment.
*   **Cross-Platform**: Linux and macOS (AMD64 & ARM64).

## Installation

### Quick Install (Linux / macOS)

```bash
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
```

### Manual Installation

Download the latest release for your platform from the [Releases Page](https://github.com/Gu1llaum-3/tinymonitor/releases).

```bash
# Linux AMD64
wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_x86_64.tar.gz
tar -xzf tinymonitor_Linux_x86_64.tar.gz
sudo mv tinymonitor /usr/local/bin/

# Linux ARM64
wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_arm64.tar.gz
tar -xzf tinymonitor_Linux_arm64.tar.gz
sudo mv tinymonitor /usr/local/bin/

# macOS (Apple Silicon)
curl -L -o tinymonitor.tar.gz https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Darwin_arm64.tar.gz
tar -xzf tinymonitor.tar.gz
sudo mv tinymonitor /usr/local/bin/
```

## Configuration

TinyMonitor uses TOML configuration and searches for files in this order:

1.  **CLI Flag**: `-c /path/to/config.toml`
2.  **Current Directory**: `./config.toml`
3.  **User Config**: `~/.config/tinymonitor/config.toml`
4.  **System Config**: `/etc/tinymonitor/config.toml`

**Example `config.toml`:**

```toml
refresh = 5
cooldown = 60

[cpu]
enabled = true
warning = 70
critical = 90
duration = 30

[memory]
enabled = true
warning = 80
critical = 95
duration = 60

[filesystem]
enabled = true
warning = 85
critical = 95

[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/your_topic"
```

See [configs/config.example.toml](configs/config.example.toml) for a complete example.

## Usage

```bash
# Run with auto-detected config
tinymonitor

# Run with specific config
tinymonitor -c /path/to/config.toml

# Display configuration summary
tinymonitor info -c /path/to/config.toml

# Validate configuration before deployment
tinymonitor validate -c /path/to/config.toml

# Test alert notifications
tinymonitor test-alert                    # Test all providers
tinymonitor test-alert --provider ntfy    # Test specific provider

# Check for updates
tinymonitor update --check

# Update to latest version
tinymonitor update

# Show version
tinymonitor version
```

## Running as a Service (Linux Systemd)

TinyMonitor includes built-in commands to manage the systemd service.

### Quick Setup

```bash
# Install as a service with default configuration
sudo tinymonitor service install

# Or use a custom configuration file
sudo tinymonitor service install -c /path/to/config.toml
```

### Service Management

```bash
# Check service status
tinymonitor service status

# Stop and remove the service (keeps config)
sudo tinymonitor service uninstall

# Standard systemctl commands work too
sudo systemctl status tinymonitor
sudo systemctl restart tinymonitor
sudo journalctl -u tinymonitor -f
```

### Complete Uninstallation

```bash
# Remove service and binary (keeps config)
sudo tinymonitor uninstall

# Remove everything including configuration
sudo tinymonitor uninstall --purge
```

## Updating

TinyMonitor can update itself to the latest version:

```bash
# Check for updates
tinymonitor update --check

# Update to latest version
tinymonitor update

# Update without confirmation
tinymonitor update --yes
```

Your configuration files are never modified during updates. If the service is running, you'll be reminded to restart it after the update.

## Building from Source

Requires Go 1.25+.

```bash
git clone https://github.com/Gu1llaum-3/tinymonitor.git
cd tinymonitor

# Build
make build

# Run tests
make test

# Build for all platforms
make release
```

## License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
