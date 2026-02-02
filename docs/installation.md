# Installation Guide

TinyMonitor is distributed as a standalone binary. No dependencies required.

## Quick Install (Linux / macOS)

The easiest way to install TinyMonitor:

```bash
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
```

This script will:

- Detect your system architecture (AMD64 or ARM64)
- Download the latest release from GitHub
- Install the binary to `/usr/local/bin`

### Environment Variables

You can customize the installation:

```bash
# Install a specific version
TINYMONITOR_VERSION=v1.0.0 curl -sSL ... | bash

# Install to a different directory
INSTALL_DIR=/opt/bin curl -sSL ... | bash
```

## Manual Installation

Download the latest release for your platform from the [Releases Page](https://github.com/Gu1llaum-3/tinymonitor/releases).

=== "Linux (AMD64)"
    ```bash
    wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_x86_64.tar.gz
    tar -xzf tinymonitor_Linux_x86_64.tar.gz
    sudo mv tinymonitor /usr/local/bin/
    ```

=== "Linux (ARM64 / RPi)"
    ```bash
    wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_arm64.tar.gz
    tar -xzf tinymonitor_Linux_arm64.tar.gz
    sudo mv tinymonitor /usr/local/bin/
    ```

=== "macOS (Apple Silicon)"
    ```bash
    curl -L -o tinymonitor.tar.gz https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Darwin_arm64.tar.gz
    tar -xzf tinymonitor.tar.gz
    sudo mv tinymonitor /usr/local/bin/
    ```

=== "macOS (Intel)"
    ```bash
    curl -L -o tinymonitor.tar.gz https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Darwin_x86_64.tar.gz
    tar -xzf tinymonitor.tar.gz
    sudo mv tinymonitor /usr/local/bin/
    ```

## Build from Source

If you prefer to build from source, you need Go 1.21 or higher.

```bash
git clone https://github.com/Gu1llaum-3/tinymonitor.git
cd tinymonitor
make build
```

## Running as a Service (Linux Systemd)

TinyMonitor includes built-in commands to manage the systemd service.

### Quick Setup

```bash
# Install service with default configuration
sudo tinymonitor service install

# Or use a custom configuration file
sudo tinymonitor service install -c /path/to/config.toml
```

This will:

1. Create `/etc/tinymonitor/` directory
2. Copy your configuration (or create a default one)
3. Create and enable the systemd service
4. Start monitoring

### Service Commands

```bash
# Check service status
tinymonitor service status

# Stop and remove the service (keeps configuration)
sudo tinymonitor service uninstall
```

### Standard Systemctl Commands

```bash
sudo systemctl status tinymonitor    # Check status
sudo systemctl restart tinymonitor   # Restart service
sudo systemctl stop tinymonitor      # Stop service
sudo journalctl -u tinymonitor -f    # View logs
```

## Running as a Service (macOS Launchd)

Create `~/Library/LaunchAgents/com.tinymonitor.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tinymonitor</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/tinymonitor</string>
        <string>-c</string>
        <string>/Users/YOUR_USER/.config/tinymonitor/config.toml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.tinymonitor.plist
```

## Uninstallation

### Linux

```bash
# Remove service and binary (keeps configuration)
sudo tinymonitor uninstall

# Remove everything including configuration
sudo tinymonitor uninstall --purge
```

### macOS

```bash
# Unload the service
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist
rm ~/Library/LaunchAgents/com.tinymonitor.plist

# Remove the binary
sudo rm /usr/local/bin/tinymonitor

# Remove configuration (optional)
rm -rf ~/.config/tinymonitor
```
