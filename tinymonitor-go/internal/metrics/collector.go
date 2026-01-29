package metrics

import "github.com/Gu1llaum-3/tinymonitor/internal/models"

// Collector defines the interface for metric collectors
type Collector interface {
	// Name returns the name of the collector
	Name() string

	// Check executes the check and returns a list of results
	// Each result contains: component name, severity level (nil if OK), and formatted value
	Check() []models.MetricResult

	// Duration returns the configured duration threshold in seconds
	Duration() int
}
