package metrics

import (
	"os"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

// RebootCollector checks if system requires reboot
type RebootCollector struct {
	name   string
	config config.RebootConfig
}

// NewRebootCollector creates a new reboot collector
func NewRebootCollector(cfg config.RebootConfig) *RebootCollector {
	return &RebootCollector{
		name:   "reboot",
		config: cfg,
	}
}

// Name returns the collector name
func (c *RebootCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *RebootCollector) Duration() int {
	return c.config.Duration
}

// Check executes the reboot check
func (c *RebootCollector) Check() []models.MetricResult {
	rebootRequired := false
	details := "OK"

	// Debian / Ubuntu / Mint standard
	// Check for the flag file created by apt/dpkg
	if fileExists("/var/run/reboot-required") || fileExists("/run/reboot-required") {
		rebootRequired = true
		details = "System requires a reboot (updates installed)"
	}

	var level *models.Severity
	if rebootRequired {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult("REBOOT", level, details),
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
