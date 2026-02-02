# tinymonitor service

Manage the TinyMonitor systemd service (Linux only).

## Usage

```bash
tinymonitor service <command> [flags]
```

## Subcommands

### install

Install and start the systemd service.

```bash
sudo tinymonitor service install [flags]
```

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file to use |

If no config is specified, uses `/etc/tinymonitor/config.toml`.

**What it does:**

1. Creates `/etc/tinymonitor/` directory (if needed)
2. Copies the config file (if specified)
3. Creates the systemd service file
4. Enables the service to start on boot
5. Starts the service

### uninstall

Stop and remove the systemd service.

```bash
sudo tinymonitor service uninstall
```

**What it does:**

1. Stops the running service
2. Disables the service
3. Removes the service file

Configuration files are **preserved**.

### status

Show current service status.

```bash
tinymonitor service status
```

**Output includes:**

- Binary location and status
- Configuration file status
- Service status (running/stopped/not installed)

## Examples

```bash
# Install with default config location
sudo tinymonitor service install

# Install with custom config
sudo tinymonitor service install -c /opt/monitoring/config.toml

# Check status
tinymonitor service status

# Remove service (keeps config)
sudo tinymonitor service uninstall
```

## Requirements

- **Linux only**: This command uses systemd
- **Root privileges**: Required for install/uninstall

## Service File Location

The service file is created at:

```
/etc/systemd/system/tinymonitor.service
```

## See Also

- [Running as a systemd Service](../guides/systemd.md) - Manual setup and advanced options
- [Running on macOS](../guides/launchd.md) - macOS launchd setup
