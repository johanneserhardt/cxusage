package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/johanneserhardt/cxusage/internal/types"
)

const (
	DefaultLogLevel = "warn" // Reduce default log noise
	DefaultLogsDir  = "logs"
	ConfigFileName  = "cxusage"
)

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*types.Config, error) {
	// Set default values
	viper.SetDefault("log_level", DefaultLogLevel)
	viper.SetDefault("local_logging", false)
	viper.SetDefault("logs_dir", DefaultLogsDir)

	// Set config file name and type
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")

	// Add config search paths
	homeDir, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(homeDir, ".config"))
		viper.AddConfigPath(homeDir)
	}
	viper.AddConfigPath(".")

	// Enable environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvPrefix("OPENAI_USAGE")

	// Read config file (optional) - ignore errors for missing or invalid config
	if err := viper.ReadInConfig(); err != nil {
		// Only return error for serious issues, not missing/invalid config files
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Log warning but continue with defaults for invalid config files
			fmt.Fprintf(os.Stderr, "Warning: Invalid config file found, using defaults: %v\n", err)
		}
	}

	// Unmarshal config
	var config types.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// No API key validation needed for local file reading

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *types.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, ConfigFileName+".yaml")

	viper.Set("log_level", config.LogLevel)
	viper.Set("local_logging", config.LocalLogging)
	viper.Set("logs_dir", config.LogsDir)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

// GetLogsDir returns the absolute path to the logs directory
func GetLogsDir(config *types.Config) (string, error) {
	logsDir := config.LogsDir
	if !filepath.IsAbs(logsDir) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		logsDir = filepath.Join(homeDir, ".local", "share", "cxusage", logsDir)
	}

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return "", fmt.Errorf("could not create logs directory: %w", err)
	}

	return logsDir, nil
}