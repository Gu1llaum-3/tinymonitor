# CLI Reference

TinyMonitor provides a command-line interface for monitoring and management.

## Commands

| Command | Description |
|---------|-------------|
| [`tinymonitor`](run.md) | Start the monitoring agent |
| [`tinymonitor info`](info.md) | Display configuration summary |
| [`tinymonitor validate`](validate.md) | Validate configuration file |
| [`tinymonitor test-alert`](test-alert.md) | Send test notifications |
| [`tinymonitor update`](update.md) | Update to latest version |
| [`tinymonitor service`](service.md) | Manage systemd service |
| [`tinymonitor uninstall`](uninstall.md) | Remove TinyMonitor |
| [`tinymonitor completion`](completion.md) | Generate shell completions |
| [`tinymonitor version`](version.md) | Print version information |

## Global Flags

These flags are available for all commands:

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file |
| `-h, --help` | Help for any command |

## Getting Help

```bash
# General help
tinymonitor --help

# Help for a specific command
tinymonitor <command> --help
```
