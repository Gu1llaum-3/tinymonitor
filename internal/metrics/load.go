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
	config config.MetricConfig
}

// NewLoadCollector creates a new load average collector
func NewLoadCollector(cfg config.MetricConfig) *LoadCollector {
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
	var level *models.Severity

	if load1 >= c.config.Critical {
		sev := models.SeverityCritical
		level = &sev
	} else if load1 >= c.config.Warning {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult("LOAD", level, fmt.Sprintf("%.2f", load1)),
	}
}
