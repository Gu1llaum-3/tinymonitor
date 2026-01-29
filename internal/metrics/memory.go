package metrics

import (
	"fmt"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryCollector monitors memory usage
type MemoryCollector struct {
	name   string
	config config.MetricConfig
}

// NewMemoryCollector creates a new memory collector
func NewMemoryCollector(cfg config.MetricConfig) *MemoryCollector {
	return &MemoryCollector{
		name:   "memory",
		config: cfg,
	}
}

// Name returns the collector name
func (c *MemoryCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *MemoryCollector) Duration() int {
	return c.config.Duration
}

// Check executes the memory check
func (c *MemoryCollector) Check() []models.MetricResult {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}

	memPercent := vmem.UsedPercent
	var level *models.Severity

	if memPercent >= c.config.Critical {
		sev := models.SeverityCritical
		level = &sev
	} else if memPercent >= c.config.Warning {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult("MEMORY", level, fmt.Sprintf("%.1f%%", memPercent)),
	}
}
