package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/johanneserhardt/cxusage/internal/utils"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Show demo of beautiful output formatting",
	Long:  "Display sample usage data to demonstrate beautiful table formatting and colors.",
	Run:   runDemo,
}

func runDemo(cmd *cobra.Command, args []string) {
	// Create sample daily usage data
	sampleDaily := []types.DailyUsage{
		{
			Date:         "2025-08-31",
			RequestCount: 45,
			TotalTokens:  12450,
			TotalCost:    0.1245,
			ModelUsage: map[string]types.Usage{
				"gpt-4o": {
					PromptTokens:     8900,
					CompletionTokens: 2100,
					TotalTokens:      11000,
				},
				"gpt-4": {
					PromptTokens:     1200,
					CompletionTokens: 250,
					TotalTokens:      1450,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4o": 0.0995,
				"gpt-4":  0.0250,
			},
		},
		{
			Date:         "2025-08-30",
			RequestCount: 23,
			TotalTokens:  8750,
			TotalCost:    0.0875,
			ModelUsage: map[string]types.Usage{
				"gpt-4o": {
					PromptTokens:     6200,
					CompletionTokens: 1800,
					TotalTokens:      8000,
				},
				"gpt-3.5-turbo": {
					PromptTokens:     600,
					CompletionTokens: 150,
					TotalTokens:      750,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4o":        0.0800,
				"gpt-3.5-turbo": 0.0075,
			},
		},
		{
			Date:         "2025-08-29",
			RequestCount: 67,
			TotalTokens:  18900,
			TotalCost:    0.2890,
			ModelUsage: map[string]types.Usage{
				"gpt-4": {
					PromptTokens:     12000,
					CompletionTokens: 3500,
					TotalTokens:      15500,
				},
				"gpt-4o-mini": {
					PromptTokens:     2800,
					CompletionTokens: 600,
					TotalTokens:      3400,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4":       0.2650,
				"gpt-4o-mini": 0.0240,
			},
		},
		{
			Date:         "2025-08-28",
			RequestCount: 12,
			TotalTokens:  3200,
			TotalCost:    0.0320,
			ModelUsage: map[string]types.Usage{
				"gpt-3.5-turbo": {
					PromptTokens:     2400,
					CompletionTokens: 800,
					TotalTokens:      3200,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-3.5-turbo": 0.0320,
			},
		},
		{
			Date:         "2025-08-27",
			RequestCount: 89,
			TotalTokens:  25600,
			TotalCost:    1.2800,
			ModelUsage: map[string]types.Usage{
				"gpt-4": {
					PromptTokens:     18000,
					CompletionTokens: 5200,
					TotalTokens:      23200,
				},
				"gpt-4o": {
					PromptTokens:     2000,
					CompletionTokens: 400,
					TotalTokens:      2400,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4":  1.1600,
				"gpt-4o": 0.1200,
			},
		},
	}
	
	// Create sample monthly data
	sampleMonthly := []types.MonthlyUsage{
		{
			Month:        "2025-08",
			RequestCount: 1245,
			TotalTokens:  342000,
			TotalCost:    15.67,
			DailyBreakdown: sampleDaily,
			ModelUsage: map[string]types.Usage{
				"gpt-4": {
					PromptTokens:     240000,
					CompletionTokens: 60000,
					TotalTokens:      300000,
				},
				"gpt-4o": {
					PromptTokens:     32000,
					CompletionTokens: 10000,
					TotalTokens:      42000,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4":  12.50,
				"gpt-4o": 3.17,
			},
		},
		{
			Month:        "2025-07",
			RequestCount: 892,
			TotalTokens:  198000,
			TotalCost:    8.43,
			DailyBreakdown: []types.DailyUsage{}, // Simplified
			ModelUsage: map[string]types.Usage{
				"gpt-3.5-turbo": {
					PromptTokens:     150000,
					CompletionTokens: 48000,
					TotalTokens:      198000,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-3.5-turbo": 8.43,
			},
		},
		{
			Month:        "2025-06",
			RequestCount: 567,
			TotalTokens:  89000,
			TotalCost:    3.21,
			DailyBreakdown: []types.DailyUsage{},
			ModelUsage: map[string]types.Usage{
				"gpt-4o-mini": {
					PromptTokens:     70000,
					CompletionTokens: 19000,
					TotalTokens:      89000,
				},
			},
			ModelCosts: map[string]float64{
				"gpt-4o-mini": 3.21,
			},
		},
	}
	
	fmt.Println(utils.BoldWhite("🎨 cxusage - ccusage-style Output Demo"))
	
	utils.FormatDailyUsageTableProper(sampleDaily)
	
	fmt.Println()
	utils.FormatMonthlyUsageTableProper(sampleMonthly)
	
	fmt.Println()
	showBlocksDemo()
}

func showBlocksDemo() {
	fmt.Println(utils.BoldCyan("📊 5-Hour Blocks Live Monitoring"))
	fmt.Println()
	fmt.Printf("Use %s to see real-time usage tracking!\n", utils.BoldWhite("cxusage blocks --live"))
	fmt.Println()
	fmt.Printf("Features:\n")
	fmt.Printf("• %s - Real-time 5-hour billing block tracking\n", utils.Green("Live Dashboard"))
	fmt.Printf("• %s - Current usage vs projected usage\n", utils.Green("Projections"))
	fmt.Printf("• %s - Visual progress bar and burn rate\n", utils.Green("Progress Tracking"))
	fmt.Printf("• %s - Token usage and cost monitoring\n", utils.Green("Usage Alerts"))
	fmt.Println()
	fmt.Printf("Other blocks commands:\n")
	fmt.Printf("• %s - Show recent 5-hour blocks\n", utils.Gray("cxusage blocks"))
	fmt.Printf("• %s - Show only active block\n", utils.Gray("cxusage blocks --active"))
	fmt.Printf("• %s - Show blocks from last 3 days\n", utils.Gray("cxusage blocks --recent"))
}

func init() {
	rootCmd.AddCommand(demoCmd)
}