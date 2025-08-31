package commands

import (
	"fmt"

	"github.com/johanneserhardt/cxusage/internal/codex"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Codex CLI installation and configuration",
	Long: `Validate your Codex CLI installation and configuration.
This command checks if Codex CLI is properly installed and has usage data available.`,
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	logger.Info("Validating Codex CLI installation")

	// Check if Codex directory exists
	fmt.Println("🔍 Checking Codex CLI installation...")
	exists, err := codex.CodexDirExists(cfg)
	if err != nil {
		fmt.Printf("❌ Failed to check Codex directory: %v\n", err)
		return err
	}
	
	if !exists {
		fmt.Println("❌ Codex CLI directory not found")
		fmt.Println("   Make sure Codex CLI is installed and you have used it at least once.")
		fmt.Println("   Expected location: ~/.codex")
		return fmt.Errorf("Codex CLI not found")
	}
	fmt.Println("✅ Codex CLI directory found")

	// Get Codex paths
	paths, err := codex.GetCodexPaths(cfg)
	if err != nil {
		fmt.Printf("❌ Failed to get Codex paths: %v\n", err)
		return err
	}

	// Check for usage log files
	fmt.Println("🔍 Checking for usage log files...")
	files, err := codex.GetUsageLogFiles(cfg)
	if err != nil {
		fmt.Printf("⚠️ Warning: Could not check for log files: %v\n", err)
	} else if len(files) == 0 {
		fmt.Println("⚠️ No usage log files found")
		fmt.Println("   Usage data will be available after you use Codex CLI.")
	} else {
		fmt.Printf("✅ Found %d usage log files\n", len(files))
	}

	// Show configuration
	fmt.Println("\n⚙️ Configuration:")
	fmt.Printf("   Codex Directory: %s\n", paths.ConfigDir)
	fmt.Printf("   Config File: %s\n", paths.ConfigFile)
	fmt.Printf("   Instructions File: %s\n", paths.InstructionsFile)
	fmt.Printf("   Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   Local Logging: %v\n", cfg.LocalLogging)

	fmt.Println("\n🎉 Codex CLI validation complete!")
	return nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
}