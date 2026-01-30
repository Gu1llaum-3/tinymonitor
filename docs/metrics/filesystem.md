# Filesystem Metric

The Filesystem metric monitors disk space usage percentage for all mounted partitions.

## How it works

It uses the [gopsutil](https://github.com/shirou/gopsutil) library to iterate over all mounted partitions and check usage.

### Smart Filtering

To avoid noise, TinyMonitor automatically ignores the following filesystem types and mount points:

*   **Snap packages**: `/snap/*`, `squashfs`
*   **Docker**: `overlay`, `/var/lib/docker/*`
*   **Virtual**: `tmpfs`, `devtmpfs`, `proc`, `sysfs`

You can also exclude additional mount points in the configuration.

## Configuration

```toml
[filesystem]
enabled = true
warning = 85
critical = 95
duration = 0
exclude = ["/mnt/backup", "/media/usb"]
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | `80` | Percentage threshold for WARNING alert. |
| `critical` | `float` | `90` | Percentage threshold for CRITICAL alert. |
| `duration` | `int` | `0` | Time in seconds the value must be above threshold before alerting. |
| `exclude` | `list` | `[]` | List of mount points to exclude from monitoring. |

### Recommendations

Disk usage rarely fluctuates rapidly. A `duration` of `0` (immediate) is usually fine.
