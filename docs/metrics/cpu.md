# CPU Metric

The CPU metric monitors the global CPU usage percentage of the system.

## How it works

It uses the [gopsutil](https://github.com/shirou/gopsutil) library to measure the system-wide CPU utilization over a 1-second interval.

## Configuration

```toml
[cpu]
enabled = true
warning = 70
critical = 90
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

CPU usage can spike momentarily. The default `duration` of **2 minutes (120 seconds)** helps avoid false positives due to short bursts of activity (garbage collection, temporary load spikes).
