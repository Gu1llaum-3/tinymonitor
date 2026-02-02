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

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

# Detect OS
detect_os() {
    local os
    os=$(uname -s | tr '[:upper:]' '[:lower:]')

    case "$os" in
        linux)
            echo "linux"
            ;;
        darwin)
            echo "darwin"
            ;;
        *)
            error "Unsupported operating system: $os. Supported: linux, darwin (macOS)."
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch=$(uname -m)

    case "$arch" in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            error "Unsupported architecture: $arch. Supported: amd64, arm64."
            ;;
    esac
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

# Download binary
download_binary() {
    local os=$1
    local arch=$2
    local version=$3
    local tmp_dir

    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT

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

    if ! curl -sSL -o "${tmp_dir}/${archive_name}" "$download_url"; then
        error "Failed to download from ${download_url}"
    fi

    info "Extracting archive..."
    tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"

    # Find the binary (could be in root or subdirectory)
    local binary_path
    binary_path=$(find "$tmp_dir" -name "$BINARY_NAME" -type f -executable 2>/dev/null | head -n1)

    if [[ -z "$binary_path" ]]; then
        # Try without executable flag (might not be set in archive)
        binary_path=$(find "$tmp_dir" -name "$BINARY_NAME" -type f 2>/dev/null | head -n1)
    fi

    if [[ -z "$binary_path" ]]; then
        error "Binary not found in archive."
    fi

    echo "$binary_path"
}

# Install binary
install_binary() {
    local binary_path=$1
    local install_path="${INSTALL_DIR}/${BINARY_NAME}"

    info "Installing to ${install_path}..."

    # Check if we need sudo
    if [[ -w "$INSTALL_DIR" ]]; then
        cp "$binary_path" "$install_path"
        chmod +x "$install_path"
    else
        warn "Need elevated privileges to install to ${INSTALL_DIR}"
        sudo cp "$binary_path" "$install_path"
        sudo chmod +x "$install_path"
    fi

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
    local os arch
    os=$(detect_os)
    arch=$(detect_arch)

    info "Detected system: ${os}/${arch}"

    # Get version
    if [[ "$VERSION" == "latest" ]]; then
        VERSION=$(get_latest_version)
    fi

    info "Version to install: ${VERSION}"

    # Download and install
    local binary_path
    binary_path=$(download_binary "$os" "$arch" "$VERSION")
    install_binary "$binary_path"

    # Verify
    verify_installation

    echo ""
    info "TinyMonitor has been installed successfully!"
    echo ""
    echo "  Next steps:"
    echo "    1. Validate your configuration file:"
    echo "       ${BINARY_NAME} validate -c /path/to/config.toml"
    echo ""

    # Show OS-specific service instructions
    if [[ "$os" == "linux" ]]; then
        echo "    2. Install as a systemd service (optional):"
        echo "       sudo ${BINARY_NAME} service install"
        echo ""
    elif [[ "$os" == "darwin" ]]; then
        echo "    2. For running as a service on macOS, see:"
        echo "       https://github.com/${GITHUB_REPO}#running-as-a-service-macos-launchd"
        echo ""
    fi

    echo "  Documentation: https://github.com/${GITHUB_REPO}"
    echo ""
}


main "$@"
