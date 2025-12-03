# Memory Metric

The Memory metric monitors the physical RAM usage percentage.

## How it works
It uses `psutil.virtual_memory().percent` to get the percentage of used physical memory. Swap memory is not included in this metric.

## Configuration

```json
"memory": {
    "enabled": true,
    "warning": 80,
    "critical": 95,
    "duration": 60
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
High memory usage is often normal for servers (caching). Ensure your thresholds reflect actual memory pressure rather than just cache usage.
