# Development Guide

Want to contribute? Here is how to build and test TinyMonitor.

## Prerequisites

*   Go 1.21+
*   Make (optional but recommended)

## Setup Environment

1.  Clone the repository:
    ```bash
    git clone https://github.com/Gu1llaum-3/tinymonitor.git
    cd tinymonitor
    ```

2.  Build the project:
    ```bash
    make build
    # Or directly with Go
    go build -o tinymonitor ./cmd/tinymonitor
    ```

3.  Run tests:
    ```bash
    make test
    # Or directly with Go
    go test ./...
    ```

## Project Structure

```
tinymonitor/
├── cmd/tinymonitor/
│   ├── main.go             # Entry point
│   └── cmd/                 # CLI commands (Cobra)
│       ├── root.go         # Main monitoring command
│       ├── version.go      # Version command
│       └── validate.go     # Config validation command
├── internal/
│   ├── config/             # TOML configuration loading & validation
│   ├── monitor/            # Main monitoring loop
│   ├── metrics/            # Metric collectors (CPU, memory, disk, etc.)
│   ├── alerts/             # Alert providers (Ntfy, SMTP, etc.)
│   ├── models/             # Shared types
│   └── utils/              # Utility functions
├── configs/                # Example configurations
├── docs/                   # Documentation
├── go.mod
└── Makefile
```

## Commands

We use a `Makefile` to automate tasks:

| Command | Description |
| :--- | :--- |
| `make build` | Build the binary. |
| `make test` | Run unit tests. |
| `make vet` | Run static analysis. |
| `make release` | Build binaries for all platforms. |
| `make clean` | Remove build artifacts. |

## CLI Commands

TinyMonitor uses Cobra for CLI management:

```bash
tinymonitor              # Run the monitoring agent
tinymonitor version      # Show version information
tinymonitor validate     # Validate configuration file
tinymonitor -c config.toml  # Run with specific config
```

## Adding a New Metric

1.  Create a new file in `internal/metrics/` (e.g., `network.go`)
2.  Implement the `Collector` interface:
    ```go
    type Collector interface {
        Name() string
        Check() []models.MetricResult
        Duration() int
    }
    ```
3.  Register the collector in `internal/monitor/monitor.go`

## Adding a New Alert Provider

1.  Create a new file in `internal/alerts/` (e.g., `slack.go`)
2.  Implement the `Provider` interface:
    ```go
    type Provider interface {
        Name() string
        Send(alert models.Alert) error
        ShouldSend(component string, level models.Severity) bool
    }
    ```
3.  Register the provider in `internal/alerts/manager.go`

## Configuration Validation

When adding new configuration fields, update the validation in `internal/config/config.go`:

```go
func (c *Config) Validate() ValidationErrors {
    var errs ValidationErrors
    // Add your validation logic here
    return errs
}
```

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration:

1.  **Tests**: Runs on every push to `main`.
2.  **Build**: Runs on every Tag (`v*`). Uses GoReleaser to build binaries for Linux (AMD64/ARM64) and macOS (Intel/Silicon).
3.  **Release**: Automatically creates a GitHub Release with the binaries attached.
