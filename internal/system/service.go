package system

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Paths for service installation
const (
	DefaultConfigDir  = "/etc/tinymonitor"
	DefaultConfigFile = "/etc/tinymonitor/config.toml"
	ServiceFile       = "/etc/systemd/system/tinymonitor.service"
	BinaryPath        = "/usr/local/bin/tinymonitor"
	ServiceName       = "tinymonitor"
)

// ServiceConfig contains configuration for the systemd service
type ServiceConfig struct {
	ConfigPath string // Path to config.toml
	User       string // User to run service as (default: nobody)
	Group      string // Group to run service as (default: nogroup)
}

// systemd service template
const serviceTemplate = `[Unit]
Description=TinyMonitor - Lightweight System Monitoring
Documentation=https://github.com/Gu1llaum-3/tinymonitor
After=network.target

[Service]
Type=simple
User={{.User}}
Group={{.Group}}
ExecStart={{.BinaryPath}} -c {{.ConfigPath}}
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=read-only
PrivateTmp=true
ReadOnlyPaths=/

[Install]
WantedBy=multi-user.target
`

// IsRoot checks if the current user is root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// IsSystemd checks if systemd is available on the system
func IsSystemd() bool {
	_, err := exec.LookPath("systemctl")
	if err != nil {
		return false
	}

	// Check if systemd is the init system
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return true
	}

	return false
}

// IsBinaryInstalled checks if tinymonitor is installed in the expected location
func IsBinaryInstalled() bool {
	_, err := os.Stat(BinaryPath)
	return err == nil
}

// IsServiceInstalled checks if the systemd service exists
func IsServiceInstalled() bool {
	_, err := os.Stat(ServiceFile)
	return err == nil
}

// InstallService installs and starts the systemd service
func InstallService(cfg ServiceConfig) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required to install service")
	}

	if !IsSystemd() {
		return fmt.Errorf("systemd is not available on this system")
	}

	if !IsBinaryInstalled() {
		return fmt.Errorf("tinymonitor binary not found at %s", BinaryPath)
	}

	// Set defaults
	if cfg.User == "" {
		cfg.User = "nobody"
	}
	if cfg.Group == "" {
		cfg.Group = "nogroup"
	}
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = DefaultConfigFile
	}

	// Create config directory if needed
	if err := os.MkdirAll(DefaultConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Generate service file content
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	data := struct {
		User       string
		Group      string
		BinaryPath string
		ConfigPath string
	}{
		User:       cfg.User,
		Group:      cfg.Group,
		BinaryPath: BinaryPath,
		ConfigPath: cfg.ConfigPath,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to generate service file: %w", err)
	}

	// Write service file
	if err := os.WriteFile(ServiceFile, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	if err := runSystemctl("daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	// Enable service
	if err := runSystemctl("enable", ServiceName); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	// Start service
	if err := runSystemctl("start", ServiceName); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// UninstallService stops and removes the systemd service
func UninstallService() error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required to uninstall service")
	}

	if !IsServiceInstalled() {
		return fmt.Errorf("service is not installed")
	}

	// Stop service (ignore errors if not running)
	_ = runSystemctl("stop", ServiceName)

	// Disable service
	_ = runSystemctl("disable", ServiceName)

	// Remove service file
	if err := os.Remove(ServiceFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	if err := runSystemctl("daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	return nil
}

// ServiceStatus returns the current status of the service
func ServiceStatus() (string, error) {
	if !IsSystemd() {
		return "", fmt.Errorf("systemd is not available")
	}

	if !IsServiceInstalled() {
		return "not installed", nil
	}

	cmd := exec.Command("systemctl", "status", ServiceName)
	output, err := cmd.CombinedOutput()

	// systemctl status returns non-zero exit code if service is not running
	// We still want to show the output
	return string(output), err
}

// IsServiceRunning checks if the service is currently running
func IsServiceRunning() bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", ServiceName)
	return cmd.Run() == nil
}

// runSystemctl executes a systemctl command
func runSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CopyConfig copies a config file to the default location
func CopyConfig(sourcePath string) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	// Ensure config directory exists
	if err := os.MkdirAll(DefaultConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source config: %w", err)
	}

	// Write to default location
	if err := os.WriteFile(DefaultConfigFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// WriteDefaultConfig writes the embedded default config to the default location
func WriteDefaultConfig(content []byte) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required")
	}

	// Ensure config directory exists
	if err := os.MkdirAll(DefaultConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(DefaultConfigFile); err == nil {
		return fmt.Errorf("config file already exists at %s", DefaultConfigFile)
	}

	// Write default config
	if err := os.WriteFile(DefaultConfigFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ConfigExists checks if the default config file exists
func ConfigExists() bool {
	_, err := os.Stat(DefaultConfigFile)
	return err == nil
}

// GetConfigPath returns the path to use for the service config
func GetConfigPath(customPath string) string {
	if customPath != "" {
		if filepath.IsAbs(customPath) {
			return customPath
		}
		// Convert relative path to absolute
		abs, err := filepath.Abs(customPath)
		if err == nil {
			return abs
		}
	}
	return DefaultConfigFile
}

// ValidateConfigPath checks if a config file exists and is readable
func ValidateConfigPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found: %s", path)
		}
		return fmt.Errorf("cannot access config file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("config path is a directory: %s", path)
	}

	// Try to read the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}
	file.Close()

	return nil
}

// FormatServiceStatus formats the service status for display
func FormatServiceStatus() string {
	var sb strings.Builder

	sb.WriteString("TinyMonitor Service Status\n")
	sb.WriteString("==========================\n\n")

	// Binary
	if IsBinaryInstalled() {
		sb.WriteString(fmt.Sprintf("Binary:     %s (installed)\n", BinaryPath))
	} else {
		sb.WriteString(fmt.Sprintf("Binary:     %s (not found)\n", BinaryPath))
	}

	// Config
	if ConfigExists() {
		sb.WriteString(fmt.Sprintf("Config:     %s (exists)\n", DefaultConfigFile))
	} else {
		sb.WriteString(fmt.Sprintf("Config:     %s (not found)\n", DefaultConfigFile))
	}

	// Service
	if !IsSystemd() {
		sb.WriteString("Service:    systemd not available\n")
	} else if !IsServiceInstalled() {
		sb.WriteString("Service:    not installed\n")
	} else if IsServiceRunning() {
		sb.WriteString("Service:    running\n")
	} else {
		sb.WriteString("Service:    stopped\n")
	}

	return sb.String()
}
