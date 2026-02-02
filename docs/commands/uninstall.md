# tinymonitor uninstall

Completely remove TinyMonitor from the system.

## Usage

```bash
sudo tinymonitor uninstall [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--purge` | Also remove configuration files |
| `-y, --yes` | Skip confirmation prompt |

## What Gets Removed

### Without `--purge`

- Systemd service (if installed)
- Binary at `/usr/local/bin/tinymonitor`

### With `--purge`

- All of the above
- Configuration directory at `/etc/tinymonitor/`

## Examples

```bash
# Remove service and binary, keep configuration
sudo tinymonitor uninstall

# Remove everything including configuration
sudo tinymonitor uninstall --purge

# Skip confirmation prompt
sudo tinymonitor uninstall --yes

# Remove everything without confirmation
sudo tinymonitor uninstall --purge --yes
```

## Output

```
TinyMonitor Uninstallation
==========================

The following will be removed:
  - Service: /etc/systemd/system/tinymonitor.service
  - Binary:  /usr/local/bin/tinymonitor

The following will be preserved:
  - Config:  /etc/tinymonitor/

Are you sure you want to continue? [y/N]: y

Uninstalling TinyMonitor...

TinyMonitor has been uninstalled successfully!

Configuration files were preserved at /etc/tinymonitor/
Run with --purge to remove them as well.
```

## Requirements

- **Root privileges**: Required for uninstallation

## macOS Uninstallation

On macOS, use manual removal:

```bash
# Unload the service (if using launchd)
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist
rm ~/Library/LaunchAgents/com.tinymonitor.plist

# Remove the binary
sudo rm /usr/local/bin/tinymonitor

# Remove configuration (optional)
rm -rf ~/.config/tinymonitor
```

## See Also

- [Installation](../getting-started/installation.md)
