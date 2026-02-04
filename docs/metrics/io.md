# I/O Monitoring

TinyMonitor can monitor disk I/O throughput (Read and Write speeds).

## Configuration

To enable I/O monitoring, add the `io` section to your configuration file:

```toml
[io]
enabled = true
warning = "10MB"
critical = "50MB"
duration = 120
```

You can also define a maximum speed (`max_speed`) and use percentages:

```toml
[io]
enabled = true
max_speed = "100MB"
warning = "70%"
critical = "90%"
duration = 120
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |
| `warning` | `string/int` | - | Threshold for warning alert. Can be bytes (integer), string with unit (e.g. "10MB", "500KB"), or percentage if `max_speed` is defined. |
| `critical` | `string/int` | - | Threshold for critical alert. Same format as warning. |
| `max_speed` | `string/int` | - | Optional. The maximum disk speed used for percentage calculations. |
| `duration` | `int` | `120` | Time in seconds the value must be above threshold before alerting (2 minutes, avoids temporary spikes). |

## Behavior

The monitor calculates the I/O speed by comparing the disk counters between two checks. It reports the Read and Write speeds in a human-readable format (e.g., `R: 1.2MB/s W: 500KB/s`).

The alert is triggered based on the **total** throughput (Read + Write).
