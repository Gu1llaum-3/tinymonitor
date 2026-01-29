package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the main configuration
type Config struct {
	Refresh    int              `json:"refresh"`
	Cooldown   int              `json:"cooldown"`
	LogFile    string           `json:"log_file"`
	Load       MetricConfig     `json:"load"`
	CPU        MetricConfig     `json:"cpu"`
	Memory     MetricConfig     `json:"memory"`
	Filesystem FilesystemConfig `json:"filesystem"`
	Reboot     RebootConfig     `json:"reboot"`
	IO         IOConfig         `json:"io"`
	Alerts     AlertsConfig     `json:"alerts"`
}

// MetricConfig represents configuration for a simple metric
type MetricConfig struct {
	Warning  float64 `json:"warning"`
	Critical float64 `json:"critical"`
	Enabled  bool    `json:"enabled"`
	Duration int     `json:"duration"`
}

// FilesystemConfig represents filesystem metric configuration
type FilesystemConfig struct {
	Warning  float64  `json:"warning"`
	Critical float64  `json:"critical"`
	Enabled  bool     `json:"enabled"`
	Duration int      `json:"duration"`
	Exclude  []string `json:"exclude"`
}

// RebootConfig represents reboot metric configuration
type RebootConfig struct {
	Enabled  bool `json:"enabled"`
	Duration int  `json:"duration"`
}

// IOConfig represents I/O metric configuration
type IOConfig struct {
	Warning  interface{} `json:"warning"`  // Can be number, string with unit, or percentage
	Critical interface{} `json:"critical"` // Can be number, string with unit, or percentage
	MaxSpeed interface{} `json:"max_speed"`
	Enabled  bool        `json:"enabled"`
	Duration int         `json:"duration"`
}

// AlertsConfig represents all alert providers configuration
type AlertsConfig struct {
	GoogleChat GoogleChatConfig `json:"google_chat"`
	Ntfy       NtfyConfig       `json:"ntfy"`
	SMTP       SMTPConfig       `json:"smtp"`
	Webhook    WebhookConfig    `json:"webhook"`
	Gotify     GotifyConfig     `json:"gotify"`
}

// ProviderRules represents alert filtering rules
type ProviderRules map[string][]string

// GoogleChatConfig represents Google Chat alert configuration
type GoogleChatConfig struct {
	Enabled    bool          `json:"enabled"`
	WebhookURL string        `json:"webhook_url"`
	Levels     []string      `json:"levels"`
	Rules      ProviderRules `json:"rules"`
}

// NtfyConfig represents Ntfy alert configuration
type NtfyConfig struct {
	Enabled  bool          `json:"enabled"`
	TopicURL string        `json:"topic_url"`
	Token    string        `json:"token"`
	Levels   []string      `json:"levels"`
	Rules    ProviderRules `json:"rules"`
}

// SMTPConfig represents SMTP alert configuration
type SMTPConfig struct {
	Enabled  bool          `json:"enabled"`
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	User     string        `json:"user"`
	Password string        `json:"password"`
	FromAddr string        `json:"from_addr"`
	ToAddrs  interface{}   `json:"to_addrs"` // Can be string or []string
	UseTLS   bool          `json:"use_tls"`
	Levels   []string      `json:"levels"`
	Rules    ProviderRules `json:"rules"`
}

// WebhookConfig represents generic webhook alert configuration
type WebhookConfig struct {
	Enabled bool              `json:"enabled"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Timeout int               `json:"timeout"`
	Levels  []string          `json:"levels"`
	Rules   ProviderRules     `json:"rules"`
}

// GotifyConfig represents Gotify alert configuration
type GotifyConfig struct {
	Enabled bool          `json:"enabled"`
	URL     string        `json:"url"`
	Token   string        `json:"token"`
	Levels  []string      `json:"levels"`
	Rules   ProviderRules `json:"rules"`
}

// Default returns a new Config with default values
func Default() *Config {
	cpuCount := runtime.NumCPU()
	loadWarning := float64(cpuCount) * 0.7
	loadCritical := float64(cpuCount) * 0.9

	return &Config{
		Refresh:  2,
		Cooldown: 60,
		LogFile:  "tinymonitor.log",
		Load: MetricConfig{
			Warning:  loadWarning,
			Critical: loadCritical,
			Enabled:  true,
			Duration: 60,
		},
		CPU: MetricConfig{
			Warning:  70,
			Critical: 90,
			Enabled:  true,
			Duration: 0,
		},
		Memory: MetricConfig{
			Warning:  70,
			Critical: 90,
			Enabled:  true,
			Duration: 0,
		},
		Filesystem: FilesystemConfig{
			Warning:  80,
			Critical: 90,
			Enabled:  true,
			Duration: 0,
			Exclude:  []string{},
		},
		Reboot: RebootConfig{
			Enabled:  true,
			Duration: 0,
		},
		IO: IOConfig{
			Enabled:  true,
			Duration: 0,
		},
		Alerts: AlertsConfig{
			GoogleChat: GoogleChatConfig{
				Enabled: false,
			},
			Ntfy: NtfyConfig{
				Enabled: false,
			},
			SMTP: SMTPConfig{
				Enabled: false,
				Port:    587,
				UseTLS:  true,
			},
			Webhook: WebhookConfig{
				Enabled: false,
				Timeout: 10,
			},
			Gotify: GotifyConfig{
				Enabled: false,
			},
		},
	}
}

// Load loads configuration from file with cascade search
func Load(configPath string) (*Config, error) {
	config := Default()

	// 1. Explicit path from CLI
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Loading config from: %s\n", configPath)
			if err := loadFromFile(configPath, config); err != nil {
				fmt.Printf("Error loading config file %s: %v\n", configPath, err)
			}
		} else {
			fmt.Printf("Warning: Config file %s not found. Using defaults.\n", configPath)
		}
		return config, nil
	}

	// 2. Search in standard locations (Priority order)
	searchPaths := []string{
		filepath.Join(getCurrentDir(), "config.json"),
		filepath.Join(getHomeDir(), ".config", "tinymonitor", "config.json"),
		"/etc/tinymonitor/config.json",
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Loading config from: %s\n", path)
			if err := loadFromFile(path, config); err != nil {
				fmt.Printf("Error loading config file %s: %v\n", path, err)
				return config, nil
			}
			return config, nil
		}
	}

	fmt.Println("No config file found. Using internal defaults.")
	return config, nil
}

func loadFromFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}
