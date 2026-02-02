#!/bin/bash
#
# TinyMonitor Installation Script
# https://github.com/Gu1llaum-3/tinymonitor
#
# Supported platforms: Linux, macOS (AMD64, ARM64)
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/Gu1llaum-3/tinymonitor/main/install/install.sh | bash
#
# Environment variables:
#   TINYMONITOR_VERSION  - Version to install (default: latest)
#   INSTALL_DIR          - Installation directory (default: /usr/local/bin)
#

set -euo pipefail

# Configuration
BINARY_NAME="tinymonitor"
GITHUB_REPO="Gu1llaum-3/tinymonitor"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${TINYMONITOR_VERSION:-latest}"
USE_SUDO="false"
OS=""
ARCH=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" >&2
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

# Detect OS and architecture
detect_system() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux|darwin)
            ;;
        *)
            error "Unsupported operating system: $OS. Supported: linux, darwin (macOS)."
            ;;
    esac

    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH. Supported: amd64, arm64."
            ;;
    esac

    # Determine if we need sudo
    if [ "$OS" = "linux" ] || [ "$OS" = "darwin" ]; then
        USE_SUDO="true"
    fi
}

# Run command with sudo if needed
runAsRoot() {
    local CMD="$*"
    if [ "$USE_SUDO" = "true" ] && [ "$(id -u)" != "0" ]; then
        echo -e "${PURPLE}We need sudo access to install to $INSTALL_DIR${NC}" >&2
        CMD="sudo $CMD"
    fi
    $CMD
}

# Get latest version from GitHub
get_latest_version() {
    local latest
    latest=$(curl -sSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [[ -z "$latest" ]]; then
        error "Failed to fetch latest version from GitHub."
    fi

    echo "$latest"
}

# Download and extract binary
download_binary() {
    local os=$1
    local arch=$2
    local version=$3

    # Format OS name (capitalize first letter: linux -> Linux, darwin -> Darwin)
    local os_formatted
    os_formatted="$(echo "${os:0:1}" | tr '[:lower:]' '[:upper:]')${os:1}"

    # Format arch (amd64 -> x86_64, arm64 stays arm64)
    local arch_formatted="$arch"
    if [[ "$arch" == "amd64" ]]; then
        arch_formatted="x86_64"
    fi

    local archive_name="${BINARY_NAME}_${os_formatted}_${arch_formatted}.tar.gz"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${archive_name}"

    info "Downloading ${BINARY_NAME} ${version} for ${os}/${arch}..."

    curl -sSL -o "${BINARY_NAME}-tmp.tar.gz" "$download_url" || error "Failed to download from ${download_url}"

    info "Extracting archive..."
    tar -xzf "${BINARY_NAME}-tmp.tar.gz" || error "Failed to extract archive"

    # Cleanup archive
    rm -f "${BINARY_NAME}-tmp.tar.gz"

    # Check if binary exists
    if [[ ! -f "$BINARY_NAME" ]]; then
        error "Binary not found in archive."
    fi
}

# Install binary
install_binary() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"

    info "Installing to ${install_path}..."

    chmod +x "$BINARY_NAME" || error "Failed to set permissions"
    runAsRoot mv "$BINARY_NAME" "$install_path" || error "Failed to install binary"

    info "Installation complete!"
}

# Verify installation
verify_installation() {
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"

    if [[ -x "$install_path" ]]; then
        info "Verifying installation..."
        "$install_path" version
        return 0
    else
        error "Installation verification failed."
    fi
}

# Cleanup temporary files
cleanup() {
    rm -f "${BINARY_NAME}-tmp.tar.gz" "$BINARY_NAME" 2>/dev/null
    rm -f LICENSE README.md 2>/dev/null
    rm -rf configs install 2>/dev/null
}

# Main function
main() {
    echo ""
    echo "  ╔════════════════════════════════════════╗"
    echo "  ║      TinyMonitor Installation          ║"
    echo "  ╚════════════════════════════════════════╝"
    echo ""

    # Check for required tools
    for cmd in curl tar; do
        if ! command -v "$cmd" &>/dev/null; then
            error "Required command not found: $cmd"
        fi
    done

    # Detect system
    detect_system

    info "Detected system: ${OS}/${ARCH}"

    # Get version
    if [[ "$VERSION" == "latest" ]]; then
        VERSION=$(get_latest_version)
    fi

    info "Version to install: ${VERSION}"

    # Ensure cleanup on exit
    trap cleanup EXIT

    # Download, install and verify
    download_binary "$OS" "$ARCH" "$VERSION"
    install_binary
    verify_installation

    echo ""
    info "TinyMonitor has been installed successfully!"
    echo ""
    echo "  Next steps:"
    echo "    1. Validate your configuration file:"
    echo "       ${BINARY_NAME} validate -c /path/to/config.toml"
    echo ""

    # Show OS-specific service instructions
    if [[ "$OS" == "linux" ]]; then
        echo "    2. Install as a systemd service (optional):"
        echo "       sudo ${BINARY_NAME} service install"
        echo ""
    elif [[ "$OS" == "darwin" ]]; then
        echo "    2. For running as a service on macOS, see:"
        echo "       https://github.com/${GITHUB_REPO}#running-as-a-service-macos-launchd"
        echo ""
    fi

    echo "  Documentation: https://github.com/${GITHUB_REPO}"
    echo ""
}

main "$@"
