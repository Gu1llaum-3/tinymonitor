# tinymonitor info

Display a summary of your configuration including enabled metrics, thresholds, and alert providers.

## Usage

```bash
tinymonitor info [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file |

## Output

The command displays:

- **Global settings**: refresh interval, cooldown period
- **Enabled metrics**: with their warning and critical thresholds
- **Load average**: shows effective thresholds when using auto mode
- **Alert providers**: enabled providers and their rules

## Example

```bash
$ tinymonitor info

Configuration Summary
=====================
Refresh: 5s | Cooldown: 60s

Metrics:
  CPU:        warning=70%, critical=90%, duration=30s
  Memory:     warning=80%, critical=95%, duration=0s
  Load:       warning=5.6, critical=7.2 (auto: 8 CPUs)
  Filesystem: warning=85%, critical=95%, duration=0s

Alerts:
  ntfy:       enabled (default: WARNING, CRITICAL)
  smtp:       enabled (default: CRITICAL)
  gotify:     disabled
```

## Use Cases

- **Verify configuration** before deploying
- **Check effective thresholds** for load average auto mode
- **Review alert routing** rules

## See Also

- [Configuration Reference](../configuration.md)
- [tinymonitor validate](validate.md)
