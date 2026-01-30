package alerts

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
	"github.com/Gu1llaum-3/tinymonitor/internal/utils"
)

// SMTPProvider sends alerts via email
type SMTPProvider struct {
	BaseProvider
	host     string
	port     int
	user     string
	password string
	fromAddr string
	toAddrs  []string
	useTLS   bool
}

// NewSMTPProvider creates a new SMTP provider
func NewSMTPProvider(cfg config.SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		BaseProvider: BaseProvider{
			ProviderName: "smtp",
			Enabled:      cfg.Enabled,
			Levels:       cfg.Levels,
			Rules:        cfg.Rules,
		},
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		fromAddr: cfg.FromAddr,
		toAddrs:  cfg.ToAddrs,
		useTLS:   cfg.UseTLS,
	}
}

// Send sends an alert via email
func (p *SMTPProvider) Send(alert models.Alert) error {
	if p.host == "" || p.user == "" || p.password == "" || p.fromAddr == "" || len(p.toAddrs) == 0 {
		return fmt.Errorf("missing SMTP configuration")
	}

	// System Info
	hostname := utils.GetHostname()
	executionTime := time.Now().Format("2006-01-02 15:04:05")
	ipPrivate := utils.GetPrivateIP()
	ipPublic := utils.GetPublicIP()
	loadAvg := utils.GetLoadAvg()
	uptimePretty := utils.GetUptime()

	subject := fmt.Sprintf("[%s] %s on %s - %s", alert.Level, alert.Component, hostname, alert.Value)

	htmlContent := fmt.Sprintf(`
	<html>
	<body>
		<h2>%s</h2>
		<p><strong>Component:</strong> %s</p>
		<p><strong>Value:</strong> %s</p>
		<p><strong>Level:</strong> %s</p>
		<hr>
		<h3>Machine Context</h3>
		<ul>
			<li><strong>Server:</strong> %s</li>
			<li><strong>Private IP:</strong> %s</li>
			<li><strong>Public IP:</strong> %s</li>
			<li><strong>Load Avg:</strong> %s</li>
			<li><strong>Uptime:</strong> %s</li>
			<li><strong>Time:</strong> %s</li>
		</ul>
	</body>
	</html>`,
		alert.Title, alert.Component, alert.Value, alert.Level,
		hostname, ipPrivate, ipPublic, loadAvg, uptimePretty, executionTime)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n%s",
		p.fromAddr, strings.Join(p.toAddrs, ", "), subject, htmlContent)

	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	auth := smtp.PlainAuth("", p.user, p.password, p.host)

	if p.useTLS {
		return p.sendWithTLS(addr, auth, msg)
	}

	return smtp.SendMail(addr, auth, p.fromAddr, p.toAddrs, []byte(msg))
}

func (p *SMTPProvider) sendWithTLS(addr string, auth smtp.Auth, msg string) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: p.host})
	if err != nil {
		// Try STARTTLS instead
		return p.sendWithSTARTTLS(addr, auth, msg)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(p.fromAddr); err != nil {
		return err
	}

	for _, toAddr := range p.toAddrs {
		if err := client.Rcpt(toAddr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	LogInfo(p.ProviderName, "Email sent", "to", p.toAddrs)
	return client.Quit()
}

func (p *SMTPProvider) sendWithSTARTTLS(addr string, auth smtp.Auth, msg string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Hello("localhost"); err != nil {
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: p.host}
		if err := client.StartTLS(config); err != nil {
			return err
		}
	}

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(p.fromAddr); err != nil {
		return err
	}

	for _, toAddr := range p.toAddrs {
		if err := client.Rcpt(toAddr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	LogInfo(p.ProviderName, "Email sent", "to", p.toAddrs)
	return client.Quit()
}
