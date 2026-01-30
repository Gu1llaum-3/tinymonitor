# Google Chat

Sends notifications to a Google Chat Space using Incoming Webhooks.

## Configuration

```toml
[alerts.google_chat]
enabled = true
webhook_url = "https://chat.googleapis.com/v1/spaces/..."

  [alerts.google_chat.rules]
  default = ["WARNING", "CRITICAL"]
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `false` | Enable or disable this provider. |
| `webhook_url` | `string` | `""` | The Google Chat Incoming Webhook URL. |
| `rules` | `table` | `{}` | Alert filtering rules. |

### Setup

1.  Go to Google Chat.
2.  Select the Space where you want to receive alerts.
3.  Click on the Space name > **Apps & integrations**.
4.  Click **Manage webhooks**.
5.  Create a new webhook and copy the URL.
