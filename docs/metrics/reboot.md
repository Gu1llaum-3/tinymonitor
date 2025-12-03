# Reboot Required Metric

The Reboot Required metric checks if the system needs a restart following system updates (kernel, libc, etc.).

## How it works
It checks for the existence of the file `/var/run/reboot-required`.
*   **Supported OS**: Debian, Ubuntu, and their derivatives.
*   **Behavior**: If the file exists, it triggers a **WARNING** alert.

## Configuration

```json
"reboot": {
    "enabled": true
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `true` | Enable or disable this metric. |

### Note
This metric does not have `warning` or `critical` thresholds as it is a binary state (Reboot needed: Yes/No). It always alerts at the **WARNING** level.
