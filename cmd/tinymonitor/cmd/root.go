package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gu1llaum-3/tinymonitor/internal/config"
	"github.com/Gu1llaum-3/tinymonitor/internal/monitor"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "tinymonitor",
	Short: "Lightweight system monitoring agent",
	Long: `TinyMonitor is a lightweight system monitoring agent written in Go.

It monitors CPU, memory, disk, load average, and I/O, sending alerts
via multiple channels (Ntfy, Google Chat, SMTP, Webhooks, Gotify).`,
	Run: runMonitor,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to configuration file")
}

func runMonitor(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Configuration error:")
		if validationErrs, ok := err.(config.ValidationErrors); ok {
			for _, e := range validationErrs {
				fmt.Fprintf(os.Stderr, "  - %s: %s\n", e.Field, e.Message)
			}
			fmt.Fprintln(os.Stderr, "\nRun 'tinymonitor validate -c <file>' for details.")
		} else {
			fmt.Fprintf(os.Stderr, "  %v\n", err)
		}
		os.Exit(1)
	}

	// Configure logging
	setupLogging(cfg)

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

func setupLogging(cfg *config.Config) {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler

	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not open log file %s: %v\n", cfg.LogFile, err)
			handler = slog.NewTextHandler(os.Stdout, logOpts)
		} else {
			handler = &multiHandler{
				handlers: []slog.Handler{
					slog.NewTextHandler(os.Stdout, logOpts),
					slog.NewTextHandler(logFile, logOpts),
				},
			}
		}
	} else {
		handler = slog.NewTextHandler(os.Stdout, logOpts)
	}

	slog.SetDefault(slog.New(handler))
}

// multiHandler writes to multiple slog handlers
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
