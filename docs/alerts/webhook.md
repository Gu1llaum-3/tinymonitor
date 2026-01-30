# Generic Webhook

The generic webhook allows you to send alerts to any external system (n8n, Zapier, ELK, Custom Dashboard) via a raw JSON payload.

## Configuration

```toml
[alerts.webhook]
enabled = true
url = "https://your-endpoint.com/api/alert"
timeout = 10

  [alerts.webhook.headers]
  Authorization = "Bearer your-secret-token"
  X-Custom-Header = "TinyMonitor"

  [alerts.webhook.rules]
  default = ["WARNING", "CRITICAL"]
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `false` | Enable or disable this provider. |
| `url` | `string` | `""` | The target URL (POST request). |
| `headers` | `table` | `{}` | Custom HTTP headers to include. |
| `timeout` | `int` | `10` | Request timeout in seconds. |

## Payload Format

TinyMonitor sends a POST request with the following JSON body:

```json
{
  "timestamp": "2025-12-03T14:30:00.123456",
  "alert": {
    "level": "CRITICAL",
    "component": "cpu",
    "value": "95.5",
    "title": "ALERT CRITICAL : cpu",
    "message": "Component cpu is in state CRITICAL. Value: 95.5"
  },
  "host": {
    "hostname": "prod-server-01",
    "ip_private": "192.168.1.10",
    "ip_public": "203.0.113.42",
    "uptime": "12 days, 4:02",
    "load_average": "0.5 0.4 0.3"
  }
}
```
