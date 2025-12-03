# Load Average Metric

The Load Average metric monitors the system load over 1, 5, and 15 minutes.

## How it works
It uses `os.getloadavg()` to retrieve the system load.
*   **Note**: This metric is only available on Unix-like systems (Linux, macOS). It is disabled on Windows.

## Configuration

```json
"load": {
    "enabled": true,
    "warning": 5.0,
    "critical": 10.0,
    "duration": 60
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `float` | `N/A` | Load threshold for WARNING alert. |
| `critical` | `float` | `N/A` | Load threshold for CRITICAL alert. |
| `duration` | `int` | `0` | Time in seconds the value must be above threshold before alerting. |

### Understanding Load
Load average represents the number of processes waiting for CPU time.
*   **Load < Number of Cores**: System is underutilized.
*   **Load > Number of Cores**: System is overutilized (processes are queuing).

Set your thresholds based on the number of CPU cores available on your server.
