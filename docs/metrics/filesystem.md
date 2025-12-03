# Filesystem Metric

The Filesystem metric monitors disk space usage percentage for all mounted partitions.

## How it works
It iterates over all mounted partitions using `psutil.disk_partitions()` and checks usage with `psutil.disk_usage()`.

### Smart Filtering
To avoid noise, TinyMonitor automatically ignores the following filesystem types and mount points:
*   **Snap packages**: `/snap/*`, `squashfs`
*   **Docker**: `overlay`, `/var/lib/docker/*`
*   **Virtual**: `tmpfs`, `devtmpfs`, `proc`, `sysfs`

## Configuration

```json
"filesystem": {
    "enabled": true,
    "warning": 85,
    "critical": 95,
    "duration": 0
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `int` | `80` | Percentage threshold for WARNING alert. |
| `critical` | `int` | `95` | Percentage threshold for CRITICAL alert. |
| `duration` | `int` | `0` | Time in seconds the value must be above threshold before alerting. |

### Recommendations
Disk usage rarely fluctuates rapidly. A `duration` of `0` (immediate) is usually fine.
