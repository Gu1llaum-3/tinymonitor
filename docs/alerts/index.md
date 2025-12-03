# Alert Providers

TinyMonitor supports multiple alert channels simultaneously. Alerts are sent asynchronously to avoid blocking the monitoring loop.

You can configure multiple providers at the same time. For example, receive critical alerts on your phone via Ntfy and all alerts via Email.

## Available Providers

*   [📡 Ntfy.sh](ntfy.md): Push notifications to mobile/desktop.
*   [🔔 Gotify](gotify.md): Self-hosted push notifications.
*   [💬 Google Chat](google_chat.md): Messages to Google Chat Spaces.
*   [📧 SMTP / Email](smtp.md): Classic email alerts.
*   [🔗 Generic Webhook](webhook.md): Integration with n8n, Zapier, ELK, etc.
