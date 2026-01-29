package alerts

import (
	"log/slog"
	"strings"

	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

// Provider defines the interface for alert providers
type Provider interface {
	// Name returns the name of the provider
	Name() string

	// Send sends an alert
	Send(alert models.Alert) error

	// ShouldSend checks if this provider should send an alert for the given component and level
	ShouldSend(component string, level models.Severity) bool
}

// BaseProvider provides common functionality for alert providers
type BaseProvider struct {
	ProviderName string
	Enabled      bool
	Levels       []string
	Rules        map[string][]string
}

// Name returns the provider name
func (p *BaseProvider) Name() string {
	return p.ProviderName
}

// ShouldSend checks if this provider should send an alert
func (p *BaseProvider) ShouldSend(component string, level models.Severity) bool {
	// 1. Global check
	if !p.Enabled {
		return false
	}

	// 2. Rules retrieval
	if p.Rules == nil || len(p.Rules) == 0 {
		// If no rules defined, use old system (or all by default)
		acceptedLevels := p.Levels
		if len(acceptedLevels) == 0 {
			acceptedLevels = []string{"WARNING", "CRITICAL"}
		}
		return contains(acceptedLevels, string(level))
	}

	// 3. Component logic
	configKey := normalizeComponentName(component)

	// Look for specific rule, else default, else refuse for safety
	allowedLevels, ok := p.Rules[configKey]
	if !ok {
		allowedLevels = p.Rules["default"]
	}

	return contains(allowedLevels, string(level))
}

// normalizeComponentName converts the technical component name to a configuration key
func normalizeComponentName(component string) string {
	if strings.HasPrefix(component, "DISK:") {
		return "filesystem"
	}
	if component == "LOAD" {
		return "load"
	}
	return strings.ToLower(component)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// LogError logs an error for a provider
func LogError(providerName string, msg string, args ...any) {
	slog.Error("["+providerName+"] "+msg, args...)
}

// LogInfo logs info for a provider
func LogInfo(providerName string, msg string, args ...any) {
	slog.Info("["+providerName+"] "+msg, args...)
}
