# Configuration

TinyMonitor uses TOML configuration files for clarity and ease of use.

## File Locations

TinyMonitor looks for the configuration file in the following order:

1.  **CLI Argument**: `tinymonitor -c /path/to/config.toml`
2.  **Current Directory**: `./config.toml`
3.  **User Config**: `~/.config/tinymonitor/config.toml`
4.  **System Config**: `/etc/tinymonitor/config.toml` (Recommended for Linux Services)

## Viewing Configuration Summary

To display a human-readable summary of your configuration:

```bash
tinymonitor info -c /path/to/config.toml
```

Example output:

```
Configuration: /etc/tinymonitor/config.toml

Global Settings
  Refresh:   5s
  Cooldown:  60s
  Log File:  (stdout)

Metrics
  [✓] CPU         warning: 70%    critical: 90%    duration: 30s
  [✓] Memory      warning: 80%    critical: 95%    duration: 60s
  [✓] Filesystem  warning: 85%    critical: 95%
  [✗] Load        (disabled)

Alert Providers
  [✓] Ntfy        https://ntfy.sh/my_topic
  [✗] Google Chat
  [✓] SMTP        smtp.gmail.com:587 → 2 recipient(s)
```

## Validating Configuration

Before deploying, you can validate your configuration file:

```bash
tinymonitor validate -c /path/to/config.toml
```

This checks for:

*   Valid TOML syntax
*   Required fields when providers are enabled
*   Threshold values (0-100 for percentages)
*   `warning < critical` for all metrics
*   Valid port numbers

## Example Configuration

```toml
# Global settings
refresh = 5          # How often (in seconds) to check metrics
cooldown = 60        # Minimum time between repeat notifications (-1 = alert once per incident)
log_file = ""        # Leave empty for stdout, or set a path like "/var/log/tinymonitor.log"

# CPU monitoring
[cpu]
enabled = true
warning = 70         # Percentage threshold for WARNING alert
critical = 90        # Percentage threshold for CRITICAL alert
duration = 30        # Condition must persist for X seconds before alerting

# Memory monitoring
[memory]
enabled = true
warning = 80
critical = 95
duration = 60

# Filesystem monitoring
[filesystem]
enabled = true
warning = 85
critical = 95
duration = 0         # Duration 0 = alert immediately

# Alert provider
[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/my_secret_topic"

  [alerts.ntfy.rules]
  default = ["WARNING", "CRITICAL"]
```

## Configuration Reference

### Global Settings

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `refresh` | `int` | `2` | How often (in seconds) to check metrics. |
| `cooldown` | `int` | `60` | Minimum time (in seconds) between repeat notifications. Set to `-1` to alert only once per incident. |
| `log_file` | `string` | `""` | Path to log file. Empty = stdout only. |

### Metric Settings

Each metric (`cpu`, `memory`, `filesystem`, `load`) supports:

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | varies | Threshold for WARNING alert (0-100 for percentages). |
| `critical` | `float` | varies | Threshold for CRITICAL alert (0-100 for percentages). |
| `duration` | `int` | `0` | Time in seconds the value must exceed threshold before alerting. |

### Alert Rules

Each alert provider supports filtering rules to control which alerts are sent:

```toml
[alerts.ntfy.rules]
default = ["WARNING", "CRITICAL"]  # Default levels for all metrics
cpu = ["CRITICAL"]                  # Only CRITICAL for CPU
filesystem = ["WARNING", "CRITICAL"]
```

## Complete Example

See [configs/config.example.toml](https://github.com/Gu1llaum-3/tinymonitor/blob/main/configs/config.example.toml) for a complete configuration example with all available options.
