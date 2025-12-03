# Ntfy.sh

[Ntfy.sh](https://ntfy.sh) is a simple HTTP-based pub-sub notification service. It allows you to send notifications to your phone or desktop via scripts from any computer.

## Configuration

```json
"ntfy": {
    "enabled": true,
    "topic_url": "https://ntfy.sh/my_secret_topic",
    "token": "optional_token",
    "rules": {
        "default": ["WARNING", "CRITICAL"]
    }
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `false` | Enable or disable this provider. |
| `topic_url` | `string` | `""` | The full URL of your topic (e.g., `https://ntfy.sh/mytopic`). |
| `token` | `string` | `""` | Optional access token if your topic is protected. |
| `rules` | `object` | `{}` | Alert filtering rules. |

### Features
*   **Priorities**: Maps `CRITICAL` to High Priority and `WARNING` to Default Priority.
*   **Tags**: Adds emojis based on the alert level.
*   **Markdown**: Supports basic Markdown formatting in messages.
