# Installation Guide

TinyMonitor can be installed as a standalone binary (recommended for servers) or via Python pip.

## 1. Install the Binary

=== "Linux (AMD64)"
    ```bash
    # Download
    wget -O /usr/local/bin/tinymonitor https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor-linux-amd64
    
    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

=== "Linux (ARM64 / RPi)"
    ```bash
    # Download
    wget -O /usr/local/bin/tinymonitor https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor-linux-arm64
    
    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

=== "macOS (Apple Silicon)"
    ```bash
    # Download
    curl -L -o /usr/local/bin/tinymonitor https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor-darwin-arm64
    
    # Make executable
    chmod +x /usr/local/bin/tinymonitor
    ```

## 2. Install via Python (Alternative)

If you cannot use the binary or prefer to run from source, you can install TinyMonitor using `pip`.

### Prerequisites
*   Python 3.8 or higher
*   `pip`

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/Gu1llaum-3/tinymonitor.git
    cd tinymonitor
    ```

2.  Install dependencies:
    ```bash
    pip install -r requirements.txt
    ```

3.  Run the application:
    ```bash
    python3 -m tinymonitor --config config.json
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
    # Config is automatically loaded from /etc/tinymonitor/config.json
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
            <string>--config</string>
            <string>/Users/YOUR_USER/.config/tinymonitor/config.json</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <true/>
    </dict>
    </plist>
    ```
