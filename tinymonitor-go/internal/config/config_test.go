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

	cpuCount := runtime.NumCPU()
	expectedLoadWarning := float64(cpuCount) * 0.7
	if cfg.Load.Warning != expectedLoadWarning {
		t.Errorf("Expected Load.Warning=%f, got %f", expectedLoadWarning, cfg.Load.Warning)
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"refresh": 5,
		"cooldown": 120,
		"cpu": {
			"warning": 80,
			"critical": 95,
			"enabled": true,
			"duration": 30
		},
		"alerts": {
			"google_chat": {
				"enabled": true,
				"webhook_url": "https://example.com/webhook"
			}
		}
	}`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Refresh != 5 {
		t.Errorf("Expected Refresh=5, got %d", cfg.Refresh)
	}

	if cfg.Cooldown != 120 {
		t.Errorf("Expected Cooldown=120, got %d", cfg.Cooldown)
	}

	if cfg.CPU.Warning != 80 {
		t.Errorf("Expected CPU.Warning=80, got %f", cfg.CPU.Warning)
	}

	if cfg.CPU.Duration != 30 {
		t.Errorf("Expected CPU.Duration=30, got %d", cfg.CPU.Duration)
	}

	if !cfg.Alerts.GoogleChat.Enabled {
		t.Error("Expected GoogleChat.Enabled=true")
	}

	if cfg.Alerts.GoogleChat.WebhookURL != "https://example.com/webhook" {
		t.Errorf("Expected GoogleChat.WebhookURL='https://example.com/webhook', got '%s'", cfg.Alerts.GoogleChat.WebhookURL)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("Load should not return error for missing file: %v", err)
	}

	// Should return defaults
	if cfg.Refresh != 2 {
		t.Errorf("Expected default Refresh=2, got %d", cfg.Refresh)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load should not return error for empty path: %v", err)
	}

	// Should return defaults (or find config in standard locations)
	if cfg == nil {
		t.Error("Expected non-nil config")
	}
}
