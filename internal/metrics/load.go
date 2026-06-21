package metrics

import (
	"fmt"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/load"
)

// LoadCollector monitors a single system load-average window (5 or 15 minutes).
// Window selection, thresholds and duration are resolved once at construction
// so Check() stays branch-free on the hot path.
type LoadCollector struct {
	name      string
	component string
	useLoad15 bool
	duration  int
	warning   float64
	critical  float64
}

// NewLoadCollector creates a load average collector for the given window
// (5 or 15 minutes). The 1-minute average is intentionally not supported as it
// is too noisy for alerting. Any window other than 15 is treated as the
// 5-minute window.
func NewLoadCollector(window int, cfg config.LoadConfig) *LoadCollector {
	wc := cfg.Window5
	component := "LOAD5"
	useLoad15 := false
	if window == 15 {
		wc = cfg.Window15
		component = "LOAD15"
		useLoad15 = true
	}

	warning, critical := cfg.ThresholdsFor(wc)
	return &LoadCollector{
		name:      fmt.Sprintf("load%d", window),
		component: component,
		useLoad15: useLoad15,
		duration:  wc.Duration,
		warning:   warning,
		critical:  critical,
	}
}

// Name returns the collector name
func (c *LoadCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold for this window
func (c *LoadCollector) Duration() int {
	return c.duration
}

// Check executes the load average check for this window
func (c *LoadCollector) Check() []models.MetricResult {
	avg, err := load.Avg()
	if err != nil {
		// Windows doesn't support load average
		return nil
	}

	value := avg.Load5
	if c.useLoad15 {
		value = avg.Load15
	}

	var level *models.Severity
	if value >= c.critical {
		sev := models.SeverityCritical
		level = &sev
	} else if value >= c.warning {
		sev := models.SeverityWarning
		level = &sev
	}

	return []models.MetricResult{
		models.NewMetricResult(c.component, level, fmt.Sprintf("%.2f", value)),
	}
}
