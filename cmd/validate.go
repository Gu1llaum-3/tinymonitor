package cmd

import (
	"fmt"
	"os"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long: `Validate a TinyMonitor configuration file.

Checks for:
  - Valid TOML syntax
  - Required fields when providers are enabled
  - Threshold values (0-100 for percentages)
  - warning < critical for all metrics
  - Valid port numbers

Examples:
  tinymonitor validate
  tinymonitor validate -c /etc/tinymonitor/config.toml`,
	Run: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	configPath := cfgFile

	// If no config specified, try to find one
	if configPath == "" {
		searchPaths := []string{
			"config.toml",
			os.ExpandEnv("$HOME/.config/tinymonitor/config.toml"),
			"/etc/tinymonitor/config.toml",
		}

		for _, path := range searchPaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}

		if configPath == "" {
			fmt.Println("No configuration file found.")
			fmt.Println("Specify a file with: tinymonitor validate -c /path/to/config.toml")
			os.Exit(1)
		}
	}

	fmt.Printf("Validating: %s\n", configPath)

	_, err := config.LoadAndValidate(configPath)
	if err != nil {
		fmt.Println()
		fmt.Println("Configuration errors:")
		if validationErrs, ok := err.(config.ValidationErrors); ok {
			for _, e := range validationErrs {
				fmt.Printf("  - %s: %s\n", e.Field, e.Message)
			}
		} else {
			fmt.Printf("  - %s\n", err)
		}
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println("Configuration is valid.")
}
