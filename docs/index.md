![tinymonitor](assets/images/logo.png "Tinymonitor")
# Welcome to TinyMonitor

**TinyMonitor** is a lightweight system monitoring agent written in **Go**, designed for simplicity and ease of deployment.

It monitors your server's vital signs and alerts you immediately via your favorite channels when something goes wrong.

## Key Features

*   **Lightweight**: Single binary (~9MB), minimal CPU/RAM footprint.
*   **Zero Dependencies**: No runtime dependencies - just download and run.
*   **Multi-Platform**: Runs on Linux (AMD64/ARM64) and macOS (Intel/Silicon).
*   **Multi-Channel Alerts**: Ntfy, Gotify, Google Chat, SMTP, and Generic Webhooks.
*   **TOML Configuration**: Human-readable config format, easy to write and maintain.
*   **Self-Updating**: Built-in `update` command to stay current.
*   **Flexible Rules**: Route specific metrics to specific alert channels.

## Quick Start

Get up and running in seconds:

```bash
# Install
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash

# Configure and test
tinymonitor test-alert

# Run as a service
sudo tinymonitor service install
```

[Get Started](getting-started/index.md){ .md-button .md-button--primary }
[View All Commands](commands/index.md){ .md-button }
