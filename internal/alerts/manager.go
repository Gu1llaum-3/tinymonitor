package alerts

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

// Manager distributes alerts to configured providers
type Manager struct {
	providers []Provider
	alertChan chan alertTask
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

type alertTask struct {
	provider Provider
	alert    models.Alert
}

// NewManager creates a new alert manager
func NewManager(cfg config.AlertsConfig) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		providers: make([]Provider, 0),
		alertChan: make(chan alertTask, 100),
		ctx:       ctx,
		cancel:    cancel,
	}

	m.loadProviders(cfg)
	m.startWorkers(5)

	return m
}

func (m *Manager) loadProviders(cfg config.AlertsConfig) {
	// Google Chat
	if cfg.GoogleChat.Enabled {
		m.providers = append(m.providers, NewGoogleChatProvider(cfg.GoogleChat))
		slog.Info("Alert Provider loaded: Google Chat")
	}

	// Ntfy
	if cfg.Ntfy.Enabled {
		m.providers = append(m.providers, NewNtfyProvider(cfg.Ntfy))
		slog.Info("Alert Provider loaded: Ntfy")
	}

	// SMTP
	if cfg.SMTP.Enabled {
		m.providers = append(m.providers, NewSMTPProvider(cfg.SMTP))
		slog.Info("Alert Provider loaded: SMTP")
	}

	// Webhook
	if cfg.Webhook.Enabled {
		m.providers = append(m.providers, NewWebhookProvider(cfg.Webhook))
		slog.Info("Alert Provider loaded: Webhook")
	}

	// Gotify
	if cfg.Gotify.Enabled {
		m.providers = append(m.providers, NewGotifyProvider(cfg.Gotify))
		slog.Info("Alert Provider loaded: Gotify")
	}
}

func (m *Manager) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		m.wg.Add(1)
		go m.worker()
	}
}

func (m *Manager) worker() {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return
		case task, ok := <-m.alertChan:
			if !ok {
				return
			}
			if err := task.provider.Send(task.alert); err != nil {
				slog.Error("Failed to send alert",
					"provider", task.provider.Name(),
					"error", err)
			}
		}
	}
}

// SendAlert distributes an alert to all configured and eligible providers
func (m *Manager) SendAlert(component string, level models.Severity, value string) {
	alert := models.NewAlert(component, level, value)

	for _, provider := range m.providers {
		if provider.ShouldSend(component, level) {
			slog.Info("Triggering alert",
				"provider", provider.Name(),
				"component", component,
				"level", level)

			select {
			case m.alertChan <- alertTask{provider: provider, alert: alert}:
			default:
				slog.Warn("Alert channel full, dropping alert",
					"provider", provider.Name(),
					"component", component)
			}
		}
	}
}

// SendRecovery distributes a recovery notification to all configured providers
func (m *Manager) SendRecovery(component string, previousLevel models.Severity, value string) {
	alert := models.NewRecoveryAlert(component, previousLevel, value)

	for _, provider := range m.providers {
		// Send recovery to providers that would have received the original alert
		if provider.ShouldSend(component, previousLevel) {
			slog.Info("Triggering recovery",
				"provider", provider.Name(),
				"component", component,
				"previous_level", previousLevel)

			select {
			case m.alertChan <- alertTask{provider: provider, alert: alert}:
			default:
				slog.Warn("Alert channel full, dropping recovery",
					"provider", provider.Name(),
					"component", component)
			}
		}
	}
}

// Shutdown gracefully shuts down the manager
func (m *Manager) Shutdown() {
	m.cancel()
	close(m.alertChan)
	m.wg.Wait()
}
