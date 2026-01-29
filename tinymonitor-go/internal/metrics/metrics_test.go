package metrics

import (
	"testing"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/models"
)

func TestCPUCollector(t *testing.T) {
	cfg := config.MetricConfig{
		Warning:  70,
		Critical: 90,
		Enabled:  true,
		Duration: 0,
	}

	collector := NewCPUCollector(cfg)

	if collector.Name() != "cpu" {
		t.Errorf("Expected name 'cpu', got '%s'", collector.Name())
	}

	if collector.Duration() != 0 {
		t.Errorf("Expected duration 0, got %d", collector.Duration())
	}

	results := collector.Check()
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Component != "CPU" {
		t.Errorf("Expected component 'CPU', got '%s'", results[0].Component)
	}
}

func TestMemoryCollector(t *testing.T) {
	cfg := config.MetricConfig{
		Warning:  70,
		Critical: 90,
		Enabled:  true,
		Duration: 0,
	}

	collector := NewMemoryCollector(cfg)

	if collector.Name() != "memory" {
		t.Errorf("Expected name 'memory', got '%s'", collector.Name())
	}

	results := collector.Check()
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Component != "MEMORY" {
		t.Errorf("Expected component 'MEMORY', got '%s'", results[0].Component)
	}
}

func TestDiskCollector(t *testing.T) {
	cfg := config.FilesystemConfig{
		Warning:  80,
		Critical: 90,
		Enabled:  true,
		Duration: 0,
		Exclude:  []string{},
	}

	collector := NewDiskCollector(cfg)

	if collector.Name() != "filesystem" {
		t.Errorf("Expected name 'filesystem', got '%s'", collector.Name())
	}

	results := collector.Check()
	// Should have at least one disk
	if len(results) == 0 {
		t.Log("Warning: No disk partitions found (may be expected in some environments)")
	}

	for _, r := range results {
		if r.Component[:5] != "DISK:" {
			t.Errorf("Expected component to start with 'DISK:', got '%s'", r.Component)
		}
	}
}

func TestLoadCollector(t *testing.T) {
	cfg := config.MetricConfig{
		Warning:  7.0,
		Critical: 9.0,
		Enabled:  true,
		Duration: 60,
	}

	collector := NewLoadCollector(cfg)

	if collector.Name() != "load" {
		t.Errorf("Expected name 'load', got '%s'", collector.Name())
	}

	if collector.Duration() != 60 {
		t.Errorf("Expected duration 60, got %d", collector.Duration())
	}

	results := collector.Check()
	// May be empty on Windows
	if len(results) == 1 {
		if results[0].Component != "LOAD" {
			t.Errorf("Expected component 'LOAD', got '%s'", results[0].Component)
		}
	}
}

func TestIOCollector(t *testing.T) {
	cfg := config.IOConfig{
		Enabled:  true,
		Duration: 0,
	}

	collector := NewIOCollector(cfg)

	if collector.Name() != "io" {
		t.Errorf("Expected name 'io', got '%s'", collector.Name())
	}

	// First call initializes counters
	results := collector.Check()
	// May return nil on first call
	t.Logf("First check results: %v", results)

	// Second call should return results
	results = collector.Check()
	if len(results) == 1 {
		if results[0].Component != "I/O" {
			t.Errorf("Expected component 'I/O', got '%s'", results[0].Component)
		}
	}
}

func TestIOCollectorParseThreshold(t *testing.T) {
	cfg := config.IOConfig{
		Enabled: true,
	}
	collector := NewIOCollector(cfg)

	tests := []struct {
		input    interface{}
		expected float64
	}{
		{100.0, 100.0},
		{100, 100.0},
		{"100", 100.0},
		{"1K", 1024.0},
		{"1KB", 1024.0},
		{"1M", 1024 * 1024},
		{"1MB", 1024 * 1024},
		{"1G", 1024 * 1024 * 1024},
		{"1GB", 1024 * 1024 * 1024},
	}

	for _, test := range tests {
		result := collector.parseThreshold(test.input, nil)
		if result != test.expected {
			t.Errorf("parseThreshold(%v) = %v, expected %v", test.input, result, test.expected)
		}
	}

	// Test percentage with max value
	maxVal := 1000.0
	result := collector.parseThreshold("50%", &maxVal)
	if result != 500.0 {
		t.Errorf("parseThreshold('50%%', 1000) = %v, expected 500", result)
	}
}

func TestRebootCollector(t *testing.T) {
	cfg := config.RebootConfig{
		Enabled:  true,
		Duration: 0,
	}

	collector := NewRebootCollector(cfg)

	if collector.Name() != "reboot" {
		t.Errorf("Expected name 'reboot', got '%s'", collector.Name())
	}

	results := collector.Check()
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Component != "REBOOT" {
		t.Errorf("Expected component 'REBOOT', got '%s'", results[0].Component)
	}
}

func TestCollectorInterface(t *testing.T) {
	// Verify all collectors implement the Collector interface
	var _ Collector = (*CPUCollector)(nil)
	var _ Collector = (*MemoryCollector)(nil)
	var _ Collector = (*DiskCollector)(nil)
	var _ Collector = (*LoadCollector)(nil)
	var _ Collector = (*IOCollector)(nil)
	var _ Collector = (*RebootCollector)(nil)
}

func TestSeverityLevels(t *testing.T) {
	// Test that severity detection works correctly
	cfg := config.MetricConfig{
		Warning:  50,
		Critical: 80,
		Enabled:  true,
	}

	// The actual metrics would need to be mocked for proper threshold testing
	// This test just verifies the severity constants are correct
	if models.SeverityWarning != "WARNING" {
		t.Errorf("Expected SeverityWarning='WARNING', got '%s'", models.SeverityWarning)
	}
	if models.SeverityCritical != "CRITICAL" {
		t.Errorf("Expected SeverityCritical='CRITICAL', got '%s'", models.SeverityCritical)
	}

	_ = cfg // Use cfg to avoid unused variable warning
}
