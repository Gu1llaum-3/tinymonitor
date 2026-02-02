package system

import (
	"fmt"
	"os"
)

// UninstallOptions contains options for uninstallation
type UninstallOptions struct {
	Purge bool // Remove configuration files as well
}

// Uninstall removes TinyMonitor from the system
func Uninstall(opts UninstallOptions) error {
	if !IsRoot() {
		return fmt.Errorf("root privileges required to uninstall")
	}

	var errors []string

	// 1. Stop and remove service if installed
	if IsServiceInstalled() {
		if err := UninstallService(); err != nil {
			errors = append(errors, fmt.Sprintf("service: %v", err))
		}
	}

	// 2. Remove binary
	if IsBinaryInstalled() {
		if err := os.Remove(BinaryPath); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("binary: %v", err))
		}
	}

	// 3. Remove config directory if purge is enabled
	if opts.Purge {
		if err := os.RemoveAll(DefaultConfigDir); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("config: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("uninstallation completed with errors: %v", errors)
	}

	return nil
}

// UninstallSummary returns a summary of what will be removed
func UninstallSummary(purge bool) string {
	var sb string

	sb += "The following will be removed:\n\n"

	if IsServiceInstalled() {
		sb += fmt.Sprintf("  - Systemd service: %s\n", ServiceFile)
	}

	if IsBinaryInstalled() {
		sb += fmt.Sprintf("  - Binary: %s\n", BinaryPath)
	}

	if purge && ConfigExists() {
		sb += fmt.Sprintf("  - Config directory: %s\n", DefaultConfigDir)
	} else if ConfigExists() {
		sb += fmt.Sprintf("\n  Config will be kept at: %s\n", DefaultConfigDir)
		sb += "  (use --purge to remove)\n"
	}

	return sb
}
