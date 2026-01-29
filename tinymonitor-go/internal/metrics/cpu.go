package metrics

import (
	"fmt"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUCollector monitors CPU usage
type CPUCollector struct {
	name   string
	config config.MetricConfig
}

// NewCPUCollector creates a new CPU collector
func NewCPUCollector(cfg config.MetricConfig) *CPUCollector {
	// Initialize CPU percent calculation (first call returns 0)
	cpu.Percent(0, false)

	return &CPUCollector{
		name:   "cpu",
		config: cfg,
	}
}

// Name returns the collector name
func (c *CPUCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *CPUCollector) Duration() int {
	return c.config.Duration
}

// Check executes the CPU check
func (c *CPUCollector) Check() []models.MetricResult {
	percents, err := cpu.Percent(0, false)
	if err != nil || len(percents) == 0 {
		return nil
	}

	cpuPercent := percents[0]
	var level *models.Severity

	if cpuPercent >= c.config.Critical {
		sev := models.SeverityCritical
		level = &sev
	} else if cpuPercent >= c.config.Warning {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult("CPU", level, fmt.Sprintf("%.1f%%", cpuPercent)),
	}
}
