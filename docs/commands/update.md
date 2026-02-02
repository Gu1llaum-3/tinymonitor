# tinymonitor update

Check for updates and install the latest version from GitHub releases.

## Usage

```bash
tinymonitor update [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--check` | Check for updates without installing |
| `-y, --yes` | Update without confirmation |

## Examples

```bash
# Check if updates are available
tinymonitor update --check

# Interactive update (asks for confirmation)
tinymonitor update

# Automatic update (no confirmation)
tinymonitor update --yes
```

## Output

```
Checking for updates...

Current version: v1.2.0
Latest version:  v1.3.0

A new version is available!

Changelog: https://github.com/Gu1llaum-3/tinymonitor/releases/tag/v1.3.0

Do you want to update? [y/N]: y

Downloading tinymonitor v1.3.0...
Installing to /usr/local/bin/tinymonitor...

Update complete!

Note: TinyMonitor service is running.
Run 'sudo systemctl restart tinymonitor' to apply the update.
```

## Already Up to Date

```
Checking for updates...

Current version: v1.3.0
Latest version:  v1.3.0

Already up to date (v1.3.0)!
```

## Important Notes

- **Configuration preserved**: Your configuration files are never modified during updates
- **Service restart required**: If running as a service, restart it after updating
- **Permissions**: May require sudo to write to `/usr/local/bin`
- **Internet required**: Downloads from GitHub Releases

## After Updating

If TinyMonitor is running as a service:

```bash
# Linux (systemd)
sudo systemctl restart tinymonitor

# macOS (launchd)
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist
launchctl load ~/Library/LaunchAgents/com.tinymonitor.plist
```

## See Also

- [Installation](../getting-started/installation.md)
- [GitHub Releases](https://github.com/Gu1llaum-3/tinymonitor/releases)
