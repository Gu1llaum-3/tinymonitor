package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config represents the main configuration
type Config struct {
	Refresh    int              `toml:"refresh"`
	Cooldown   int              `toml:"cooldown"`
	LogFile    string           `toml:"log_file"`
	Load       LoadConfig       `toml:"load"`
	CPU        MetricConfig     `toml:"cpu"`
	Memory     MetricConfig     `toml:"memory"`
	Filesystem FilesystemConfig `toml:"filesystem"`
	Reboot     RebootConfig     `toml:"reboot"`
	IO         IOConfig         `toml:"io"`
	Alerts     AlertsConfig     `toml:"alerts"`
}

// MetricConfig represents configuration for a simple metric
type MetricConfig struct {
	Warning  float64 `toml:"warning"`
	Critical float64 `toml:"critical"`
	Enabled  bool    `toml:"enabled"`
	Duration int     `toml:"duration"`
}

// LoadConfig represents configuration for the load average metric
type LoadConfig struct {
	Enabled       bool    `toml:"enabled"`
	Auto          bool    `toml:"auto"`
	WarningRatio  float64 `toml:"warning_ratio"`
	CriticalRatio float64 `toml:"critical_ratio"`
	Warning       float64 `toml:"warning"`
	Critical      float64 `toml:"critical"`
	Duration      int     `toml:"duration"`
}

// GetThresholds returns the effective warning and critical thresholds
// If Auto is true, thresholds are calculated based on CPU count
func (c *LoadConfig) GetThresholds() (warning, critical float64) {
	if c.Auto {
		cpuCount := float64(runtime.NumCPU())
		return cpuCount * c.WarningRatio, cpuCount * c.CriticalRatio
	}
	return c.Warning, c.Critical
}

// FilesystemConfig represents filesystem metric configuration
type FilesystemConfig struct {
	Warning  float64  `toml:"warning"`
	Critical float64  `toml:"critical"`
	Enabled  bool     `toml:"enabled"`
	Duration int      `toml:"duration"`
	Exclude  []string `toml:"exclude"`
}

// RebootConfig represents reboot metric configuration
type RebootConfig struct {
	Enabled  bool `toml:"enabled"`
	Duration int  `toml:"duration"`
}

// IOConfig represents I/O metric configuration
type IOConfig struct {
	Warning  interface{} `toml:"warning"`
	Critical interface{} `toml:"critical"`
	MaxSpeed interface{} `toml:"max_speed"`
	Enabled  bool        `toml:"enabled"`
	Duration int         `toml:"duration"`
}

// AlertsConfig represents all alert providers configuration
type AlertsConfig struct {
	SendRecovery bool             `toml:"send_recovery"`
	GoogleChat   GoogleChatConfig `toml:"google_chat"`
	Ntfy         NtfyConfig       `toml:"ntfy"`
	SMTP         SMTPConfig       `toml:"smtp"`
	Webhook      WebhookConfig    `toml:"webhook"`
	Gotify       GotifyConfig     `toml:"gotify"`
}

// ProviderRules represents alert filtering rules
type ProviderRules map[string][]string

// GoogleChatConfig represents Google Chat alert configuration
type GoogleChatConfig struct {
	Enabled    bool          `toml:"enabled"`
	WebhookURL string        `toml:"webhook_url"`
	Levels     []string      `toml:"levels"`
	Rules      ProviderRules `toml:"rules"`
}

// NtfyConfig represents Ntfy alert configuration
type NtfyConfig struct {
	Enabled  bool          `toml:"enabled"`
	TopicURL string        `toml:"topic_url"`
	Token    string        `toml:"token"`
	Levels   []string      `toml:"levels"`
	Rules    ProviderRules `toml:"rules"`
}

// SMTPConfig represents SMTP alert configuration
type SMTPConfig struct {
	Enabled  bool          `toml:"enabled"`
	Host     string        `toml:"host"`
	Port     int           `toml:"port"`
	User     string        `toml:"user"`
	Password string        `toml:"password"`
	FromAddr string        `toml:"from_addr"`
	ToAddrs  []string      `toml:"to_addrs"`
	UseTLS   bool          `toml:"use_tls"`
	Levels   []string      `toml:"levels"`
	Rules    ProviderRules `toml:"rules"`
}

// WebhookConfig represents generic webhook alert configuration
type WebhookConfig struct {
	Enabled bool              `toml:"enabled"`
	URL     string            `toml:"url"`
	Headers map[string]string `toml:"headers"`
	Timeout int               `toml:"timeout"`
	Levels  []string          `toml:"levels"`
	Rules   ProviderRules     `toml:"rules"`
}

// GotifyConfig represents Gotify alert configuration
type GotifyConfig struct {
	Enabled bool          `toml:"enabled"`
	URL     string        `toml:"url"`
	Token   string        `toml:"token"`
	Levels  []string      `toml:"levels"`
	Rules   ProviderRules `toml:"rules"`
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "\n")
}

// Default returns a new Config with default values
func Default() *Config {
	return &Config{
		Refresh:  2,
		Cooldown: 60,
		LogFile:  "",
		Load: LoadConfig{
			Enabled:       true,
			Auto:          true,
			WarningRatio:  0.7,
			CriticalRatio: 0.9,
			Warning:       0,
			Critical:      0,
			Duration:      180,
		},
		CPU: MetricConfig{
			Warning:  70,
			Critical: 90,
			Enabled:  true,
			Duration: 120,
		},
		Memory: MetricConfig{
			Warning:  70,
			Critical: 90,
			Enabled:  true,
			Duration: 120,
		},
		Filesystem: FilesystemConfig{
			Warning:  80,
			Critical: 90,
			Enabled:  true,
			Duration: 300,
			Exclude:  []string{},
		},
		Reboot: RebootConfig{
			Enabled:  true,
			Duration: 0,
		},
		IO: IOConfig{
			Enabled:  true,
			Duration: 120,
		},
		Alerts: AlertsConfig{
			SendRecovery: true,
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
				return nil, fmt.Errorf("error loading config file %s: %w", configPath, err)
			}
			if errs := config.Validate(); len(errs) > 0 {
				return nil, errs
			}
			return config, nil
		}
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// 2. Search in standard locations (Priority order)
	searchPaths := []string{
		filepath.Join(getCurrentDir(), "config.toml"),
		filepath.Join(getHomeDir(), ".config", "tinymonitor", "config.toml"),
		"/etc/tinymonitor/config.toml",
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Loading config from: %s\n", path)
			if err := loadFromFile(path, config); err != nil {
				return nil, fmt.Errorf("error loading config file %s: %w", path, err)
			}
			if errs := config.Validate(); len(errs) > 0 {
				return nil, errs
			}
			return config, nil
		}
	}

	fmt.Println("No config file found. Using internal defaults.")
	return config, nil
}

// LoadAndValidate loads a config file and returns validation errors without printing
func LoadAndValidate(configPath string) (*Config, error) {
	config := Default()

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	if err := loadFromFile(configPath, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if errs := config.Validate(); len(errs) > 0 {
		return config, errs
	}

	return config, nil
}

func loadFromFile(path string, config *Config) error {
	_, err := toml.DecodeFile(path, config)
	return err
}

// Validate checks the configuration for errors
func (c *Config) Validate() ValidationErrors {
	var errs ValidationErrors

	// Global settings
	if c.Refresh <= 0 {
		errs = append(errs, ValidationError{"refresh", "must be greater than 0"})
	}
	if c.Cooldown < -1 {
		errs = append(errs, ValidationError{"cooldown", "must be >= -1 (-1 = alert once per incident)"})
	}

	// CPU
	if c.CPU.Enabled {
		errs = append(errs, validateThresholds("cpu", c.CPU.Warning, c.CPU.Critical)...)
	}

	// Memory
	if c.Memory.Enabled {
		errs = append(errs, validateThresholds("memory", c.Memory.Warning, c.Memory.Critical)...)
	}

	// Filesystem
	if c.Filesystem.Enabled {
		errs = append(errs, validateThresholds("filesystem", c.Filesystem.Warning, c.Filesystem.Critical)...)
	}

	// Load validation
	if c.Load.Enabled {
		if c.Load.Auto {
			// Validate ratios
			if c.Load.WarningRatio <= 0 {
				errs = append(errs, ValidationError{"load.warning_ratio", "must be greater than 0"})
			}
			if c.Load.CriticalRatio <= 0 {
				errs = append(errs, ValidationError{"load.critical_ratio", "must be greater than 0"})
			}
			if c.Load.WarningRatio >= c.Load.CriticalRatio {
				errs = append(errs, ValidationError{"load", "warning_ratio must be less than critical_ratio"})
			}
		} else {
			// Validate absolute values
			if c.Load.Warning <= 0 {
				errs = append(errs, ValidationError{"load.warning", "must be greater than 0"})
			}
			if c.Load.Critical <= 0 {
				errs = append(errs, ValidationError{"load.critical", "must be greater than 0"})
			}
			if c.Load.Warning >= c.Load.Critical {
				errs = append(errs, ValidationError{"load", "warning must be less than critical"})
			}
		}
	}

	// Alert providers
	if c.Alerts.GoogleChat.Enabled {
		if c.Alerts.GoogleChat.WebhookURL == "" {
			errs = append(errs, ValidationError{"alerts.google_chat.webhook_url", "required when google_chat is enabled"})
		}
	}

	if c.Alerts.Ntfy.Enabled {
		if c.Alerts.Ntfy.TopicURL == "" {
			errs = append(errs, ValidationError{"alerts.ntfy.topic_url", "required when ntfy is enabled"})
		}
	}

	if c.Alerts.SMTP.Enabled {
		if c.Alerts.SMTP.Host == "" {
			errs = append(errs, ValidationError{"alerts.smtp.host", "required when smtp is enabled"})
		}
		if c.Alerts.SMTP.Port < 1 || c.Alerts.SMTP.Port > 65535 {
			errs = append(errs, ValidationError{"alerts.smtp.port", "must be between 1 and 65535"})
		}
		if c.Alerts.SMTP.User == "" {
			errs = append(errs, ValidationError{"alerts.smtp.user", "required when smtp is enabled"})
		}
		if c.Alerts.SMTP.Password == "" {
			errs = append(errs, ValidationError{"alerts.smtp.password", "required when smtp is enabled"})
		}
		if c.Alerts.SMTP.FromAddr == "" {
			errs = append(errs, ValidationError{"alerts.smtp.from_addr", "required when smtp is enabled"})
		}
		if len(c.Alerts.SMTP.ToAddrs) == 0 {
			errs = append(errs, ValidationError{"alerts.smtp.to_addrs", "required when smtp is enabled"})
		}
	}

	if c.Alerts.Webhook.Enabled {
		if c.Alerts.Webhook.URL == "" {
			errs = append(errs, ValidationError{"alerts.webhook.url", "required when webhook is enabled"})
		}
		if c.Alerts.Webhook.Timeout <= 0 {
			errs = append(errs, ValidationError{"alerts.webhook.timeout", "must be greater than 0"})
		}
	}

	if c.Alerts.Gotify.Enabled {
		if c.Alerts.Gotify.URL == "" {
			errs = append(errs, ValidationError{"alerts.gotify.url", "required when gotify is enabled"})
		}
		if c.Alerts.Gotify.Token == "" {
			errs = append(errs, ValidationError{"alerts.gotify.token", "required when gotify is enabled"})
		}
	}

	return errs
}

func validateThresholds(name string, warning, critical float64) ValidationErrors {
	var errs ValidationErrors

	if warning < 0 || warning > 100 {
		errs = append(errs, ValidationError{
			Field:   fmt.Sprintf("%s.warning", name),
			Message: fmt.Sprintf("must be between 0 and 100 (got %.1f)", warning),
		})
	}
	if critical < 0 || critical > 100 {
		errs = append(errs, ValidationError{
			Field:   fmt.Sprintf("%s.critical", name),
			Message: fmt.Sprintf("must be between 0 and 100 (got %.1f)", critical),
		})
	}
	if warning >= critical {
		errs = append(errs, ValidationError{
			Field:   name,
			Message: fmt.Sprintf("warning (%.1f) must be less than critical (%.1f)", warning, critical),
		})
	}

	return errs
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
