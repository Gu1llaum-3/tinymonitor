# Running on macOS (launchd)

This guide explains how to run TinyMonitor as a background service on macOS using launchd.

## Create the Launch Agent

Create a plist file at `~/Library/LaunchAgents/com.tinymonitor.plist`:

```bash
nano ~/Library/LaunchAgents/com.tinymonitor.plist
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tinymonitor</string>

    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/tinymonitor</string>
        <string>-c</string>
        <string>/Users/YOUR_USERNAME/.config/tinymonitor/config.toml</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>/tmp/tinymonitor.log</string>

    <key>StandardErrorPath</key>
    <string>/tmp/tinymonitor.err</string>
</dict>
</plist>
```

**Important**: Replace `YOUR_USERNAME` with your actual username.

## Load the Service

```bash
launchctl load ~/Library/LaunchAgents/com.tinymonitor.plist
```

The service will now:

- Start immediately
- Start automatically on login
- Restart if it crashes

## Management Commands

### Check if running

```bash
launchctl list | grep tinymonitor
```

If running, you'll see output like:

```
12345   0   com.tinymonitor
```

The first number is the PID, the second is the exit code (0 = running).

### View logs

```bash
# Standard output
tail -f /tmp/tinymonitor.log

# Error output
tail -f /tmp/tinymonitor.err
```

### Stop the service

```bash
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist
```

### Restart the service

```bash
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist
launchctl load ~/Library/LaunchAgents/com.tinymonitor.plist
```

## Configuration Location

On macOS, the recommended configuration location is:

```
~/.config/tinymonitor/config.toml
```

The install script creates this automatically.

## Running as a System Service

To run TinyMonitor for all users (requires admin privileges):

1. Create the plist in `/Library/LaunchDaemons/` instead
2. Use absolute paths (not `~`)
3. Load with `sudo launchctl load`

```bash
sudo nano /Library/LaunchDaemons/com.tinymonitor.plist
```

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tinymonitor</string>

    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/tinymonitor</string>
        <string>-c</string>
        <string>/etc/tinymonitor/config.toml</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

```bash
sudo launchctl load /Library/LaunchDaemons/com.tinymonitor.plist
```

## Uninstall

```bash
# Unload the service
launchctl unload ~/Library/LaunchAgents/com.tinymonitor.plist

# Remove the plist
rm ~/Library/LaunchAgents/com.tinymonitor.plist

# Remove the binary (optional)
sudo rm /usr/local/bin/tinymonitor

# Remove configuration (optional)
rm -rf ~/.config/tinymonitor
```

## Troubleshooting

### Service doesn't start

Check the error log:

```bash
cat /tmp/tinymonitor.err
```

### Permission denied

Ensure the config file is readable:

```bash
chmod 644 ~/.config/tinymonitor/config.toml
```

### launchctl error on load

Validate the plist syntax:

```bash
plutil -lint ~/Library/LaunchAgents/com.tinymonitor.plist
```

## See Also

- [Running on Linux](systemd.md) - systemd setup
- [Troubleshooting](troubleshooting.md) - Common issues
