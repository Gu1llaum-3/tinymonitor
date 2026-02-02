package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/alerts"
	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
	"github.com/spf13/cobra"
)

var testAlertCmd = &cobra.Command{
	Use:   "test-alert",
	Short: "Send a test alert to verify configuration",
	Long: `Send a test alert to all enabled alert providers to verify your configuration.

This is useful for:
  - Verifying credentials are correct
  - Checking that notifications arrive properly
  - Testing without waiting for a real issue

Examples:
  tinymonitor test-alert                     # Test all enabled providers
  tinymonitor test-alert --provider ntfy     # Test only Ntfy
  tinymonitor test-alert --provider smtp     # Test only SMTP
  tinymonitor test-alert -c config.toml      # Use specific config`,
	Run: runTestAlert,
}

var testProvider string

func init() {
	rootCmd.AddCommand(testAlertCmd)
	testAlertCmd.Flags().StringVarP(&testProvider, "provider", "p", "", "Test only a specific provider (ntfy, smtp, google_chat, webhook, gotify)")
}

func runTestAlert(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Build list of providers to test
	providers := buildProviderList(cfg.Alerts, testProvider)

	if len(providers) == 0 {
		if testProvider != "" {
			fmt.Fprintf(os.Stderr, "Provider '%s' is not enabled or does not exist.\n", testProvider)
			fmt.Fprintln(os.Stderr, "Available providers: ntfy, smtp, google_chat, webhook, gotify")
		} else {
			fmt.Fprintln(os.Stderr, "No alert providers are enabled in your configuration.")
			fmt.Fprintln(os.Stderr, "Enable at least one provider in your config file.")
		}
		os.Exit(1)
	}

	// Create test alert
	alert := createTestAlert()

	fmt.Println("TinyMonitor Test Alert")
	fmt.Println("======================")
	fmt.Println()
	fmt.Printf("Sending test alert to %d provider(s)...\n\n", len(providers))

	// Send to each provider and report results
	var success, failed int
	for _, provider := range providers {
		fmt.Printf("  %-12s ", provider.Name()+":")

		err := provider.Send(alert)
		if err != nil {
			fmt.Printf("FAILED - %v\n", err)
			failed++
		} else {
			fmt.Println("OK")
			success++
		}
	}

	// Summary
	fmt.Println()
	if failed == 0 {
		fmt.Printf("All %d provider(s) successfully sent the test alert.\n", success)
	} else {
		fmt.Printf("Results: %d succeeded, %d failed\n", success, failed)
		os.Exit(1)
	}
}

func buildProviderList(cfg config.AlertsConfig, filterProvider string) []alerts.Provider {
	var providers []alerts.Provider
	filter := strings.ToLower(filterProvider)

	// Ntfy
	if cfg.Ntfy.Enabled && (filter == "" || filter == "ntfy") {
		providers = append(providers, alerts.NewNtfyProvider(cfg.Ntfy))
	}

	// Google Chat
	if cfg.GoogleChat.Enabled && (filter == "" || filter == "google_chat" || filter == "googlechat") {
		providers = append(providers, alerts.NewGoogleChatProvider(cfg.GoogleChat))
	}

	// SMTP
	if cfg.SMTP.Enabled && (filter == "" || filter == "smtp" || filter == "email") {
		providers = append(providers, alerts.NewSMTPProvider(cfg.SMTP))
	}

	// Webhook
	if cfg.Webhook.Enabled && (filter == "" || filter == "webhook") {
		providers = append(providers, alerts.NewWebhookProvider(cfg.Webhook))
	}

	// Gotify
	if cfg.Gotify.Enabled && (filter == "" || filter == "gotify") {
		providers = append(providers, alerts.NewGotifyProvider(cfg.Gotify))
	}

	return providers
}

func createTestAlert() models.Alert {
	hostname := utils.GetHostname()

	return models.Alert{
		Component: "TEST",
		Level:     models.SeverityWarning,
		Value:     "This is a test alert",
		Title:     fmt.Sprintf("Test Alert from %s", hostname),
		Message:   fmt.Sprintf("This is a test alert from TinyMonitor on %s. If you receive this, your alert configuration is working correctly.", hostname),
		Timestamp: time.Now(),
	}
}
