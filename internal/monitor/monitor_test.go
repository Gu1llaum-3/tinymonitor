package monitor

import (
	"testing"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

func TestIsLowerSeverity(t *testing.T) {
	tests := []struct {
		name     string
		newLevel models.Severity
		oldLevel models.Severity
		want     bool
	}{
		{
			name:     "WARNING is lower than CRITICAL",
			newLevel: models.SeverityWarning,
			oldLevel: models.SeverityCritical,
			want:     true,
		},
		{
			name:     "CRITICAL is not lower than WARNING",
			newLevel: models.SeverityCritical,
			oldLevel: models.SeverityWarning,
			want:     false,
		},
		{
			name:     "WARNING is not lower than WARNING",
			newLevel: models.SeverityWarning,
			oldLevel: models.SeverityWarning,
			want:     false,
		},
		{
			name:     "CRITICAL is not lower than CRITICAL",
			newLevel: models.SeverityCritical,
			oldLevel: models.SeverityCritical,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLowerSeverity(tt.newLevel, tt.oldLevel)
			if got != tt.want {
				t.Errorf("isLowerSeverity(%v, %v) = %v, want %v",
					tt.newLevel, tt.oldLevel, got, tt.want)
			}
		})
	}
}

func TestProcessState_RecoveryNotifications(t *testing.T) {
	tests := []struct {
		name           string
		setupStates    func(*Monitor) // Setup initial state
		component      string
		level          *models.Severity
		value          string
		duration       int
		wantShouldAlert bool
		wantIsRecovery bool
		wantLevel      models.Severity
		description    string
	}{
		{
			name:      "CRITICAL to OK should send recovery",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				critical := models.SeverityCritical
				m.alertStates["test_metric"] = &models.AlertState{
					Level:          critical,
					StartTime:      time.Now().Add(-10 * time.Second),
					AlertTriggered: true,
				}
			},
			level:           nil, // nil = OK
			value:           "50%",
			duration:        0,
			wantShouldAlert: true,
			wantIsRecovery:  true,
			wantLevel:       models.SeverityCritical, // Previous level
			description:     "Direct transition from CRITICAL to OK should send recovery",
		},
		{
			name:      "WARNING to OK should send recovery",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				warning := models.SeverityWarning
				m.alertStates["test_metric"] = &models.AlertState{
					Level:          warning,
					StartTime:      time.Now().Add(-10 * time.Second),
					AlertTriggered: true,
				}
			},
			level:           nil,
			value:           "50%",
			duration:        0,
			wantShouldAlert: true,
			wantIsRecovery:  true,
			wantLevel:       models.SeverityWarning,
			description:     "Direct transition from WARNING to OK should send recovery",
		},
		{
			name:      "CRITICAL to WARNING preserves AlertTriggered",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				critical := models.SeverityCritical
				m.alertStates["test_metric"] = &models.AlertState{
					Level:          critical,
					StartTime:      time.Now().Add(-10 * time.Second),
					AlertTriggered: true,
				}
			},
			level:           ptrSeverity(models.SeverityWarning),
			value:           "75%",
			duration:        0,
			wantShouldAlert: true,
			wantIsRecovery:  false,
			wantLevel:       models.SeverityWarning,
			description:     "Transition to lower severity should preserve AlertTriggered and send alert immediately",
		},
		{
			name:      "WARNING to CRITICAL should not preserve AlertTriggered",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				warning := models.SeverityWarning
				m.alertStates["test_metric"] = &models.AlertState{
					Level:          warning,
					StartTime:      time.Now().Add(-10 * time.Second),
					AlertTriggered: true,
				}
			},
			level:           ptrSeverity(models.SeverityCritical),
			value:           "95%",
			duration:        0,
			wantShouldAlert: true,
			wantIsRecovery:  false,
			wantLevel:       models.SeverityCritical,
			description:     "Transition to higher severity should reset state",
		},
		{
			name:      "First alert with duration should wait",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				// No existing state
			},
			level:           ptrSeverity(models.SeverityCritical),
			value:           "95%",
			duration:        30,
			wantShouldAlert: false,
			wantIsRecovery:  false,
			description:     "New alert with duration should not trigger immediately",
		},
		{
			name:      "Alert after duration elapsed should trigger",
			component: "test_metric",
			setupStates: func(m *Monitor) {
				critical := models.SeverityCritical
				m.alertStates["test_metric"] = &models.AlertState{
					Level:          critical,
					StartTime:      time.Now().Add(-40 * time.Second),
					AlertTriggered: false,
				}
			},
			level:           ptrSeverity(models.SeverityCritical),
			value:           "95%",
			duration:        30,
			wantShouldAlert: true,
			wantIsRecovery:  false,
			wantLevel:       models.SeverityCritical,
			description:     "Alert should trigger after duration has elapsed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create monitor with minimal config
			cfg := &config.Config{
				Refresh:  5,
				Cooldown: 60,
				Alerts: config.AlertsConfig{
					SendRecovery: true,
				},
			}
			m := New(cfg)

			// Setup initial state
			if tt.setupStates != nil {
				tt.setupStates(m)
			}

			// Process state
			change := m.processState(tt.component, tt.level, tt.value, tt.duration)

			// Verify results
			if change.ShouldAlert != tt.wantShouldAlert {
				t.Errorf("%s: ShouldAlert = %v, want %v",
					tt.description, change.ShouldAlert, tt.wantShouldAlert)
			}

			if change.IsRecovery != tt.wantIsRecovery {
				t.Errorf("%s: IsRecovery = %v, want %v",
					tt.description, change.IsRecovery, tt.wantIsRecovery)
			}

			if tt.wantShouldAlert && !tt.wantIsRecovery {
				if change.Level != tt.wantLevel {
					t.Errorf("%s: Level = %v, want %v",
						tt.description, change.Level, tt.wantLevel)
				}
			}

			if tt.wantIsRecovery && change.PreviousLevel != tt.wantLevel {
				t.Errorf("%s: PreviousLevel = %v, want %v",
					tt.description, change.PreviousLevel, tt.wantLevel)
			}
		})
	}
}

func TestProcessState_MultiLevelTransition(t *testing.T) {
	// This test simulates the bug scenario: CRITICAL → WARNING → OK
	// The fix should ensure recovery is sent even when passing through WARNING

	cfg := &config.Config{
		Refresh:  5,
		Cooldown: 60,
		Alerts: config.AlertsConfig{
			SendRecovery: true,
		},
	}
	m := New(cfg)

	component := "load"
	critical := models.SeverityCritical
	warning := models.SeverityWarning

	// Step 1: Metric goes CRITICAL
	change := m.processState(component, &critical, "5.0", 30)
	if change.ShouldAlert {
		t.Error("Step 1: Should not alert immediately (duration=30)")
	}

	// Simulate time passing (duration elapsed)
	state := m.alertStates[component]
	state.StartTime = time.Now().Add(-40 * time.Second)

	// Step 2: Duration elapsed, should trigger CRITICAL alert
	change = m.processState(component, &critical, "5.0", 30)
	if !change.ShouldAlert || change.IsRecovery {
		t.Error("Step 2: Should trigger CRITICAL alert")
	}
	if !state.AlertTriggered {
		t.Error("Step 2: AlertTriggered should be true")
	}

	// Step 3: Metric drops to WARNING (this is where the bug occurred)
	change = m.processState(component, &warning, "3.0", 30)
	if !change.ShouldAlert {
		t.Error("Step 3: Should alert for WARNING level (AlertTriggered preserved)")
	}
	if change.IsRecovery {
		t.Error("Step 3: Should not be recovery (still in alert state)")
	}

	// Verify AlertTriggered is preserved
	state = m.alertStates[component]
	if !state.AlertTriggered {
		t.Fatal("Step 3: AlertTriggered should be preserved when transitioning to lower severity")
	}

	// Step 4: Metric returns to OK (before WARNING duration elapses)
	change = m.processState(component, nil, "1.0", 30)
	if !change.ShouldAlert {
		t.Error("Step 4: Should send alert (recovery)")
	}
	if !change.IsRecovery {
		t.Error("Step 4: Should be recovery")
	}
	if change.PreviousLevel != warning {
		t.Errorf("Step 4: PreviousLevel = %v, want %v", change.PreviousLevel, warning)
	}

	// Verify state is cleaned up
	if _, exists := m.alertStates[component]; exists {
		t.Error("Step 4: Alert state should be deleted after recovery")
	}
}

func TestProcessState_NoRecoveryIfNeverTriggered(t *testing.T) {
	// If an alert was detected but never triggered (duration not elapsed),
	// returning to OK should NOT send recovery

	cfg := &config.Config{
		Refresh:  5,
		Cooldown: 60,
		Alerts: config.AlertsConfig{
			SendRecovery: true,
		},
	}
	m := New(cfg)

	component := "cpu"
	critical := models.SeverityCritical

	// Step 1: Metric goes CRITICAL (with duration)
	change := m.processState(component, &critical, "95%", 30)
	if change.ShouldAlert {
		t.Error("Should not alert immediately with duration")
	}

	// Step 2: Metric returns to OK before duration elapsed
	change = m.processState(component, nil, "50%", 30)
	if change.ShouldAlert {
		t.Error("Should not send recovery if alert was never triggered")
	}
	if change.IsRecovery {
		t.Error("Should not be recovery")
	}

	// Verify state is cleaned up
	if _, exists := m.alertStates[component]; exists {
		t.Error("Alert state should be deleted")
	}
}

// ptrSeverity is a helper to create a pointer to a Severity value
func ptrSeverity(s models.Severity) *models.Severity {
	return &s
}
