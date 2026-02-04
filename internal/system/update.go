package system

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	GitHubRepo   = "Gu1llaum-3/tinymonitor"
	GitHubAPIURL = "https://api.github.com/repos/Gu1llaum-3/tinymonitor/releases/latest"
	DownloadURL  = "https://github.com/Gu1llaum-3/tinymonitor/releases/download"
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// GetLatestVersion fetches the latest version from GitHub API
func GetLatestVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(GitHubAPIURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	return release.TagName, nil
}

// CompareVersions returns true if latest is newer than current
// Handles versions like "v1.2.0", "v1.2.0-dirty", "dev"
func CompareVersions(current, latest string) bool {
	// Clean versions (remove 'v' prefix if present)
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// If current is "dev" or contains "dirty", always consider update available
	if current == "dev" || strings.Contains(current, "dirty") {
		return true
	}

	// Simple string comparison - works for semver
	// For more complex cases, use a semver library
	return latest > current
}

// GetChangelogURL returns the URL to the release page
func GetChangelogURL(version string) string {
	return fmt.Sprintf("https://github.com/%s/releases/tag/%s", GitHubRepo, version)
}

// DownloadBinary downloads the binary for the current OS/arch and returns the temp file path
func DownloadBinary(version string) (string, error) {
	// Determine OS and arch
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Format OS name (capitalize first letter)
	osFormatted := strings.ToUpper(osName[:1]) + osName[1:]

	// Format arch (amd64 -> x86_64)
	archFormatted := arch
	if arch == "amd64" {
		archFormatted = "x86_64"
	}

	// Build download URL
	archiveName := fmt.Sprintf("tinymonitor_%s_%s.tar.gz", osFormatted, archFormatted)
	downloadURL := fmt.Sprintf("%s/%s/%s", DownloadURL, version, archiveName)

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "tinymonitor-update-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Download archive
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Extract binary from tar.gz
	binaryPath, err := extractBinary(resp.Body, tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to extract binary: %w", err)
	}

	return binaryPath, nil
}

// extractBinary extracts the tinymonitor binary from a tar.gz stream
func extractBinary(r io.Reader, destDir string) (string, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the tinymonitor binary
		if header.Typeflag == tar.TypeReg && filepath.Base(header.Name) == "tinymonitor" {
			destPath := filepath.Join(destDir, "tinymonitor")
			outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()

			return destPath, nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

// InstallBinary replaces the current binary with the new one
func InstallBinary(newBinaryPath string) error {
	// Verify the new binary works
	cmd := exec.Command(newBinaryPath, "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("new binary verification failed: %w", err)
	}

	// Install to BinaryPath (defined in service.go as /usr/local/bin/tinymonitor)
	// Use atomic rename instead of copy to avoid "text file busy" errors
	tmpPath := BinaryPath + ".new"

	if IsRoot() {
		// Direct copy then rename if root
		if err := copyFile(newBinaryPath, tmpPath); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}
		if err := os.Rename(tmpPath, BinaryPath); err != nil {
			os.Remove(tmpPath) // Clean up on failure
			return fmt.Errorf("failed to install binary: %w", err)
		}
		return nil
	}

	// Use sudo if not root
	// First copy to temporary location
	cmd = exec.Command("sudo", "cp", newBinaryPath, tmpPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	// Then atomically rename (mv can replace running binaries)
	cmd = exec.Command("sudo", "mv", tmpPath, BinaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		// Clean up temp file on failure
		exec.Command("sudo", "rm", "-f", tmpPath).Run()
		return fmt.Errorf("failed to install binary: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// CleanupUpdate removes temporary files from an update
func CleanupUpdate(tmpPath string) {
	if tmpPath != "" {
		// Remove the parent temp directory
		dir := filepath.Dir(tmpPath)
		if strings.HasPrefix(filepath.Base(dir), "tinymonitor-update-") {
			os.RemoveAll(dir)
		}
	}
}
