package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/johanneserhardt/cxusage/internal/utils"
)

var dailyCmd = &cobra.Command{
	Use:   "daily [days]",
	Short: "Show daily usage reports",
	Long: `Display daily usage reports for OpenAI API.
By default shows the last 7 days. You can specify a different number of days.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDaily,
}

func runDaily(cmd *cobra.Command, args []string) error {
	days := 7 // default
	if len(args) > 0 {
		var err error
		days, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid number of days: %s", args[0])
		}
		if days < 1 || days > 365 {
			return fmt.Errorf("days must be between 1 and 365")
		}
	}

	// Get flags
	outputFormat, _ := cmd.Flags().GetString("output")
	offline, _ := cmd.Flags().GetBool("offline")
	
	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	logger.WithFields(map[string]interface{}{
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"days":       days,
		"offline":    offline,
	}).Info("Generating daily usage report")

	var dailyUsage []types.DailyUsage
	var err error

	// Load from Codex CLI local files (no API needed)
	dailyUsage, err = utils.LoadDailyUsageFromCodex(cfg, startDate, endDate, logger)

	if err != nil {
		return fmt.Errorf("failed to load daily usage data: %w", err)
	}

	// Handle empty data with helpful message
	if len(dailyUsage) == 0 {
		if outputFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Printf("%s\n", utils.Yellow("No Codex CLI usage data found"))
			fmt.Println()
			fmt.Printf("Try:\n")
			fmt.Printf("• %s - Check if Codex CLI is set up\n", utils.Cyan("cxusage validate"))
			fmt.Printf("• Use Codex CLI first, then run %s\n", utils.Cyan("cxusage daily"))
			fmt.Printf("• Run %s to see sample output\n", utils.Cyan("cxusage demo"))
		}
		return nil
	}

	// Output results
	switch types.OutputFormat(outputFormat) {
	case types.OutputFormatJSON:
		return outputDailyJSON(dailyUsage)
	case types.OutputFormatTable:
		return outputDailyTable(dailyUsage)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputDailyJSON(dailyUsage []types.DailyUsage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(dailyUsage)
}

func outputDailyTable(dailyUsage []types.DailyUsage) error {
	utils.FormatDailyUsageTable(dailyUsage)
	return nil
}

func init() {
	rootCmd.AddCommand(dailyCmd)
	
	// Daily-specific flags
	dailyCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
	dailyCmd.Flags().String("end-date", "", "End date (YYYY-MM-DD)")
	dailyCmd.Flags().StringSlice("models", []string{}, "Filter by specific models")
}