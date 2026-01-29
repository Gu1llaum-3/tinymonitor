package metrics

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/shirou/gopsutil/v3/disk"
)

// IOCollector monitors disk I/O
type IOCollector struct {
	name         string
	config       config.IOConfig
	lastCounters *disk.IOCountersStat
	lastTime     time.Time
	mu           sync.Mutex
}

// NewIOCollector creates a new I/O collector
func NewIOCollector(cfg config.IOConfig) *IOCollector {
	counters, _ := disk.IOCounters()
	var total *disk.IOCountersStat
	if len(counters) > 0 {
		t := aggregateCounters(counters)
		total = &t
	}

	return &IOCollector{
		name:         "io",
		config:       cfg,
		lastCounters: total,
		lastTime:     time.Now(),
	}
}

func aggregateCounters(counters map[string]disk.IOCountersStat) disk.IOCountersStat {
	var total disk.IOCountersStat
	for _, c := range counters {
		total.ReadBytes += c.ReadBytes
		total.WriteBytes += c.WriteBytes
	}
	return total
}

// Name returns the collector name
func (c *IOCollector) Name() string {
	return c.name
}

// Duration returns the configured duration threshold
func (c *IOCollector) Duration() int {
	return c.config.Duration
}

// parseThreshold parses a threshold value (number, string with unit, or percentage)
func (c *IOCollector) parseThreshold(value interface{}, maxValue *float64) float64 {
	if value == nil {
		return math.Inf(1)
	}

	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		return c.parseStringThreshold(v, maxValue)
	}

	return math.Inf(1)
}

func (c *IOCollector) parseStringThreshold(value string, maxValue *float64) float64 {
	value = strings.TrimSpace(strings.ToUpper(value))

	// Handle percentage
	if strings.HasSuffix(value, "%") {
		if maxValue == nil {
			return math.Inf(1)
		}
		percent, err := strconv.ParseFloat(value[:len(value)-1], 64)
		if err != nil {
			return math.Inf(1)
		}
		return *maxValue * (percent / 100.0)
	}

	// Handle units
	units := map[string]float64{
		"K": 1024,
		"M": 1024 * 1024,
		"G": 1024 * 1024 * 1024,
		"T": 1024 * 1024 * 1024 * 1024,
	}

	for unit, multiplier := range units {
		if strings.HasSuffix(value, unit+"B") {
			numPart := value[:len(value)-2]
			num, err := strconv.ParseFloat(numPart, 64)
			if err == nil {
				return num * multiplier
			}
		} else if strings.HasSuffix(value, unit) {
			numPart := value[:len(value)-1]
			num, err := strconv.ParseFloat(numPart, 64)
			if err == nil {
				return num * multiplier
			}
		}
	}

	// Try plain number
	num, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return num
	}

	return math.Inf(1)
}

func formatBytes(size float64) string {
	power := 1024.0
	n := 0
	labels := []string{"", "K", "M", "G", "T"}

	for size > power && n < len(labels)-1 {
		size /= power
		n++
	}

	return fmt.Sprintf("%.1f%sB/s", size, labels[n])
}

// Check executes the I/O check
func (c *IOCollector) Check() []models.MetricResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	currentTime := time.Now()
	counters, err := disk.IOCounters()
	if err != nil || len(counters) == 0 {
		return nil
	}

	currentCounters := aggregateCounters(counters)

	if c.lastCounters == nil {
		c.lastCounters = &currentCounters
		c.lastTime = currentTime
		return nil
	}

	timeDelta := currentTime.Sub(c.lastTime).Seconds()
	if timeDelta <= 0 {
		return nil
	}

	readBytesDelta := float64(currentCounters.ReadBytes - c.lastCounters.ReadBytes)
	writeBytesDelta := float64(currentCounters.WriteBytes - c.lastCounters.WriteBytes)

	if readBytesDelta < 0 {
		readBytesDelta = 0
	}
	if writeBytesDelta < 0 {
		writeBytesDelta = 0
	}

	readSpeed := readBytesDelta / timeDelta
	writeSpeed := writeBytesDelta / timeDelta

	c.lastCounters = &currentCounters
	c.lastTime = currentTime

	formattedRead := formatBytes(readSpeed)
	formattedWrite := formatBytes(writeSpeed)

	var level *models.Severity

	// Parse max_speed if present
	var maxSpeed *float64
	if c.config.MaxSpeed != nil {
		ms := c.parseThreshold(c.config.MaxSpeed, nil)
		if !math.IsInf(ms, 1) {
			maxSpeed = &ms
		}
	}

	warningThreshold := c.parseThreshold(c.config.Warning, maxSpeed)
	criticalThreshold := c.parseThreshold(c.config.Critical, maxSpeed)

	totalSpeed := readSpeed + writeSpeed

	if totalSpeed >= criticalThreshold {
		sev := models.SeverityCritical
		level = &sev
	} else if totalSpeed >= warningThreshold {
		sev := models.SeverityWarning
		level = &sev
	}

	valueStr := fmt.Sprintf("R: %s W: %s", formattedRead, formattedWrite)
	return []models.MetricResult{
		models.NewMetricResult("I/O", level, valueStr),
	}
}
