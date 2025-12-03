# Development Guide

Want to contribute? Here is how to build and test TinyMonitor.

## Prerequisites
*   Python 3.8+
*   Make
*   Virtualenv (recommended)

## Setup Environment

1.  Clone the repository:
    ```bash
    git clone https://github.com/Gu1llaum-3/tinymonitor.git
    cd tinymonitor
    ```

2.  Create a virtual environment:
    ```bash
    python3 -m venv .venv
    source .venv/bin/activate
    ```

3.  Install development dependencies:
    ```bash
    pip install -r requirements-dev.txt
    ```

## Commands

We use a `Makefile` to automate tasks:

| Command | Description |
| :--- | :--- |
| `make install` | Install dependencies and dev tools. |
| `make test` | Run unit tests. |
| `make build` | Compile the standalone binary to `dist/`. |
| `make clean` | Remove build artifacts. |

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration:

1.  **Tests**: Runs on every push to `main`.
2.  **Build**: Runs on every Tag (`v*`). Builds binaries for Linux (AMD64/ARM64) and macOS (Intel/Silicon).
3.  **Release**: Automatically creates a GitHub Release with the binaries attached.
