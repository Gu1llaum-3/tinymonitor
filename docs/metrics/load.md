# Load Average Metric

The Load Average metric monitors the system load over 1, 5, and 15 minutes.

## How it works

It uses the [gopsutil](https://github.com/shirou/gopsutil) library to retrieve the system load average.

*   **Note**: This metric is only available on Unix-like systems (Linux, macOS). It is disabled on Windows.

## Configuration

```toml
[load]
enabled = true
warning = 5.0
critical = 10.0
duration = 60
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | `CPU_COUNT * 0.7` | Load threshold for WARNING alert. |
| `critical` | `float` | `CPU_COUNT * 0.9` | Load threshold for CRITICAL alert. |
| `duration` | `int` | `60` | Time in seconds the value must be above threshold before alerting. |

### Understanding Load

Load average represents the number of processes waiting for CPU time.

*   **Load < Number of Cores**: System is underutilized.
*   **Load = Number of Cores**: System is fully utilized.
*   **Load > Number of Cores**: System is overutilized (processes are queuing).

By default, TinyMonitor sets the warning threshold to 70% of your CPU count and critical to 90%. For example, on a 4-core system:

*   Warning: 2.8
*   Critical: 3.6
