# tinymonitor validate

Validate your configuration file before deployment.

## Usage

```bash
tinymonitor validate [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `-c, --config <path>` | Path to configuration file |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Configuration is valid |
| 1 | Configuration has errors |

## Examples

### Valid configuration

```bash
$ tinymonitor validate -c config.toml
Configuration is valid.
```

### Invalid configuration

```bash
$ tinymonitor validate -c broken.toml
Configuration errors:
  - cpu.warning: must be between 0 and 100
  - alerts.smtp.host: required when smtp is enabled
  - alerts.ntfy.topic_url: invalid URL format

Run 'tinymonitor validate -c <file>' for details.
```

## What Gets Validated

- **Syntax**: TOML parsing errors
- **Required fields**: missing mandatory values
- **Value ranges**: thresholds within valid bounds (0-100 for percentages)
- **Type checking**: correct data types for each field
- **Provider configuration**: required fields when a provider is enabled
- **URL formats**: valid URLs for webhooks and API endpoints

## Use Cases

- **CI/CD pipelines**: validate config before deployment
- **Pre-flight checks**: ensure config is valid before starting
- **Debugging**: identify configuration errors

## See Also

- [Configuration Reference](../configuration.md)
- [tinymonitor info](info.md)
