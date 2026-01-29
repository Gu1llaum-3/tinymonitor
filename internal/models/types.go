package models

import "time"

// Severity represents the alert severity level
type Severity string

const (
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// MetricResult represents the result of a metric check
type MetricResult struct {
	Component string
	Level     *Severity // nil means OK/normal
	Value     string
}

// Alert represents an alert to be sent
type Alert struct {
	Component string
	Level     Severity
	Value     string
	Title     string
	Message   string
	Timestamp time.Time
}

// AlertState tracks the state of an alert for duration-based alerting
type AlertState struct {
	Level          Severity
	StartTime      time.Time
	AlertTriggered bool
}

// NewMetricResult creates a new MetricResult
func NewMetricResult(component string, level *Severity, value string) MetricResult {
	return MetricResult{
		Component: component,
		Level:     level,
		Value:     value,
	}
}

// NewAlert creates a new Alert
func NewAlert(component string, level Severity, value string) Alert {
	title := "ALERT " + string(level) + " : " + component
	message := "Component " + component + " is in state " + string(level) + ". Value: " + value

	return Alert{
		Component: component,
		Level:     level,
		Value:     value,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
	}
}
