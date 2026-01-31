package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
)

// GotifyProvider sends alerts to Gotify
type GotifyProvider struct {
	BaseProvider
	url   string
	token string
}

// NewGotifyProvider creates a new Gotify provider
func NewGotifyProvider(cfg config.GotifyConfig) *GotifyProvider {
	return &GotifyProvider{
		BaseProvider: BaseProvider{
			ProviderName: "gotify",
			Enabled:      cfg.Enabled,
			Levels:       cfg.Levels,
			Rules:        cfg.Rules,
		},
		url:   cfg.URL,
		token: cfg.Token,
	}
}

// Send sends an alert to Gotify
func (p *GotifyProvider) Send(alert models.Alert) error {
	if p.url == "" || p.token == "" {
		return fmt.Errorf("no url or token provided")
	}

	// Ensure URL ends with /message
	serverURL := p.url
	if !strings.HasSuffix(serverURL, "/message") {
		if !strings.HasSuffix(serverURL, "/") {
			serverURL += "/"
		}
		serverURL += "message"
	}

	// Priority mapping
	var priority int
	switch alert.Level {
	case models.SeverityCritical:
		priority = 8
	case models.SeverityWarning:
		priority = 5
	case models.SeverityRecovery:
		priority = 3
	default:
		priority = 2
	}

	// System Info
	hostname := utils.GetHostname()
	executionTime := time.Now().Format("2006-01-02 15:04:05")
	ipPrivate := utils.GetPrivateIP()
	ipPublic := utils.GetPublicIP()
	loadAvg := utils.GetLoadAvg()
	uptimePretty := utils.GetUptime()

	// Enriched message (Markdown supported)
	fullMessage := fmt.Sprintf(`**Component** : %s
**Value**     : %s
**Level**     : %s

__Machine Context__
🖥️ **Server**    : `+"`%s`"+`
🏠 **Private IP**: `+"`%s`"+`
🌍 **Public IP** : `+"`%s`"+`
⚙️ **Load Avg**  : `+"`%s`"+`
⏱️ **Uptime**    : `+"`%s`"+`
🕒 **Time**      : %s`,
		alert.Component, alert.Value, alert.Level,
		hostname, ipPrivate, ipPublic, loadAvg, uptimePretty, executionTime)

	payload := map[string]interface{}{
		"title":    alert.Title,
		"message":  fullMessage,
		"priority": priority,
		"extras": map[string]interface{}{
			"client::display": map[string]string{
				"contentType": "text/markdown",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", serverURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", p.token)

	client := &http.Client{Timeout: 10 * time.Second}
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
