# Load Average Metric

The Load Average metric monitors the system load over 1, 5, and 15 minutes.

## How it works

It uses the [gopsutil](https://github.com/shirou/gopsutil) library to retrieve the system load average.

*   **Note**: This metric is only available on Unix-like systems (Linux, macOS). It is disabled on Windows.

## Understanding Load Average

Load average represents the average number of processes waiting for CPU time:

*   **Load < Number of Cores**: System is underutilized
*   **Load = Number of Cores**: System is fully utilized
*   **Load > Number of Cores**: System is overloaded (processes are queuing)

For example, on a 4-core system:
- Load of 2.0 = 50% utilization
- Load of 4.0 = 100% utilization
- Load of 8.0 = 200% utilization (overloaded)

## Configuration

### Auto Mode (Recommended)

By default, TinyMonitor automatically calculates thresholds based on your CPU count:

```toml
[load]
enabled = true
auto = true           # Calculate thresholds based on CPU count
warning_ratio = 0.7   # warning = CPU_COUNT × 0.7
critical_ratio = 0.9  # critical = CPU_COUNT × 0.9
duration = 60
```

With auto mode, the same configuration works on any machine:

| CPUs | Warning (0.7×) | Critical (0.9×) |
|------|----------------|-----------------|
| 2    | 1.4            | 1.8             |
| 4    | 2.8            | 3.6             |
| 8    | 5.6            | 7.2             |
| 16   | 11.2           | 14.4            |

### Manual Mode

For specific thresholds, disable auto mode:

```toml
[load]
enabled = true
auto = false
warning = 5.0
critical = 10.0
duration = 60
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `auto` | `bool` | `true` | Calculate thresholds based on CPU count. |
| `warning_ratio` | `float` | `0.7` | Multiplier for warning threshold (auto mode). |
| `critical_ratio` | `float` | `0.9` | Multiplier for critical threshold (auto mode). |
| `warning` | `float` | - | Absolute warning threshold (manual mode). |
| `critical` | `float` | - | Absolute critical threshold (manual mode). |
| `duration` | `int` | `60` | Time in seconds the value must exceed threshold before alerting. |

## Viewing Effective Thresholds

Use the `info` command to see the calculated thresholds on your system:

```bash
tinymonitor info
```

Example output on an 8-core system:

```
Metrics
  [✓] Load        warning: 5.6     critical: 7.2    duration: 60s    (auto: 8 CPUs)
```

## Tips

*   **Use auto mode** for portable configurations that work across different machines
*   **Set duration to 60s** or more to avoid false positives from temporary load spikes
*   **Adjust ratios** if you want more/less headroom:
    - Conservative: `warning_ratio = 0.5`, `critical_ratio = 0.7`
    - Aggressive: `warning_ratio = 0.8`, `critical_ratio = 1.0`
