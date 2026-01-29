package metrics

import (
	"fmt"
	"strings"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/disk"
)

// DiskCollector monitors filesystem usage
type DiskCollector struct {
	name   string
	config config.FilesystemConfig
}

// NewDiskCollector creates a new disk/filesystem collector
func NewDiskCollector(cfg config.FilesystemConfig) *DiskCollector {
	return &DiskCollector{
		name:   "filesystem",
		config: cfg,
	}
}

// Name returns the collector name
func (c *DiskCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *DiskCollector) Duration() int {
	return c.config.Duration
}

// Check executes the filesystem check
func (c *DiskCollector) Check() []models.MetricResult {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil
	}

	var results []models.MetricResult

	for _, part := range partitions {
		// Filter out snap loops and squashfs
		if strings.Contains(part.Device, "loop") || part.Fstype == "squashfs" {
			continue
		}

		// Filter out Docker overlay
		if strings.Contains(part.Mountpoint, "docker") || part.Fstype == "overlay" {
			continue
		}

		// Skip cdrom or empty fstype
		if containsOpt(part.Opts, "cdrom") || part.Fstype == "" {
			continue
		}

		// User defined excludes
		excluded := false
		for _, ex := range c.config.Exclude {
			if strings.Contains(part.Mountpoint, ex) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		usage, err := disk.Usage(part.Mountpoint)
		if err != nil {
			continue
		}

		usagePercent := usage.UsedPercent
		var level *models.Severity

		if usagePercent >= c.config.Critical {
			sev := models.SeverityCritical
			level = &sev
		} else if usagePercent >= c.config.Warning {
			sev := models.SeverityWarning
			level = &sev
		}

		componentName := fmt.Sprintf("DISK:%s", part.Mountpoint)
		results = append(results, models.NewMetricResult(componentName, level, fmt.Sprintf("%.1f%%", usagePercent)))
	}

	return results
}

func containsOpt(opts []string, target string) bool {
	for _, opt := range opts {
		if strings.Contains(opt, target) {
			return true
		}
	}
	return false
}
