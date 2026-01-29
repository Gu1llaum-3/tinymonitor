package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
)

// GetHostname returns the system hostname
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// GetPrivateIP returns the private IP address
func GetPrivateIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetPublicIP returns the public IP address
func GetPublicIP() string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "N/A"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "N/A"
	}

	return string(body)
}

// GetLoadAvg returns the load average as a formatted string
func GetLoadAvg() string {
	avg, err := load.Avg()
	if err != nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f, %.2f, %.2f", avg.Load1, avg.Load5, avg.Load15)
}

// GetUptime returns the system uptime as a formatted string
func GetUptime() string {
	bootTime, err := host.BootTime()
	if err != nil {
		return "N/A"
	}

	uptimeSeconds := time.Now().Unix() - int64(bootTime)
	hours := uptimeSeconds / 3600
	minutes := (uptimeSeconds % 3600) / 60

	return fmt.Sprintf("%dh %dm", hours, minutes)
}
