package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Refresh != 2 {
		t.Errorf("Expected Refresh=2, got %d", cfg.Refresh)
	}

	if cfg.Cooldown != 60 {
		t.Errorf("Expected Cooldown=60, got %d", cfg.Cooldown)
	}

	if cfg.CPU.Warning != 70 {
		t.Errorf("Expected CPU.Warning=70, got %f", cfg.CPU.Warning)
	}

	if cfg.CPU.Critical != 90 {
		t.Errorf("Expected CPU.Critical=90, got %f", cfg.CPU.Critical)
	}

	if !cfg.CPU.Enabled {
		t.Error("Expected CPU.Enabled=true")
	}

	// Load should use auto mode by default
	if !cfg.Load.Auto {
		t.Error("Expected Load.Auto=true")
	}

	if cfg.Load.WarningRatio != 0.7 {
		t.Errorf("Expected Load.WarningRatio=0.7, got %f", cfg.Load.WarningRatio)
	}

	if cfg.Load.CriticalRatio != 0.9 {
		t.Errorf("Expected Load.CriticalRatio=0.9, got %f", cfg.Load.CriticalRatio)
	}

	// Verify GetThresholds calculates correctly
	cpuCount := runtime.NumCPU()
	expectedWarning := float64(cpuCount) * 0.7
	expectedCritical := float64(cpuCount) * 0.9
	warning, critical := cfg.Load.GetThresholds()

	if warning != expectedWarning {
		t.Errorf("Expected Load warning threshold=%f, got %f", expectedWarning, warning)
	}

	if critical != expectedCritical {
		t.Errorf("Expected Load critical threshold=%f, got %f", expectedCritical, critical)
	}
}

func TestLoadFromTOMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
refresh = 10
cooldown = 180

[cpu]
enabled = true
warning = 75
critical = 95
duration = 15

[memory]
enabled = false

[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/test"

  [alerts.ntfy.rules]
  default = ["WARNING", "CRITICAL"]
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test TOML config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Refresh != 10 {
		t.Errorf("Expected Refresh=10, got %d", cfg.Refresh)
	}

	if cfg.Cooldown != 180 {
		t.Errorf("Expected Cooldown=180, got %d", cfg.Cooldown)
	}

	if cfg.CPU.Warning != 75 {
		t.Errorf("Expected CPU.Warning=75, got %f", cfg.CPU.Warning)
	}

	if cfg.CPU.Duration != 15 {
		t.Errorf("Expected CPU.Duration=15, got %d", cfg.CPU.Duration)
	}

	if cfg.Memory.Enabled {
		t.Error("Expected Memory.Enabled=false")
	}

	if !cfg.Alerts.Ntfy.Enabled {
		t.Error("Expected Ntfy.Enabled=true")
	}

	if cfg.Alerts.Ntfy.TopicURL != "https://ntfy.sh/test" {
		t.Errorf("Expected Ntfy.TopicURL='https://ntfy.sh/test', got '%s'", cfg.Alerts.Ntfy.TopicURL)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load should not return error for empty path: %v", err)
	}

	if cfg == nil {
		t.Error("Expected non-nil config")
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expectError bool
		errorField  string
	}{
		{
			name: "valid config",
			config: `
refresh = 5
cooldown = 60

[cpu]
enabled = true
warning = 70
critical = 90
`,
			expectError: false,
		},
		{
			name: "invalid refresh",
			config: `
refresh = 0
cooldown = 60
`,
			expectError: true,
			errorField:  "refresh",
		},
		{
			name: "warning >= critical",
			config: `
refresh = 5
cooldown = 60

[cpu]
enabled = true
warning = 90
critical = 70
`,
			expectError: true,
			errorField:  "cpu",
		},
		{
			name: "threshold out of range",
			config: `
refresh = 5
cooldown = 60

[cpu]
enabled = true
warning = 150
critical = 200
`,
			expectError: true,
			errorField:  "cpu.warning",
		},
		{
			name: "ntfy enabled without topic_url",
			config: `
refresh = 5
cooldown = 60

[alerts.ntfy]
enabled = true
topic_url = ""
`,
			expectError: true,
			errorField:  "alerts.ntfy.topic_url",
		},
		{
			name: "smtp enabled without host",
			config: `
refresh = 5
cooldown = 60

[alerts.smtp]
enabled = true
host = ""
port = 587
`,
			expectError: true,
			errorField:  "alerts.smtp.host",
		},
		{
			name: "smtp invalid port",
			config: `
refresh = 5
cooldown = 60

[alerts.smtp]
enabled = true
host = "smtp.example.com"
port = 99999
user = "user"
password = "pass"
from_addr = "test@example.com"
to_addrs = ["admin@example.com"]
`,
			expectError: true,
			errorField:  "alerts.smtp.port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")

			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			_, err := LoadAndValidate(configPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error, got nil")
					return
				}
				if validationErrs, ok := err.(ValidationErrors); ok {
					found := false
					for _, e := range validationErrs {
						if e.Field == tt.errorField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error for field '%s', got: %v", tt.errorField, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
