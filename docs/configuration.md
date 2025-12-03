# Configuration

TinyMonitor uses a simple JSON configuration file.

## File Locations

TinyMonitor looks for the configuration file in the following order:

1.  **CLI Argument**: `tinymonitor -c /path/to/config.json`
2.  **Current Directory**: `./config.json`
3.  **User Config**: `~/.config/tinymonitor/config.json`
4.  **System Config**: `/etc/tinymonitor/config.json` (Recommended for Linux Services)

## Example Configuration

```json
{
    "refresh": 5,          // (1)
    "cooldown": 60,        // (2)
    "log_file": "",        // (3)
    
    "cpu": {
        "enabled": true,
        "warning": 70,
        "critical": 90,
        "duration": 30     // (4)
    },
    "memory": {
        "enabled": true,
        "warning": 80,
        "critical": 95,
        "duration": 60
    },
    "filesystem": {
        "enabled": true,
        "warning": 85,
        "critical": 95,
        "duration": 0      // (5)
    },
    "alerts": {
        "ntfy": {
            "enabled": true,
            "topic_url": "https://ntfy.sh/my_secret_topic",
            "rules": {
                "default": ["WARNING", "CRITICAL"]
            }
        }
    }
}
```

1.  **Refresh Rate**: How often (in seconds) to check metrics.
2.  **Global Cooldown**: Minimum time (in seconds) between repeat notifications for the same issue. Set to `-1` to alert only once per incident.
3.  **Logging**: Leave empty `""` to log to stdout (best for Systemd/Docker). Set a path like `/var/log/tinymonitor.log` for a file.
4.  **Duration**: The condition must persist for X seconds before alerting. Prevents false positives on CPU spikes.
5.  **Immediate Alert**: Duration 0 means alert immediately.
