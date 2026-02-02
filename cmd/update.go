package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Gu1llaum-3/tinymonitor/internal/system"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update TinyMonitor to the latest version",
	Long: `Check for updates and install the latest version of TinyMonitor.

The configuration file in /etc/tinymonitor/ is never modified.

Examples:
  tinymonitor update              # Interactive update
  tinymonitor update --check      # Check only, don't install
  tinymonitor update --yes        # Update without confirmation`,
	Run: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().Bool("check", false, "Check for updates without installing")
	updateCmd.Flags().BoolP("yes", "y", false, "Update without confirmation")
}

func runUpdate(cmd *cobra.Command, args []string) {
	checkOnly, _ := cmd.Flags().GetBool("check")
	yes, _ := cmd.Flags().GetBool("yes")

	// Get current version
	currentVersion := Version

	// Fetch latest version
	fmt.Println("Checking for updates...")
	latestVersion, err := system.GetLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("Current version: %s\n", currentVersion)
	fmt.Printf("Latest version:  %s\n", latestVersion)
	fmt.Println()

	// Check if update needed
	if !system.CompareVersions(currentVersion, latestVersion) {
		fmt.Printf("Already up to date (%s)!\n", currentVersion)
		return
	}

	fmt.Println("A new version is available!")
	fmt.Printf("\nChangelog: %s\n\n", system.GetChangelogURL(latestVersion))

	// Check-only mode
	if checkOnly {
		return
	}

	// Ask confirmation unless --yes is provided
	if !yes {
		fmt.Print("Do you want to update? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Update cancelled.")
			return
		}
	}

	// Download and install
	fmt.Printf("Downloading tinymonitor %s...\n", latestVersion)
	binaryPath, err := system.DownloadBinary(latestVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading: %v\n", err)
		os.Exit(1)
	}
	defer system.CleanupUpdate(binaryPath)

	fmt.Printf("Installing to %s...\n", system.BinaryPath)
	if err := system.InstallBinary(binaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error installing: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nUpdate complete!")

	// Remind about service restart if running
	if system.IsServiceRunning() {
		fmt.Println()
		fmt.Println("Note: TinyMonitor service is running.")
		fmt.Println("Run 'sudo systemctl restart tinymonitor' to apply the update.")
	}
}
