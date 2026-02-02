# tinymonitor test-alert

Send a test alert to verify your notification providers are working correctly.

## Usage

```bash
tinymonitor test-alert [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file |
| `--provider <name>` | Test only this specific provider |

## Examples

```bash
# Test all enabled providers
tinymonitor test-alert

# Test specific provider
tinymonitor test-alert --provider ntfy
tinymonitor test-alert --provider smtp
tinymonitor test-alert --provider gotify
```

## Available Providers

| Provider | Name for `--provider` |
|----------|----------------------|
| Ntfy.sh | `ntfy` |
| Gotify | `gotify` |
| Google Chat | `google_chat` |
| Email (SMTP) | `smtp` |
| Webhook | `webhook` |

## Test Alert Content

The test alert includes:

- **Level**: INFO
- **Message**: "This is a test alert from TinyMonitor"
- **Hostname**: Your server's hostname
- **Timestamp**: Current date and time

## Example Output

```bash
$ tinymonitor test-alert

Testing alert providers...

  ntfy:        OK
  smtp:        OK
  google_chat: FAILED - connection timeout

1 provider(s) failed. Check your configuration.
```

## Troubleshooting

If a provider fails:

1. Check the provider configuration in your config file
2. Verify network connectivity to the service
3. Check authentication credentials
4. Review [Troubleshooting](../guides/troubleshooting.md) for common issues

## See Also

- [Alerts Overview](../alerts/index.md)
- [Troubleshooting](../guides/troubleshooting.md)
