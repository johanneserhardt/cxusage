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

var monthlyCmd = &cobra.Command{
	Use:   "monthly [months]",
	Short: "Show monthly usage reports",
	Long: `Display monthly usage reports for OpenAI API.
By default shows the last 3 months. You can specify a different number of months.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMonthly,
}

func runMonthly(cmd *cobra.Command, args []string) error {
	months := 3 // default
	if len(args) > 0 {
		var err error
		months, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid number of months: %s", args[0])
		}
		if months < 1 || months > 24 {
			return fmt.Errorf("months must be between 1 and 24")
		}
	}

	// Get flags
	outputFormat, _ := cmd.Flags().GetString("output")
	offline, _ := cmd.Flags().GetBool("offline")
	
	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, -months, 0)

	logger.WithFields(map[string]interface{}{
		"start_date": startDate.Format("2006-01"),
		"end_date":   endDate.Format("2006-01"),
		"months":     months,
		"offline":    offline,
	}).Info("Generating monthly usage report")

	var monthlyUsage []types.MonthlyUsage
	var err error

	// Load from Codex CLI local files (no API needed)
	monthlyUsage, err = utils.LoadMonthlyUsageFromCodex(cfg, startDate, endDate, logger)

	if err != nil {
		return fmt.Errorf("failed to load monthly usage data: %w", err)
	}

	// Output results
	switch types.OutputFormat(outputFormat) {
	case types.OutputFormatJSON:
		return outputMonthlyJSON(monthlyUsage)
	case types.OutputFormatTable:
		return outputMonthlyTable(monthlyUsage)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputMonthlyJSON(monthlyUsage []types.MonthlyUsage) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(monthlyUsage)
}

func outputMonthlyTable(monthlyUsage []types.MonthlyUsage) error {
	utils.FormatMonthlyUsageTable(monthlyUsage)
	return nil
}

func init() {
	rootCmd.AddCommand(monthlyCmd)
	
	// Monthly-specific flags
	monthlyCmd.Flags().String("start-month", "", "Start month (YYYY-MM)")
	monthlyCmd.Flags().String("end-month", "", "End month (YYYY-MM)")
	monthlyCmd.Flags().StringSlice("models", []string{}, "Filter by specific models")
}