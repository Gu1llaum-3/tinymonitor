package cmd

import (
	"fmt"
	"os"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display configuration summary",
	Long: `Display a human-readable summary of a TinyMonitor configuration file.

Shows:
  - Global settings (refresh, cooldown, log file)
  - Enabled/disabled metrics with their thresholds
  - Configured alert providers

Examples:
  tinymonitor info
  tinymonitor info -c /etc/tinymonitor/config.toml`,
	Run: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) {
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
	}

	var cfg *config.Config
	var err error

	if configPath == "" {
		fmt.Println("Configuration: (using defaults)")
		cfg = config.Default()
	} else {
		cfg, err = config.LoadAndValidate(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration: %s\n", configPath)
	}

	fmt.Println()
	printGlobalSettings(cfg)
	fmt.Println()
	printMetrics(cfg)
	fmt.Println()
	printAlertProviders(cfg)
}

func printGlobalSettings(cfg *config.Config) {
	fmt.Println("Global Settings")

	fmt.Printf("  Refresh:   %ds\n", cfg.Refresh)

	if cfg.Cooldown == -1 {
		fmt.Println("  Cooldown:  once per incident")
	} else {
		fmt.Printf("  Cooldown:  %ds\n", cfg.Cooldown)
	}

	if cfg.LogFile == "" {
		fmt.Println("  Log File:  (stdout)")
	} else {
		fmt.Printf("  Log File:  %s\n", cfg.LogFile)
	}
}

func printMetrics(cfg *config.Config) {
	fmt.Println("Metrics")

	// CPU
	if cfg.CPU.Enabled {
		dur := formatDuration(cfg.CPU.Duration)
		fmt.Printf("  [✓] CPU         warning: %.0f%%    critical: %.0f%%%s\n",
			cfg.CPU.Warning, cfg.CPU.Critical, dur)
	} else {
		fmt.Println("  [✗] CPU         (disabled)")
	}

	// Memory
	if cfg.Memory.Enabled {
		dur := formatDuration(cfg.Memory.Duration)
		fmt.Printf("  [✓] Memory      warning: %.0f%%    critical: %.0f%%%s\n",
			cfg.Memory.Warning, cfg.Memory.Critical, dur)
	} else {
		fmt.Println("  [✗] Memory      (disabled)")
	}

	// Filesystem
	if cfg.Filesystem.Enabled {
		dur := formatDuration(cfg.Filesystem.Duration)
		excludeInfo := ""
		if len(cfg.Filesystem.Exclude) > 0 {
			excludeInfo = fmt.Sprintf("    exclude: %d paths", len(cfg.Filesystem.Exclude))
		}
		fmt.Printf("  [✓] Filesystem  warning: %.0f%%    critical: %.0f%%%s%s\n",
			cfg.Filesystem.Warning, cfg.Filesystem.Critical, dur, excludeInfo)
	} else {
		fmt.Println("  [✗] Filesystem  (disabled)")
	}

	// Load
	if cfg.Load.Enabled {
		dur := formatDuration(cfg.Load.Duration)
		fmt.Printf("  [✓] Load        warning: %.1f     critical: %.1f%s\n",
			cfg.Load.Warning, cfg.Load.Critical, dur)
	} else {
		fmt.Println("  [✗] Load        (disabled)")
	}

	// I/O
	if cfg.IO.Enabled {
		dur := formatDuration(cfg.IO.Duration)
		warning := formatIOValue(cfg.IO.Warning)
		critical := formatIOValue(cfg.IO.Critical)
		fmt.Printf("  [✓] I/O         warning: %s   critical: %s%s\n",
			warning, critical, dur)
	} else {
		fmt.Println("  [✗] I/O         (disabled)")
	}

	// Reboot
	if cfg.Reboot.Enabled {
		fmt.Println("  [✓] Reboot      (checks /var/run/reboot-required)")
	} else {
		fmt.Println("  [✗] Reboot      (disabled)")
	}
}

func printAlertProviders(cfg *config.Config) {
	fmt.Println("Alert Providers")

	// Ntfy
	if cfg.Alerts.Ntfy.Enabled {
		fmt.Printf("  [✓] Ntfy        %s\n", cfg.Alerts.Ntfy.TopicURL)
	} else {
		fmt.Println("  [✗] Ntfy")
	}

	// Google Chat
	if cfg.Alerts.GoogleChat.Enabled {
		fmt.Printf("  [✓] Google Chat %s\n", truncateURL(cfg.Alerts.GoogleChat.WebhookURL))
	} else {
		fmt.Println("  [✗] Google Chat")
	}

	// SMTP
	if cfg.Alerts.SMTP.Enabled {
		fmt.Printf("  [✓] SMTP        %s:%d → %d recipient(s)\n",
			cfg.Alerts.SMTP.Host, cfg.Alerts.SMTP.Port, len(cfg.Alerts.SMTP.ToAddrs))
	} else {
		fmt.Println("  [✗] SMTP")
	}

	// Webhook
	if cfg.Alerts.Webhook.Enabled {
		fmt.Printf("  [✓] Webhook     %s\n", truncateURL(cfg.Alerts.Webhook.URL))
	} else {
		fmt.Println("  [✗] Webhook")
	}

	// Gotify
	if cfg.Alerts.Gotify.Enabled {
		fmt.Printf("  [✓] Gotify      %s\n", cfg.Alerts.Gotify.URL)
	} else {
		fmt.Println("  [✗] Gotify")
	}
}

func formatDuration(d int) string {
	if d > 0 {
		return fmt.Sprintf("    duration: %ds", d)
	}
	return ""
}

func formatIOValue(v interface{}) string {
	if v == nil {
		return "N/A"
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func truncateURL(url string) string {
	if len(url) > 50 {
		return url[:47] + "..."
	}
	return url
}
