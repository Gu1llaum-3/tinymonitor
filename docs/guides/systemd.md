# Running as a systemd Service (Linux)

This guide explains how to run TinyMonitor as a background service on Linux using systemd.

## Quick Setup

The easiest way to install TinyMonitor as a service:

```bash
sudo tinymonitor service install
```

This will:

1. Create `/etc/tinymonitor/` directory (if needed)
2. Create the systemd service file
3. Enable the service to start on boot
4. Start the service

### With Custom Configuration

```bash
sudo tinymonitor service install -c /path/to/your/config.toml
```

## Manual Setup

If you prefer manual control over the service configuration:

### 1. Create the service file

```bash
sudo nano /etc/systemd/system/tinymonitor.service
```

```ini
[Unit]
Description=TinyMonitor - Lightweight System Monitoring
Documentation=https://github.com/Gu1llaum-3/tinymonitor
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
ExecStart=/usr/local/bin/tinymonitor -c /etc/tinymonitor/config.toml
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=read-only
PrivateTmp=true
ReadOnlyPaths=/

[Install]
WantedBy=multi-user.target
```

### 2. Reload systemd

```bash
sudo systemctl daemon-reload
```

### 3. Enable and start the service

```bash
sudo systemctl enable tinymonitor
sudo systemctl start tinymonitor
```

## Management Commands

### Check status

```bash
sudo systemctl status tinymonitor
```

Or use the built-in command:

```bash
tinymonitor service status
```

### View logs

```bash
# Live logs
sudo journalctl -u tinymonitor -f

# Last 100 lines
sudo journalctl -u tinymonitor -n 100

# Logs since last boot
sudo journalctl -u tinymonitor -b

# Logs from specific time
sudo journalctl -u tinymonitor --since "1 hour ago"
```

### Restart service

After changing the configuration:

```bash
sudo systemctl restart tinymonitor
```

### Stop service

```bash
sudo systemctl stop tinymonitor
```

### Disable service

Prevent from starting on boot:

```bash
sudo systemctl disable tinymonitor
```

## Service User

By default, the service runs as the `nobody` user for security. If you need to monitor paths that require specific permissions, you can change the user in the service file:

```ini
[Service]
User=root
Group=root
```

Then reload:

```bash
sudo systemctl daemon-reload
sudo systemctl restart tinymonitor
```

## Uninstall Service

```bash
sudo tinymonitor service uninstall
```

This removes the service but preserves your configuration files.

## Troubleshooting

### Service fails to start

Check the logs:

```bash
sudo journalctl -u tinymonitor -n 50 --no-pager
```

Common issues:

- **Configuration error**: Run `tinymonitor validate -c /etc/tinymonitor/config.toml`
- **Permission denied**: Check file permissions on config file
- **Binary not found**: Verify `which tinymonitor` returns `/usr/local/bin/tinymonitor`

### Service starts but stops immediately

The configuration might be invalid. Test manually:

```bash
tinymonitor -c /etc/tinymonitor/config.toml
```

## See Also

- [tinymonitor service](../commands/service.md) - Service command reference
- [Running on macOS](launchd.md) - macOS launchd setup
- [Troubleshooting](troubleshooting.md) - Common issues
