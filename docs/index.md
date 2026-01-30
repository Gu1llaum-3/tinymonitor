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
*   **Config Validation**: Built-in `validate` command to check your configuration.
*   **Flexible Rules**: Route specific metrics to specific alert channels.
*   **Secure**: Runs as a non-privileged user (`nobody`).

## Quick Start

Get up and running in seconds:

```bash
# Download the latest binary
wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_x86_64.tar.gz
tar -xzf tinymonitor_Linux_x86_64.tar.gz

# Run with config
./tinymonitor_Linux_x86_64/tinymonitor -c config.toml
```

[Get Started with Installation](installation.md){ .md-button .md-button--primary }
