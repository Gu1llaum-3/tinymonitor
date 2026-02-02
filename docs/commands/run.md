# tinymonitor

Start the monitoring agent.

## Usage

```bash
tinymonitor [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file |
| `-h, --help` | Help for tinymonitor |

## Configuration File Search Order

If no config file is specified, TinyMonitor searches in this order:

1. `./config.toml` (current directory)
2. `~/.config/tinymonitor/config.toml` (user config)
3. `/etc/tinymonitor/config.toml` (system config)

## Examples

```bash
# Run with auto-detected config
tinymonitor

# Run with specific config file
tinymonitor -c /path/to/config.toml

# Run with config in current directory
tinymonitor -c ./my-config.toml
```

## Behavior

When started, TinyMonitor will:

1. Load and validate the configuration
2. Initialize all enabled metrics collectors
3. Initialize all enabled alert providers
4. Start the monitoring loop at the configured refresh interval
5. Send alerts when thresholds are exceeded

## Graceful Shutdown

TinyMonitor handles these signals for graceful shutdown:

- `SIGINT` (Ctrl+C)
- `SIGTERM`

When a shutdown signal is received, TinyMonitor will:

1. Stop the monitoring loop
2. Wait for pending alerts to be sent
3. Exit cleanly

## Logging

By default, logs are written to stdout. You can configure a log file in `config.toml`:

```toml
log_file = "/var/log/tinymonitor.log"
```

## See Also

- [Configuration Reference](../configuration.md)
- [Running as a Service](../guides/systemd.md)
