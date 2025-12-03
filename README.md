# TinyMonitor

<p align="center">
  <img src="docs/assets/images/logo.png" alt="TinyMonitor Logo" width="200"/>
</p>

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Python](https://img.shields.io/badge/python-3.8%2B-blue)
![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos-lightgrey)

**TinyMonitor** is a lightweight, zero-dependency (runtime) system monitoring agent designed for simplicity and performance. It runs silently in the background, watching your system resources, and alerts you immediately when something goes wrong.

Think of it as a minimalist `glances action` that never sleeps and knows how to send notifications.

## ✨ Features

*   **🚀 Lightweight**: Minimal footprint, written in pure Python.
*   **🔌 Plugin Architecture**: Easily extensible metrics and alert providers.
*   **🔔 Multi-Channel Alerts**: Native support for Google Chat, ntfy.sh, and SMTP (Email).
*   **📦 Standalone Binary**: Available as a single executable file (no Python installation required).
*   **🐧 Linux & macOS**: Fully compatible with major Unix-like systems (AMD64 & ARM64).

## 📥 Installation

### Option 1: Standalone Binary (Recommended)
Download the latest release for your OS from the [Releases Page](https://github.com/Gu1llaum-3/tinymonitor/releases).

```bash
# Example for Linux AMD64
wget https://github.com/Gu1llaum-3/tinymonitor/releases/latest/download/tinymonitor-linux-amd64
chmod +x tinymonitor-linux-amd64
sudo mv tinymonitor-linux-amd64 /usr/local/bin/tinymonitor
```

### Option 2: Via Pip (Python 3.8+)

```bash
pip install .
# OR with pipx (recommended for isolation)
pipx install .
```

## ⚙️ Configuration

TinyMonitor automatically searches for a configuration file in the following order:

1.  **Command Line Flag**: `--config /path/to/config.json`
2.  **Current Directory**: `./config.json`
3.  **User Config**: `~/.config/tinymonitor/config.json`
4.  **System Config**: `/etc/tinymonitor/config.json`

**Minimal `config.json` example:**

```json
{
    "check_interval": 5,
    "metrics": {
        "cpu": { "enabled": true, "threshold": 90 },
        "memory": { "enabled": true, "threshold": 85 },
        "disk": { "enabled": true, "threshold": 90, "path": "/" }
    },
    "alerts": {
        "console": { "enabled": true },
        "google_chat": {
            "enabled": false,
            "webhook_url": "YOUR_WEBHOOK_URL"
        }
    }
}
```

## 🚀 Usage

Run it directly in your terminal:

```bash
tinymonitor --config /path/to/config.json
```

Check the version:

```bash
tinymonitor --version
```

## 🤖 Running as a Service (Linux Systemd)

To keep TinyMonitor running in the background and starting automatically on boot, create a Systemd service.

1.  **Create the configuration directory:**
    ```bash
    sudo mkdir -p /etc/tinymonitor
    sudo cp config.json /etc/tinymonitor/config.json
    ```

2.  **Configure Logging (Important):**
    TinyMonitor runs as the `nobody` user for security.
    
    *   **Option A: Use Systemd Journal (Recommended)**
        Set `"log_file": ""` (empty string) in your `config.json`. Logs will be handled by `journalctl`.
        
    *   **Option B: Use a Log File**
        If you want a specific file (e.g., `/var/log/tinymonitor.log`), you must create it and give permissions to `nobody`:
        ```bash
        sudo touch /var/log/tinymonitor.log
        sudo chown nobody:nogroup /var/log/tinymonitor.log
        ```
        Then set `"log_file": "/var/log/tinymonitor.log"` in `config.json`.

3.  **Create the service file:**
    Create a file at `/etc/systemd/system/tinymonitor.service`:

    ```ini
    [Unit]
    Description=TinyMonitor System Monitoring Service
    After=network.target

    [Service]
    Type=simple
    # Adjust path if you installed via pip/pipx
    # Config is automatically loaded from /etc/tinymonitor/config.json
    ExecStart=/usr/local/bin/tinymonitor
    Restart=on-failure
    User=nobody
    Group=nogroup

    [Install]
    WantedBy=multi-user.target
    ```

4.  **Enable and Start:**
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable --now tinymonitor
    ```

5.  **Check Status & Logs:**
    ```bash
    # Check service status
    systemctl status tinymonitor
    
    # View logs (if using Option A)
    journalctl -u tinymonitor -f
    ```

## 🛠️ Development

Clone the repository and install dependencies:

```bash
git clone https://github.com/Gu1llaum-3/tinymonitor.git
cd tinymonitor
python -m venv .venv
source .venv/bin/activate
pip install -e .
```

Run tests:

```bash
python -m unittest discover tests
```

## 📄 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
