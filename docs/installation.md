# Installation Guide

TinyMonitor is distributed as a standalone binary. No dependencies required.

## 1. Download the Binary

=== "Linux (AMD64)"
    ```bash
    # Download
    wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_x86_64.tar.gz
    tar -xzf tinymonitor_Linux_x86_64.tar.gz
    sudo mv tinymonitor_Linux_x86_64/tinymonitor /usr/local/bin/

    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

=== "Linux (ARM64 / RPi)"
    ```bash
    # Download
    wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Linux_arm64.tar.gz
    tar -xzf tinymonitor_Linux_arm64.tar.gz
    sudo mv tinymonitor_Linux_arm64/tinymonitor /usr/local/bin/

    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

=== "macOS (Apple Silicon)"
    ```bash
    # Download
    curl -L -o tinymonitor.tar.gz https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Darwin_arm64.tar.gz
    tar -xzf tinymonitor.tar.gz
    sudo mv tinymonitor_Darwin_arm64/tinymonitor /usr/local/bin/

    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

=== "macOS (Intel)"
    ```bash
    # Download
    curl -L -o tinymonitor.tar.gz https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor_Darwin_x86_64.tar.gz
    tar -xzf tinymonitor.tar.gz
    sudo mv tinymonitor_Darwin_x86_64/tinymonitor /usr/local/bin/

    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

## 2. Build from Source

If you prefer to build from source, you need Go 1.21 or higher.

```bash
# Clone the repository
git clone https://github.com/Gu1llaum-3/tinymonitor.git
cd tinymonitor

# Build
make build

# Or build manually
go build -o tinymonitor ./cmd/tinymonitor
```

## 3. Run as a Service

To ensure TinyMonitor runs in the background and starts on boot, configure it as a service.

=== "Linux (Systemd)"
    Create the file `/etc/systemd/system/tinymonitor.service`:

    ```ini
    [Unit]
    Description=TinyMonitor System Monitoring Service
    After=network.target

    [Service]
    Type=simple
    ExecStart=/usr/local/bin/tinymonitor
    Restart=on-failure

    # Security: Run as unprivileged user
    User=nobody
    Group=nogroup

    [Install]
    WantedBy=multi-user.target
    ```

    Enable and start the service:
    ```bash
    sudo mkdir -p /etc/tinymonitor
    sudo cp config.toml /etc/tinymonitor/config.toml
    sudo systemctl daemon-reload
    sudo systemctl enable --now tinymonitor
    sudo systemctl status tinymonitor
    ```

=== "macOS (Launchd)"
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
