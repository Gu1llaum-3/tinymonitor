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
duration = 120       # Condition must persist for 2 minutes before alerting (avoids short spikes)

# Memory monitoring
[memory]
enabled = true
warning = 80
critical = 95
duration = 120       # Condition must persist for 2 minutes before alerting (avoids short spikes)

# Filesystem monitoring
[filesystem]
enabled = true
warning = 85
critical = 95
duration = 300       # Condition must persist for 5 minutes (disk fills slowly)

# Load average monitoring (auto-adjusts to CPU count)
[load]
enabled = true
auto = true          # Calculate thresholds based on CPU count
warning_ratio = 0.7  # warning = CPU_COUNT × 0.7
critical_ratio = 0.9 # critical = CPU_COUNT × 0.9
duration = 180       # Condition must persist for 3 minutes (load has natural inertia)

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

Each metric (`cpu`, `memory`, `filesystem`) supports:

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | varies | Threshold for WARNING alert (0-100 for percentages). |
| `critical` | `float` | varies | Threshold for CRITICAL alert (0-100 for percentages). |
| `duration` | `int` | varies | Time in seconds the value must exceed threshold before alerting. Default: 120s (CPU, Memory), 300s (Filesystem). |

### Load Average Settings

The `load` metric has an auto mode that calculates thresholds based on CPU count:

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `auto` | `bool` | `true` | Calculate thresholds based on CPU count. |
| `warning_ratio` | `float` | `0.7` | Multiplier for warning (auto mode): `CPU_COUNT × ratio`. |
| `critical_ratio` | `float` | `0.9` | Multiplier for critical (auto mode): `CPU_COUNT × ratio`. |
| `warning` | `float` | - | Absolute threshold (manual mode, requires `auto = false`). |
| `critical` | `float` | - | Absolute threshold (manual mode, requires `auto = false`). |
| `duration` | `int` | `180` | Time in seconds before alerting (3 minutes, load has natural inertia). |

See [Load Average Metric](metrics/load.md) for more details.

### Alert Settings

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `send_recovery` | `bool` | `true` | Send notification when a metric returns to normal. |

### Recovery Notifications

When a metric returns to normal after triggering an alert, TinyMonitor can send a recovery notification:

```toml
[alerts]
send_recovery = true  # Enable recovery notifications (default: true)
```

Recovery notifications are sent to the same providers that received the original alert. They have distinct formatting:

*   **Ntfy**: Low priority, green checkmark emoji
*   **Google Chat**: Green color, checkmark icon
*   **SMTP**: Green header, "[RECOVERED]" in subject
*   **Webhook**: `"level": "RECOVERED"` with `"previous_level"` field
*   **Gotify**: Lower priority (3)

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
