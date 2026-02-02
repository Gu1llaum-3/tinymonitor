package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/system"
	"github.com/spf13/cobra"
)

//go:embed config.example.toml
var defaultConfigContent []byte

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage TinyMonitor systemd service",
	Long: `Manage the TinyMonitor systemd service.

Available subcommands:
  install    - Install and start the systemd service
  uninstall  - Stop and remove the systemd service
  status     - Show service status`,
}

var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and start the systemd service",
	Long: `Install TinyMonitor as a systemd service.

This command will:
  1. Create /etc/tinymonitor/ directory if needed
  2. Copy the configuration file (or create default)
  3. Create the systemd service file
  4. Enable and start the service

Examples:
  sudo tinymonitor service install                    # Use default config
  sudo tinymonitor service install -c myconfig.toml  # Use custom config`,
	Run: runServiceInstall,
}

var serviceUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Stop and remove the systemd service",
	Long: `Stop and remove the TinyMonitor systemd service.

This command will:
  1. Stop the service if running
  2. Disable the service
  3. Remove the service file

Note: Configuration files in /etc/tinymonitor/ are preserved.
Use 'tinymonitor uninstall --purge' to remove everything.`,
	Run: runServiceUninstall,
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show service status",
	Long:  `Display the current status of the TinyMonitor systemd service.`,
	Run:   runServiceStatus,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)
	serviceCmd.AddCommand(serviceStatusCmd)

	serviceInstallCmd.Flags().StringP("config", "c", "", "Path to configuration file to use")
}

func runServiceInstall(cmd *cobra.Command, args []string) {
	// Check root
	if !system.IsRoot() {
		fmt.Fprintln(os.Stderr, "Error: This command requires root privileges.")
		fmt.Fprintln(os.Stderr, "Run with: sudo tinymonitor service install")
		os.Exit(1)
	}

	// Check systemd
	if !system.IsSystemd() {
		fmt.Fprintln(os.Stderr, "Error: systemd is not available on this system.")
		os.Exit(1)
	}

	// Check binary
	if !system.IsBinaryInstalled() {
		fmt.Fprintf(os.Stderr, "Error: TinyMonitor binary not found at %s\n", system.BinaryPath)
		fmt.Fprintln(os.Stderr, "Please install the binary first.")
		os.Exit(1)
	}

	configPath, _ := cmd.Flags().GetString("config")

	// Handle configuration
	if configPath != "" {
		// Validate and copy custom config
		if err := system.ValidateConfigPath(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Validate the config syntax
		if _, err := config.LoadAndValidate(configPath); err != nil {
			fmt.Fprintln(os.Stderr, "Error: Invalid configuration file:")
			fmt.Fprintf(os.Stderr, "  %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Copying configuration from %s...\n", configPath)
		if err := system.CopyConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else if !system.ConfigExists() {
		// Create default config
		fmt.Println("Creating default configuration...")
		if err := system.WriteDefaultConfig(defaultConfigContent); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Default configuration written to %s\n", system.DefaultConfigFile)
		fmt.Println("Please edit this file to configure your monitoring settings.")
	} else {
		fmt.Printf("Using existing configuration at %s\n", system.DefaultConfigFile)
	}

	// Install service
	fmt.Println("Installing systemd service...")
	cfg := system.ServiceConfig{
		ConfigPath: system.DefaultConfigFile,
	}

	if err := system.InstallService(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("TinyMonitor service installed successfully!")
	fmt.Println()
	fmt.Println("Useful commands:")
	fmt.Println("  sudo systemctl status tinymonitor   # Check status")
	fmt.Println("  sudo systemctl restart tinymonitor  # Restart service")
	fmt.Println("  sudo journalctl -u tinymonitor -f   # View logs")
}

func runServiceUninstall(cmd *cobra.Command, args []string) {
	// Check root
	if !system.IsRoot() {
		fmt.Fprintln(os.Stderr, "Error: This command requires root privileges.")
		fmt.Fprintln(os.Stderr, "Run with: sudo tinymonitor service uninstall")
		os.Exit(1)
	}

	// Check if service is installed
	if !system.IsServiceInstalled() {
		fmt.Println("Service is not installed.")
		return
	}

	fmt.Println("Uninstalling TinyMonitor service...")

	if err := system.UninstallService(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("TinyMonitor service uninstalled successfully!")
	fmt.Printf("Configuration files in %s were preserved.\n", system.DefaultConfigDir)
}

func runServiceStatus(cmd *cobra.Command, args []string) {
	fmt.Print(system.FormatServiceStatus())

	// Show detailed systemd status if available
	if system.IsSystemd() && system.IsServiceInstalled() {
		fmt.Println()
		fmt.Println("Systemd Status:")
		fmt.Println("---------------")
		status, _ := system.ServiceStatus()
		fmt.Println(status)
	}
}
