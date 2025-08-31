package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/johanneserhardt/cxusage/internal/config"
	"github.com/johanneserhardt/cxusage/internal/types"
)

var (
	cfg    *types.Config
	logger *logrus.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cxusage",
	Short: "Beautiful usage analysis tool for Codex CLI with live monitoring",
	Long: `cxusage is a CLI tool that analyzes Codex CLI usage data from local files.
It provides beautiful reports on token usage, costs, and usage patterns with 
real-time live monitoring, similar to ccusage for Claude Code.

✨ Key Features:
• Daily/Monthly reports with gorgeous formatting
• Live monitoring dashboard (cxusage blocks --live) 
• 5-hour billing blocks tracking
• No API key needed - reads local Codex files`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for version, help, demo and completion commands
		if cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "completion" || cmd.Name() == "demo" {
			logger = logrus.New()
			logger.SetLevel(logrus.WarnLevel) // Reduce log noise for help commands
			return nil
		}

		// Initialize logger
		logger = logrus.New()
		
		// Load configuration
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Set log level (default to warn to reduce noise)
		level, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			level = logrus.WarnLevel
		}
		logger.SetLevel(level)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format (table, json)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().Bool("offline", false, "Use local logs only (no API calls)")
	
	// Bind flags to viper
	// viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
}