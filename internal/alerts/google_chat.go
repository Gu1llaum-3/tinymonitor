package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
)

// GoogleChatProvider sends alerts to Google Chat
type GoogleChatProvider struct {
	BaseProvider
	webhookURL string
}

// NewGoogleChatProvider creates a new Google Chat provider
func NewGoogleChatProvider(cfg config.GoogleChatConfig) *GoogleChatProvider {
	return &GoogleChatProvider{
		BaseProvider: BaseProvider{
			ProviderName: "google_chat",
			Enabled:      cfg.Enabled,
			Levels:       cfg.Levels,
			Rules:        cfg.Rules,
		},
		webhookURL: cfg.WebhookURL,
	}
}

// Send sends an alert to Google Chat
func (p *GoogleChatProvider) Send(alert models.Alert) error {
	if p.webhookURL == "" {
		return fmt.Errorf("no webhook_url provided")
	}

	// Visual decoration based on status
	var icon, fontColor, titleText string
	switch alert.Level {
	case models.SeverityCritical:
		icon = "🚨"
		fontColor = "#FF0000"
		titleText = "CRITICAL ALERT : " + alert.Component
	case models.SeverityWarning:
		icon = "⚠️"
		fontColor = "#FFA500"
		titleText = "WARNING : " + alert.Component
	case models.SeverityRecovery:
		icon = "✅"
		fontColor = "#00AA00"
		titleText = "RECOVERED : " + alert.Component
	default:
		icon = "ℹ️"
		fontColor = "#000000"
		titleText = "INFO : " + alert.Component
	}

	// System Info
	hostname := utils.GetHostname()
	executionTime := time.Now().Format("2006-01-02 15:04:05")
	ipPrivate := utils.GetPrivateIP()
	ipPublic := utils.GetPublicIP()
	loadAvg := utils.GetLoadAvg()
	uptimePretty := utils.GetUptime()

	// Sanitize for cardId
	safeHostname := sanitizeForCardID(hostname)
	safeComponent := sanitizeForCardID(alert.Component)

	payload := map[string]interface{}{
		"cardsV2": []map[string]interface{}{
			{
				"cardId": fmt.Sprintf("tinymonitor-%s-%s", safeHostname, safeComponent),
				"card": map[string]interface{}{
					"header": map[string]interface{}{
						"title":     fmt.Sprintf("%s %s", icon, titleText),
						"subtitle":  fmt.Sprintf("Server : %s", hostname),
						"imageUrl":  "https://upload.wikimedia.org/wikipedia/commons/thumb/3/35/Tux.svg/1200px-Tux.svg.png",
						"imageType": "CIRCLE",
					},
					"sections": []map[string]interface{}{
						{
							"header": "Incident Details",
							"widgets": []map[string]interface{}{
								{
									"decoratedText": map[string]interface{}{
										"topLabel":  "Monitored Component",
										"text":      fmt.Sprintf("<b>%s</b>", alert.Component),
										"startIcon": map[string]string{"knownIcon": "MEMBERSHIP"},
									},
								},
								{
									"decoratedText": map[string]interface{}{
										"topLabel":  "Current Value",
										"text":      fmt.Sprintf("<font color=\"%s\"><b>%s</b></font>", fontColor, alert.Value),
										"startIcon": map[string]string{"knownIcon": "DESCRIPTION"},
									},
								},
								{
									"decoratedText": map[string]interface{}{
										"topLabel":  "Alert Level",
										"text":      fmt.Sprintf("<b>%s</b>", alert.Level),
										"startIcon": map[string]string{"knownIcon": "STAR"},
									},
								},
							},
						},
						{
							"header":                    "Machine Context",
							"collapsible":               true,
							"uncollapsibleWidgetsCount": 2,
							"widgets": []map[string]interface{}{
								{
									"textParagraph": map[string]string{
										"text": fmt.Sprintf("<b>Private IP:</b> %s<br><b>Public IP:</b> %s", ipPrivate, ipPublic),
									},
								},
								{
									"textParagraph": map[string]string{
										"text": fmt.Sprintf("<b>Load:</b> %s<br><b>Uptime:</b> %s", loadAvg, uptimePretty),
									},
								},
								{
									"textParagraph": map[string]string{
										"text": fmt.Sprintf("<font color=\"#808080\">Alert Time: %s</font>", executionTime),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", p.webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

func sanitizeForCardID(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return re.ReplaceAllString(s, "_")
}
