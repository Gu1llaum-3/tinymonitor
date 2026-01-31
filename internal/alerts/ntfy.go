package alerts

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
)

// NtfyProvider sends alerts to Ntfy
type NtfyProvider struct {
	BaseProvider
	topicURL string
	token    string
}

// NewNtfyProvider creates a new Ntfy provider
func NewNtfyProvider(cfg config.NtfyConfig) *NtfyProvider {
	return &NtfyProvider{
		BaseProvider: BaseProvider{
			ProviderName: "ntfy",
			Enabled:      cfg.Enabled,
			Levels:       cfg.Levels,
			Rules:        cfg.Rules,
		},
		topicURL: cfg.TopicURL,
		token:    cfg.Token,
	}
}

// Send sends an alert to Ntfy
func (p *NtfyProvider) Send(alert models.Alert) error {
	if p.topicURL == "" {
		return fmt.Errorf("no topic_url provided")
	}

	// Priority mapping
	var priority string
	var tags string
	switch alert.Level {
	case models.SeverityCritical:
		priority = "5"
		tags = "rotating_light,critical"
	case models.SeverityWarning:
		priority = "3"
		tags = "warning"
	case models.SeverityRecovery:
		priority = "2"
		tags = "white_check_mark,recovered"
	default:
		priority = "1"
		tags = "information_source"
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

	req, err := http.NewRequest("POST", p.topicURL, bytes.NewReader([]byte(fullMessage)))
	if err != nil {
		return err
	}

	req.Header.Set("Title", alert.Title)
	req.Header.Set("Priority", priority)
	req.Header.Set("Tags", tags)
	req.Header.Set("Markdown", "yes")

	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send alert: status %d", resp.StatusCode)
	}

	LogInfo(p.ProviderName, "Alert sent successfully")
	return nil
}
