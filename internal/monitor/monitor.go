package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/alerts"
	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/metrics"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

// StateChange represents the result of processing a metric state
type StateChange struct {
	ShouldAlert   bool
	IsRecovery    bool
	Level         models.Severity
	PreviousLevel models.Severity
}

// Monitor is the main monitoring loop
type Monitor struct {
	config       *config.Config
	alertManager *alerts.Manager
	collectors   []metrics.Collector
	lastAlert    map[string]time.Time
	alertStates  map[string]*models.AlertState
}

// New creates a new Monitor
func New(cfg *config.Config) *Monitor {
	m := &Monitor{
		config:       cfg,
		alertManager: alerts.NewManager(cfg.Alerts),
		collectors:   make([]metrics.Collector, 0),
		lastAlert:    make(map[string]time.Time),
		alertStates:  make(map[string]*models.AlertState),
	}

	m.loadCollectors()
	return m
}

func (m *Monitor) loadCollectors() {
	if m.config.CPU.Enabled {
		m.collectors = append(m.collectors, metrics.NewCPUCollector(m.config.CPU))
	}

	if m.config.Memory.Enabled {
		m.collectors = append(m.collectors, metrics.NewMemoryCollector(m.config.Memory))
	}

	if m.config.Filesystem.Enabled {
		m.collectors = append(m.collectors, metrics.NewDiskCollector(m.config.Filesystem))
	}

	if m.config.Load.Enabled {
		m.collectors = append(m.collectors, metrics.NewLoadCollector(m.config.Load))
	}

	if m.config.Reboot.Enabled {
		m.collectors = append(m.collectors, metrics.NewRebootCollector(m.config.Reboot))
	}

	if m.config.IO.Enabled {
		m.collectors = append(m.collectors, metrics.NewIOCollector(m.config.IO))
	}
}

// processState manages alert state persistence
// Returns StateChange with information about what action to take
func (m *Monitor) processState(component string, level *models.Severity, value string, duration int) StateChange {
	if level == nil {
		// Return to normal: check if we need to send recovery
		if state, exists := m.alertStates[component]; exists {
			if state.AlertTriggered {
				previousLevel := state.Level
				delete(m.alertStates, component)
				return StateChange{
					ShouldAlert:   true,
					IsRecovery:    true,
					PreviousLevel: previousLevel,
				}
			}
			delete(m.alertStates, component)
		}
		return StateChange{ShouldAlert: false}
	}

	now := time.Now()
	currentState := m.alertStates[component]

	if currentState == nil || currentState.Level != *level {
		// Preserve AlertTriggered if transitioning to a lower severity
		preserveTriggered := false
		if currentState != nil && currentState.AlertTriggered {
			// Check if we're going to a lower severity
			if isLowerSeverity(*level, currentState.Level) {
				preserveTriggered = true
			}
		}

		m.alertStates[component] = &models.AlertState{
			Level:          *level,
			StartTime:      now,
			AlertTriggered: preserveTriggered,
		}

		// If duration is 0 OR we preserved an already-triggered alert, send immediately
		if duration <= 0 || preserveTriggered {
			m.alertStates[component].AlertTriggered = true
			return StateChange{
				ShouldAlert: true,
				IsRecovery:  false,
				Level:       *level,
			}
		}

		slog.Debug("Detected alert, waiting for duration",
			"component", component,
			"level", *level,
			"duration", duration)
		return StateChange{ShouldAlert: false}
	}

	elapsed := now.Sub(currentState.StartTime)
	if elapsed >= time.Duration(duration)*time.Second && !currentState.AlertTriggered {
		m.alertStates[component].AlertTriggered = true
		return StateChange{
			ShouldAlert: true,
			IsRecovery:  false,
			Level:       *level,
		}
	}

	return StateChange{ShouldAlert: false}
}

// isLowerSeverity returns true if newLevel is less severe than oldLevel
func isLowerSeverity(newLevel, oldLevel models.Severity) bool {
	severityOrder := map[models.Severity]int{
		models.SeverityWarning:  1,
		models.SeverityCritical: 2,
	}
	return severityOrder[newLevel] < severityOrder[oldLevel]
}

// triggerAlert sends an alert with rate limiting
func (m *Monitor) triggerAlert(component string, level models.Severity, value string) {
	currentTime := time.Now()
	lastTime := m.lastAlert[component]

	cooldown := m.config.Cooldown
	shouldAlert := false

	if cooldown < 0 {
		// "Alert Once" Mode
		state := m.alertStates[component]
		if state != nil {
			if lastTime.Before(state.StartTime) {
				shouldAlert = true
			}
		}
	} else {
		// Classic Mode
		if currentTime.Sub(lastTime) > time.Duration(cooldown)*time.Second {
			shouldAlert = true
		}
	}

	if shouldAlert {
		slog.Info("ALERT",
			"component", component,
			"level", level,
			"value", value)
		m.alertManager.SendAlert(component, level, value)
		m.lastAlert[component] = currentTime
	} else {
		if cooldown >= 0 {
			slog.Debug("Alert suppressed (cooldown)", "component", component)
		} else {
			slog.Debug("Alert suppressed (already sent)", "component", component)
		}
	}
}

// triggerRecovery sends a recovery notification (no cooldown)
func (m *Monitor) triggerRecovery(component string, previousLevel models.Severity, value string) {
	if !m.config.Alerts.SendRecovery {
		slog.Debug("Recovery notification disabled", "component", component)
		return
	}

	slog.Info("RECOVERY",
		"component", component,
		"previous_level", previousLevel,
		"value", value)
	m.alertManager.SendRecovery(component, previousLevel, value)

	// Clear the lastAlert time so future alerts aren't affected
	delete(m.lastAlert, component)
}

// Run starts the monitoring loop
func (m *Monitor) Run(ctx context.Context) {
	slog.Info("Starting TinyMonitor...")

	ticker := time.NewTicker(time.Duration(m.config.Refresh) * time.Second)
	defer ticker.Stop()

	// Run initial check immediately
	m.runChecks()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping TinyMonitor...")
			m.alertManager.Shutdown()
			return
		case <-ticker.C:
			m.runChecks()
		}
	}
}

func (m *Monitor) runChecks() {
	for _, collector := range m.collectors {
		results := collector.Check()
		for _, result := range results {
			duration := collector.Duration()
			change := m.processState(result.Component, result.Level, result.Value, duration)

			if change.ShouldAlert {
				if change.IsRecovery {
					m.triggerRecovery(result.Component, change.PreviousLevel, result.Value)
				} else {
					m.triggerAlert(result.Component, change.Level, result.Value)
				}
			}
		}
	}
}
