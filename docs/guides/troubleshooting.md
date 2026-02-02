# Troubleshooting

Common issues and their solutions.

## Alerts Not Received

### 1. Verify configuration

```bash
tinymonitor validate -c /etc/tinymonitor/config.toml
```

### 2. Test alert providers

```bash
# Test all providers
tinymonitor test-alert

# Test specific provider
tinymonitor test-alert --provider ntfy
```

### 3. Check provider-specific issues

#### Ntfy

- Verify the topic URL is correct and accessible
- If using a private topic, ensure the token is valid
- Check if the ntfy server is reachable

#### SMTP / Email

- **Gmail**: Requires [App Password](https://support.google.com/accounts/answer/185833), not your regular password
- Check spam/junk folder
- Verify port: 587 for TLS, 465 for SSL
- Ensure `use_tls = true` for secure connections

#### Gotify

- Use the **Application Token**, not Client Token
- Verify the URL ends without trailing slash
- Check if the Gotify server is accessible

#### Webhook

- Verify the URL is accessible from your server
- Check required headers (Authorization, etc.)
- Test with curl: `curl -X POST -H "Content-Type: application/json" -d '{}' YOUR_URL`

#### Google Chat

- Webhook URLs expire if not used
- Regenerate the webhook in Google Chat settings

## Service Won't Start

### Check logs

```bash
# systemd (Linux)
sudo journalctl -u tinymonitor -n 50

# Direct run to see errors
tinymonitor -c /etc/tinymonitor/config.toml
```

### Common causes

#### Invalid configuration

```bash
tinymonitor validate -c /etc/tinymonitor/config.toml
```

#### Permission issues

```bash
# Check config file permissions
ls -la /etc/tinymonitor/config.toml

# Should be readable (at least -r--r--r--)
sudo chmod 644 /etc/tinymonitor/config.toml
```

#### Binary not found

```bash
which tinymonitor
# Should return: /usr/local/bin/tinymonitor

# If not found, reinstall
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
```

## Configuration Errors

### TOML syntax errors

Common mistakes:

```toml
# Wrong: missing quotes for strings with special characters
topic_url = https://ntfy.sh/topic  # Error!

# Correct:
topic_url = "https://ntfy.sh/topic"
```

```toml
# Wrong: using = in section headers
[alerts.ntfy = true]  # Error!

# Correct:
[alerts.ntfy]
enabled = true
```

### Invalid threshold values

```toml
# Wrong: percentage over 100
warning = 150  # Error!

# Correct: must be 0-100 for percentages
warning = 80
```

### Missing required fields

When a provider is enabled, certain fields are required:

```toml
[alerts.smtp]
enabled = true
# Error: missing host, port, user, password, from_addr, to_addrs
```

## Metric Not Reporting

### Load Average on macOS

Load average works on macOS but values are interpreted differently than Linux. A load of 4.0 on a 4-core Mac is equivalent to 100% utilization.

### Reboot Required

This metric only works on **Debian/Ubuntu** and derivatives. It checks for the file `/var/run/reboot-required`.

On other systems, this metric will always report "OK".

### Filesystem Exclusions

Some paths are automatically excluded:

- `/snap/*` (Snap packages)
- Docker overlay filesystems
- Virtual filesystems (`/proc`, `/sys`, `/dev`)

To exclude additional paths:

```toml
[filesystem]
exclude = ["/mnt/backup", "/media/external"]
```

### Disk I/O

If I/O metrics seem wrong:

- Ensure the system has actual disk activity
- Check if thresholds are set correctly (bytes vs. MB)

```toml
[io]
warning = "50MB"    # 50 megabytes per second
critical = "100MB"  # 100 megabytes per second
```

## High CPU/Memory Usage

TinyMonitor is designed to be lightweight (<10MB RAM, <1% CPU). If you see high usage:

1. **Increase refresh interval** (minimum recommended: 5 seconds)
   ```toml
   refresh = 5
   ```

2. **Disable unused metrics**
   ```toml
   [io]
   enabled = false
   ```

3. **Check for disk issues** - slow disks can cause delays

## Update Issues

### Permission denied

```bash
sudo tinymonitor update
```

### Network issues

Check connectivity to GitHub:

```bash
curl -I https://api.github.com/repos/Gu1llaum-3/tinymonitor/releases/latest
```

### Manual update

If automatic update fails:

```bash
curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
```

### After updating

Don't forget to restart the service:

```bash
sudo systemctl restart tinymonitor
```

## Getting Help

If you can't resolve your issue:

1. Search [existing issues](https://github.com/Gu1llaum-3/tinymonitor/issues)
2. Open a [new issue](https://github.com/Gu1llaum-3/tinymonitor/issues/new) with:
   - TinyMonitor version: `tinymonitor version`
   - OS and architecture: `uname -a`
   - Relevant log output
   - Configuration (remove sensitive data like tokens/passwords)
