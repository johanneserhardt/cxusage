package commands

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strings"

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

    // Optional: Analyze a sample of files to report explicit vs estimated usage presence
    if len(files) > 0 {
        explicit, estimated, costPresent, totalMessages := analyzeUsageQuality(files)

        fmt.Println("\n📊 Usage Data Quality:")
        fmt.Printf("   Messages analyzed: %d\n", totalMessages)
        fmt.Printf("   With explicit usage: %d\n", explicit)
        fmt.Printf("   Estimated (no usage in logs): %d\n", estimated)
        fmt.Printf("   With cost present: %d\n", costPresent)
    }

    fmt.Println("\n🎉 Codex CLI validation complete!")
    return nil
}

func init() {
    rootCmd.AddCommand(validateCmd)
}

// analyzeUsageQuality scans a subset of log files and reports presence of explicit usage/cost
func analyzeUsageQuality(files []string) (explicit int, estimated int, costPresent int, totalMessages int) {
    // Limit scan to a reasonable number of files to keep validate fast
    maxFiles := 50
    if len(files) < maxFiles {
        maxFiles = len(files)
    }

    for i := 0; i < maxFiles; i++ {
        f, err := os.Open(files[i])
        if err != nil {
            continue
        }
        scanner := bufio.NewScanner(f)
        // Increase buffer for long lines
        const maxCapacity = 1024 * 1024
        buf := make([]byte, maxCapacity)
        scanner.Buffer(buf, maxCapacity)

        for scanner.Scan() {
            line := strings.TrimSpace(scanner.Text())
            if line == "" || !strings.HasPrefix(line, "{") {
                continue
            }
            var raw map[string]interface{}
            if err := json.Unmarshal([]byte(line), &raw); err != nil {
                continue
            }

            if !looksLikeMessage(raw) {
                continue
            }
            totalMessages++

            if hasExplicitUsage(raw) {
                explicit++
            } else {
                estimated++
            }

            if hasCost(raw) {
                costPresent++
            }
        }
        f.Close()
    }

    return explicit, estimated, costPresent, totalMessages
}

func looksLikeMessage(raw map[string]interface{}) bool {
    if t, ok := raw["type"].(string); ok && t == "message" {
        return true
    }
    if msg, ok := raw["message"].(map[string]interface{}); ok {
        // Consider it a message if it has content or usage
        if _, ok := msg["usage"]; ok {
            return true
        }
        if _, ok := msg["content"]; ok {
            return true
        }
    }
    return false
}

func hasExplicitUsage(raw map[string]interface{}) bool {
    if usage, ok := raw["usage"].(map[string]interface{}); ok {
        if hasAnyTokenKeys(usage) {
            return true
        }
    }
    if msg, ok := raw["message"].(map[string]interface{}); ok {
        if usage, ok := msg["usage"].(map[string]interface{}); ok {
            if hasAnyTokenKeys(usage) {
                return true
            }
        }
    }
    return false
}

func hasAnyTokenKeys(usage map[string]interface{}) bool {
    keys := []string{"input_tokens", "output_tokens", "prompt_tokens", "completion_tokens", "total_tokens"}
    for _, k := range keys {
        if _, ok := usage[k]; ok {
            return true
        }
    }
    return false
}

func hasCost(raw map[string]interface{}) bool {
    if _, ok := raw["costUSD"].(float64); ok {
        return true
    }
    if _, ok := raw["cost"].(float64); ok {
        return true
    }
    if msg, ok := raw["message"].(map[string]interface{}); ok {
        if _, ok := msg["costUSD"].(float64); ok {
            return true
        }
    }
    return false
}
