package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Gu1llaum-3/tinymonitor/internal/system"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall TinyMonitor completely",
	Long: `Completely uninstall TinyMonitor from the system.

This command will:
  1. Stop and remove the systemd service (if installed)
  2. Remove the binary from /usr/local/bin

Options:
  --purge  Also remove configuration files from /etc/tinymonitor/
  --yes    Skip confirmation prompt

Examples:
  sudo tinymonitor uninstall           # Keep configuration
  sudo tinymonitor uninstall --purge   # Remove everything
  sudo tinymonitor uninstall --yes     # Skip confirmation`,
	Run: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().Bool("purge", false, "Remove configuration files as well")
	uninstallCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
}

func runUninstall(cmd *cobra.Command, args []string) {
	// Check root
	if !system.IsRoot() {
		fmt.Fprintln(os.Stderr, "Error: This command requires root privileges.")
		fmt.Fprintln(os.Stderr, "Run with: sudo tinymonitor uninstall")
		os.Exit(1)
	}

	purge, _ := cmd.Flags().GetBool("purge")
	yes, _ := cmd.Flags().GetBool("yes")

	// Check if there's anything to uninstall
	if !system.IsBinaryInstalled() && !system.IsServiceInstalled() && !system.ConfigExists() {
		fmt.Println("TinyMonitor is not installed on this system.")
		return
	}

	// Show what will be removed
	fmt.Println("TinyMonitor Uninstallation")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Print(system.UninstallSummary(purge))
	fmt.Println()

	// Ask for confirmation unless --yes is provided
	if !yes {
		fmt.Print("Are you sure you want to continue? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Uninstallation cancelled.")
			return
		}
	}

	fmt.Println()
	fmt.Println("Uninstalling TinyMonitor...")

	opts := system.UninstallOptions{
		Purge: purge,
	}

	if err := system.Uninstall(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("TinyMonitor has been uninstalled successfully!")

	if !purge && system.ConfigExists() {
		fmt.Printf("\nConfiguration files were preserved at %s\n", system.DefaultConfigDir)
		fmt.Println("Run with --purge to remove them as well.")
	}
}
