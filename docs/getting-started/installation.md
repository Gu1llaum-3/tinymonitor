# Installation

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
- Copy a default configuration file

### Environment Variables

You can customize the installation:

```bash
# Install a specific version
TINYMONITOR_VERSION=v1.0.0 curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash

# Install to a different directory
INSTALL_DIR=/opt/bin curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
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

The binary will be created in the current directory.

## Verify Installation

```bash
tinymonitor version
```

## Next Steps

- [First Configuration](first-config.md) - Set up your configuration file
- [Running as a Service](../guides/systemd.md) - Run TinyMonitor in the background
