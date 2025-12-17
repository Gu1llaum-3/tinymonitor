![tinymonitor](assets/images/logo.png "Tinymonitor")
# Welcome to TinyMonitor

**TinyMonitor** is a lightweight, zero-dependency (binary) system monitoring agent designed for simplicity and ease of deployment.

It monitors your server's vital signs and alerts you immediately via your favorite channels when something goes wrong.

## ✨ Key Features

*   **🚀 Lightweight**: Consumes minimal CPU/RAM (~30MB).
*   **📦 Zero Dependency**: Available as a standalone binary (no Python required on target).
*   **🐧 Multi-Platform**: Runs on Linux (AMD64/ARM64) and macOS (Intel/Silicon).
*   **🔔 Multi-Channel Alerts**: Ntfy, Gotify, Google Chat, SMTP, and Generic Webhooks.
*   **🔌 Plugin System**: Modular architecture for easy extension.
*   **🛡️ Secure**: Runs as a non-privileged user (`nobody`).

## 🚀 Quick Start

Get up and running in seconds:

```bash
# Download the latest binary
wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor-linux-amd64
chmod +x tinymonitor-linux-amd64

# Run with default config
./tinymonitor-linux-amd64 --config config.json
```

[Get Started with Installation](installation.md){ .md-button .md-button--primary }
