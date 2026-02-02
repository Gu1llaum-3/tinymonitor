# tinymonitor version

Print version, commit hash, and build date.

## Usage

```bash
tinymonitor version
```

## Output

```
TinyMonitor v1.2.0
Commit: abc1234
Built: 2024-01-15T10:30:00Z
```

## Fields

| Field | Description |
|-------|-------------|
| Version | Semantic version (e.g., `v1.2.0`) |
| Commit | Git commit hash of the build |
| Built | Build timestamp in ISO 8601 format |

## Development Builds

When building from source without version flags:

```
TinyMonitor dev
Commit: none
Built: unknown
```

## Checking for Updates

To check if a newer version is available:

```bash
tinymonitor update --check
```

## See Also

- [tinymonitor update](update.md)
- [GitHub Releases](https://github.com/Gu1llaum-3/tinymonitor/releases)
