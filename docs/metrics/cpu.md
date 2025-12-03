# CPU Metric

The CPU metric monitors the global CPU usage percentage of the system.

## How it works
It uses `psutil.cpu_percent(interval=1)` to measure the system-wide CPU utilization over a 1-second interval.

## Configuration

```json
"cpu": {
    "enabled": true,
    "warning": 70,
    "critical": 90,
    "duration": 30
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
CPU usage can spike momentarily. It is recommended to set a `duration` of at least **30 seconds** to avoid false positives due to short bursts of activity.
