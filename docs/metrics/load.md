# Load Average Metric

The Load Average metric alerts on the system load average. Two windows are
supported: the **5-minute** average (monitored by default) and the **15-minute**
average (opt-in). The 1-minute average is intentionally not exposed for alerting
as it is too noisy and prone to false positives.

> This mirrors how tools like Prometheus/Nagios alert on load: a longer, smoother
> window rather than the spiky 1-minute value.

## Migrating from a single `[load]` section

Earlier versions exposed a single flat `[load]` section keyed on the 1-minute
average. If you are upgrading, two **breaking changes** apply:

- **`[load] duration` moved to `[load.window5] duration`.** A top-level
  `duration` is now silently ignored, so the 5-minute window falls back to its
  default of `300`. Move your value into `[load.window5]` (and/or `[load.window15]`).
- **Alert routing key `load` is gone.** Rules must now key on `load5` and/or
  `load15`. A rule still keyed on `load` no longer matches and its alerts fall
  through to `default` (or are dropped if no `default` rule exists).

The shared threshold keys (`auto`, `warning_ratio`, `critical_ratio`, `warning`,
`critical`) are unchanged.

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

Thresholds are defined once at the `[load]` level (shared by both windows), while
each window controls whether it is monitored and after how long it alerts.

### Auto Mode (Recommended)

By default, TinyMonitor automatically calculates thresholds based on your CPU count:

```toml
[load]
enabled = true        # Master switch for load monitoring
auto = true           # Calculate thresholds based on CPU count
warning_ratio = 0.7   # warning = CPU_COUNT × 0.7
critical_ratio = 0.9  # critical = CPU_COUNT × 0.9

# 5-minute average (Prometheus-style: alert when sustained above threshold)
[load.window5]
enabled = true
duration = 300        # "for: 5m" — confirmation window before alerting

# 15-minute average (opt-in)
[load.window15]
enabled = false
duration = 0          # 0 = alert immediately (the 15m window is its own smoothing)
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

[load.window5]
enabled = true
duration = 300
```

### Per-window threshold overrides

Each window may override the shared thresholds — useful to give the 15-minute
window a lower bar (sustained load matters more than a short burst). An unset
(zero) override inherits the `[load]` default.

```toml
[load.window15]
enabled = true
duration = 0
warning_ratio = 0.6   # overrides [load].warning_ratio for the 15m window only
critical_ratio = 0.8
```

### Parameters

#### `[load]` (shared)

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Master switch for load monitoring. |
| `auto` | `bool` | `true` | Calculate thresholds based on CPU count. |
| `warning_ratio` | `float` | `0.7` | Multiplier for warning threshold (auto mode). |
| `critical_ratio` | `float` | `0.9` | Multiplier for critical threshold (auto mode). |
| `warning` | `float` | - | Absolute warning threshold (manual mode). |
| `critical` | `float` | - | Absolute critical threshold (manual mode). |

#### `[load.window5]` / `[load.window15]` (per window)

| Parameter | Type | Default (5m / 15m) | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` / `false` | Monitor this window. |
| `duration` | `int` | `300` / `0` | Seconds the value must stay above threshold before alerting (`0` = immediate). |
| `warning` / `critical` | `float` | inherit | Optional per-window threshold override (manual mode). |
| `warning_ratio` / `critical_ratio` | `float` | inherit | Optional per-window ratio override (auto mode). |

## Alert routing

Each window is a distinct component for alert routing rules: use the keys
`load5` and `load15` (the 1-minute key `load` no longer exists).

```toml
[alerts.ntfy.rules]
load5 = ["WARNING", "CRITICAL"]
load15 = ["CRITICAL"]
```

## Viewing Effective Thresholds

Use the `info` command to see the calculated thresholds on your system:

```bash
tinymonitor info
```

Example output on an 8-core system (15m window disabled):

```
Metrics
  [✓] Load 5m  warning: 5.6     critical: 7.2    duration: 300s   (auto: 8 CPUs)
  [✗] Load 15m (disabled)
```

## Tips

*   **Use auto mode** for portable configurations that work across different machines
*   **The 5-minute window + `duration = 300`** ("for: 5m") avoids false positives from temporary load spikes — the same conservative default used by common Prometheus rules
*   **Enable the 15-minute window** when you care about sustained, long-running load rather than short bursts; `duration = 0` is fine there since the window itself smooths
*   **Adjust ratios** if you want more/less headroom:
    - Conservative: `warning_ratio = 0.5`, `critical_ratio = 0.7`
    - Aggressive: `warning_ratio = 0.8`, `critical_ratio = 1.0`
