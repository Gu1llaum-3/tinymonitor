# SMTP (Email)

Sends alerts via standard email using SMTP. Supports TLS/SSL encryption.

## Configuration

```json
"smtp": {
    "enabled": true,
    "host": "smtp.gmail.com",
    "port": 587,
    "user": "your_email@gmail.com",
    "password": "your_app_password",
    "from_addr": "monitor@server01.com",
    "to_addrs": ["admin@company.com", "oncall@company.com"],
    "use_tls": true,
    "rules": {
        "default": ["CRITICAL"]
    }
}
```

### Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `enabled` | `bool` | `false` | Enable or disable this provider. |
| `host` | `string` | `""` | SMTP Server Hostname. |
| `port` | `int` | `587` | SMTP Server Port. |
| `user` | `string` | `""` | SMTP Username. |
| `password` | `string` | `""` | SMTP Password (or App Password). |
| `from_addr` | `string` | `""` | Sender email address. |
| `to_addrs` | `list` | `[]` | List of recipient email addresses. |
| `use_tls` | `bool` | `true` | Enable STARTTLS security. |

### Gmail Note
If you are using Gmail, you must enable 2-Step Verification and generate an **App Password**. You cannot use your regular login password.
