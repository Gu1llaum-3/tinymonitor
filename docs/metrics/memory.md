# Memory Metric

The Memory metric monitors the physical RAM usage percentage.

## How it works

It uses the [gopsutil](https://github.com/shirou/gopsutil) library to get the percentage of used physical memory. Swap memory is not included in this metric.

## Configuration

```toml
[memory]
enabled = true
warning = 80
critical = 95
duration = 120
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | `70` | Percentage threshold for WARNING alert. |
| `critical` | `float` | `90` | Percentage threshold for CRITICAL alert. |
| `duration` | `int` | `120` | Time in seconds the value must be above threshold before alerting. |

### Recommendations

High memory usage is often normal for servers (caching). The default `duration` of **2 minutes (120 seconds)** helps avoid false positives from temporary memory spikes. Ensure your thresholds reflect actual memory pressure rather than just cache usage.
