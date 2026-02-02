package metrics

import (
	"fmt"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/load"
)

// LoadCollector monitors system load average
type LoadCollector struct {
	name   string
	config config.LoadConfig
}

// NewLoadCollector creates a new load average collector
func NewLoadCollector(cfg config.LoadConfig) *LoadCollector {
	return &LoadCollector{
		name:   "load",
		config: cfg,
	}
}

// Name returns the collector name
func (c *LoadCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *LoadCollector) Duration() int {
	return c.config.Duration
}

// Check executes the load average check
func (c *LoadCollector) Check() []models.MetricResult {
	avg, err := load.Avg()
	if err != nil {
		// Windows doesn't support load average
		return nil
	}

	load1 := avg.Load1
	warning, critical := c.config.GetThresholds()
	var level *models.Severity

	if load1 >= critical {
		sev := models.SeverityCritical
		level = &sev
	} else if load1 >= warning {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult("LOAD", level, fmt.Sprintf("%.2f", load1)),
	}
}
