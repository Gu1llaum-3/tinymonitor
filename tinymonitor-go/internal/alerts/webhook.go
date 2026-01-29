package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
)

// WebhookProvider sends alerts to a generic webhook
type WebhookProvider struct {
	BaseProvider
	url     string
	headers map[string]string
	timeout int
}

// NewWebhookProvider creates a new webhook provider
func NewWebhookProvider(cfg config.WebhookConfig) *WebhookProvider {
	return &WebhookProvider{
		BaseProvider: BaseProvider{
			ProviderName: "webhook",
			Enabled:      cfg.Enabled,
			Levels:       cfg.Levels,
			Rules:        cfg.Rules,
		},
		url:     cfg.URL,
		headers: cfg.Headers,
		timeout: cfg.Timeout,
	}
}

// Send sends an alert to the webhook
func (p *WebhookProvider) Send(alert models.Alert) error {
	if p.url == "" {
		return fmt.Errorf("no url provided")
	}

	payload := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"alert": map[string]interface{}{
			"level":     alert.Level,
			"component": alert.Component,
			"value":     alert.Value,
			"title":     alert.Title,
			"message":   alert.Message,
		},
		"host": map[string]interface{}{
			"hostname":     utils.GetHostname(),
			"ip_private":   utils.GetPrivateIP(),
			"ip_public":    utils.GetPublicIP(),
			"uptime":       utils.GetUptime(),
			"load_average": utils.GetLoadAvg(),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", p.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Set headers
	hasContentType := false
	for key, value := range p.headers {
		req.Header.Set(key, value)
		if key == "Content-Type" {
			hasContentType = true
		}
	}
	if !hasContentType {
		req.Header.Set("Content-Type", "application/json")
	}

	timeout := p.timeout
	if timeout <= 0 {
		timeout = 10
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send alert: status %d", resp.StatusCode)
	}

	LogInfo(p.ProviderName, "Alert sent successfully")
	return nil
}
