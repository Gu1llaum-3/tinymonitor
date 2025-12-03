# Metrics Overview

TinyMonitor comes with several built-in probes to monitor the health of your system.

Each metric can be configured independently in the `config.json` file. You can define:

*   **Enabled**: Whether the metric is active.
*   **Warning Threshold**: The value at which a WARNING alert is triggered.
*   **Critical Threshold**: The value at which a CRITICAL alert is triggered.
*   **Duration**: How long (in seconds) the threshold must be exceeded before alerting.

## Available Metrics

*   [🧠 CPU](cpu.md): Global CPU usage.
*   [💾 Memory](memory.md): Physical RAM usage.
*   [💿 Filesystem](filesystem.md): Disk space usage.
*   [⚙️ Load Average](load.md): System load (Unix only).
*   [📝 I/O](io.md): Disk I/O throughput.
*   [🔄 Reboot Required](reboot.md): Pending system reboots (Debian/Ubuntu).
