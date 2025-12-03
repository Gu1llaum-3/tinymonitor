# Gotify

[Gotify](https://gotify.net) is a simple server for sending and receiving messages. It is self-hosted and open source.

## Configuration

```json
"gotify": {
    "enabled": true,
    "url": "https://gotify.yourdomain.com",
    "token": "A1b2C3d4E5f6G7h",
    "rules": {
        "default": ["WARNING", "CRITICAL"]
    }
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `false` | Enable or disable this provider. |
| `url` | `string` | `""` | The base URL of your Gotify server. |
| `token` | `string` | `""` | The Application Token (not the client token). |
| `rules` | `object` | `{}` | Alert filtering rules. |

### Features
*   **Markdown**: Messages are formatted using Markdown for better readability.
*   **Priorities**: Maps `CRITICAL` to Priority 8 and `WARNING` to Priority 5.
