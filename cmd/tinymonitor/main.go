package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/monitor"
)

const Version = "1.0.0"

func main() {
	// Parse command line flags
	configPath := flag.String("c", "", "Path to configuration file")
	flag.StringVar(configPath, "config", "", "Path to configuration file")
	showVersion := flag.Bool("v", false, "Show version")
	flag.BoolVar(showVersion, "version", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("TinyMonitor %s\n", Version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Configure logging
	var handler slog.Handler
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open log file %s: %v\n", cfg.LogFile, err)
			handler = slog.NewTextHandler(os.Stdout, logOpts)
		} else {
			// Use a multi-writer to log to both stdout and file
			handler = slog.NewTextHandler(os.Stdout, logOpts)
			fileHandler := slog.NewTextHandler(logFile, logOpts)
			handler = &multiHandler{handlers: []slog.Handler{handler, fileHandler}}
		}
	} else {
		handler = slog.NewTextHandler(os.Stdout, logOpts)
	}

	slog.SetDefault(slog.New(handler))

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Start monitor
	mon := monitor.New(cfg)
	mon.Run(ctx)
}

// multiHandler is a simple handler that writes to multiple handlers
type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}
