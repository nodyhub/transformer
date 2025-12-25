// Package cmd implements the command line interface for the transformer tool.
package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	logLevel string
	logFile  string
	verbose  bool
)

var rootCmd = &cobra.Command{
	Use:   "transformer",
	Short: "Execute a series of steps defined in a step file.",
	Long: `Transformer is a CLI tool that allows you to define and execute a series of steps
in a structured manner. You can create step files in YAML format to outline the
steps you want to perform, and Transformer will execute them sequentially.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		return setupLogger()
	},
	SilenceUsage: true,
}

// Execute runs the root command, which in turn executes the appropriate subcommand.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setupLogger() error {
	// Parse log level
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Determine output writer
	var writer *os.File
	switch logFile {
	case "", "stderr":
		writer = os.Stderr
	case "stdout":
		writer = os.Stdout
	default:
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		writer = f
	}

	// Configure and set default logger
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set the logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Set the log file path (default is stderr; use 'stdout' for standard output)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging (equivalent to --log-level=debug)")

	// If verbose is set, override log level to debug
	cobra.OnInitialize(func() {
		if verbose {
			logLevel = "debug"
		}
	})
}
