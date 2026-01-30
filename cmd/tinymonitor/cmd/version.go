package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags)
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, commit hash, and build date of TinyMonitor.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TinyMonitor %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
