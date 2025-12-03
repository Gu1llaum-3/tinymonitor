# I/O Monitoring

TinyMonitor can monitor disk I/O throughput (Read and Write speeds).

## Configuration

To enable I/O monitoring, add the `io` section to your configuration file:

```json
"io": {
    "enabled": true,
    "warning": "10MB",
    "critical": "50MB",
    "duration": 0
}
```

You can also define a maximum speed (`max_speed`) and use percentages:

```json
"io": {
    "enabled": true,
    "max_speed": "100MB",
    "warning": "70%",
    "critical": "90%",
    "duration": 0
}
```

- **enabled**: Set to `true` to enable the metric.
- **warning**: Threshold for warning alert. Can be bytes (integer), string with unit (e.g. "10MB", "500KB"), or percentage (e.g. "70%") if `max_speed` is defined.
- **critical**: Threshold for critical alert. Same format as warning.
- **max_speed**: Optional. The maximum disk speed used for percentage calculations.
- **duration**: Number of consecutive checks the threshold must be exceeded before triggering an alert.

## Behavior

The monitor calculates the I/O speed by comparing the disk counters between two checks.
It reports the Read and Write speeds in a human-readable format (e.g., `R: 1.2MB/s W: 500KB/s`).

The alert is triggered based on the **total** throughput (Read + Write).
