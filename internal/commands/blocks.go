package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/johanneserhardt/cxusage/internal/blocks"
	"github.com/johanneserhardt/cxusage/internal/codex"
	"github.com/johanneserhardt/cxusage/internal/live"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/johanneserhardt/cxusage/internal/utils"
)

var blocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "Show usage in 5-hour billing blocks (supports --live monitoring!)",
	Long: `Display usage data grouped by 5-hour billing blocks (similar to Claude's billing periods).

🔥 LIVE MONITORING: Use --live flag for real-time dashboard with gorgeous visuals!

Shows current active block, recent blocks, and supports live monitoring with:
• Real-time progress bars and burn rate tracking
• Visual alerts for token limit warnings  
• Beautiful full-screen dashboard updates
• Projections and usage forecasting`,
	RunE: runBlocks,
}

func runBlocks(cmd *cobra.Command, args []string) error {
	// Get flags
	outputFormat, _ := cmd.Flags().GetString("output")
	liveMode, _ := cmd.Flags().GetBool("live")
	activeOnly, _ := cmd.Flags().GetBool("active")
	recentOnly, _ := cmd.Flags().GetBool("recent")
	recentDays, _ := cmd.Flags().GetInt("recent-days")
	sessionHours, _ := cmd.Flags().GetInt("session-duration")
	refreshInterval, _ := cmd.Flags().GetInt("refresh-interval")
	tokenLimitStr, _ := cmd.Flags().GetString("token-limit")
	
	// Validate session duration
	if sessionHours < 1 || sessionHours > 24 {
		return fmt.Errorf("session duration must be between 1 and 24 hours")
	}
	
	// Parse token limit
	var tokenLimit *int
	if tokenLimitStr != "" && tokenLimitStr != "max" {
		if limit, err := strconv.Atoi(tokenLimitStr); err == nil && limit > 0 {
			tokenLimit = &limit
		}
	}
	
	// Handle live mode
	if liveMode {
		refreshDuration := time.Duration(refreshInterval) * time.Second
		if refreshDuration < live.MinRefreshInterval {
			refreshDuration = live.MinRefreshInterval
		}
		if refreshDuration > live.MaxRefreshInterval {
			refreshDuration = live.MaxRefreshInterval
		}
		
		config := &types.LiveMonitoringConfig{
			RefreshInterval:      refreshDuration,
			SessionDurationHours: sessionHours,
			TokenLimit:           tokenLimit,
			ShowProjections:      true,
		}
		
		monitor := live.NewLiveMonitor(config, cfg, logger)
		return monitor.Start()
	}
	
	// Calculate date range
	var startDate, endDate time.Time
	if recentOnly {
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -recentDays)
	} else if activeOnly {
		// For active only, look at last 24 hours
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -1)
	} else {
		// Default: last 7 days
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -7)
	}
	
	logger.WithFields(map[string]interface{}{
		"start_date":     startDate.Format("2006-01-02"),
		"end_date":       endDate.Format("2006-01-02"),
		"session_hours":  sessionHours,
		"active_only":    activeOnly,
		"recent_only":    recentOnly,
	}).Info("Generating blocks usage report")
	
	// Load usage data
	entries, err := codex.ParseUsageFiles(cfg, startDate, endDate, logger)
	if err != nil {
		return fmt.Errorf("failed to load usage data: %w", err)
	}
	
	if len(entries) == 0 {
		if outputFormat == "json" {
			fmt.Println("[]")
		} else {
			fmt.Printf("%s\n", utils.Yellow("No Codex CLI usage data found"))
			fmt.Println()
			fmt.Printf("This could mean:\n")
			fmt.Printf("• Codex CLI hasn't been used recently\n")
			fmt.Printf("• Codex CLI is not installed or configured\n")
			fmt.Printf("• Usage logs are stored in a different location\n")
			fmt.Println()
			fmt.Printf("Try:\n")
			fmt.Printf("• %s - Check if Codex CLI is set up\n", utils.Cyan("cxusage validate"))
			fmt.Printf("• Use Codex CLI first, then run %s\n", utils.Cyan("cxusage blocks"))
			fmt.Printf("• Run %s to see sample output\n", utils.Cyan("cxusage demo"))
		}
		return nil
	}
	
	// Aggregate into blocks
	sessionBlocks := blocks.AggregateIntoBlocks(entries, sessionHours)
	
	// Apply filters
	if recentOnly {
		sessionBlocks = blocks.FilterRecentBlocks(sessionBlocks, recentDays)
	}
	
	if activeOnly {
		if activeBlock := blocks.GetActiveBlock(sessionBlocks); activeBlock != nil {
			sessionBlocks = []types.SessionBlock{*activeBlock}
		} else {
			sessionBlocks = []types.SessionBlock{}
		}
	}
	
	// Output results
	switch types.OutputFormat(outputFormat) {
	case types.OutputFormatJSON:
		return outputBlocksJSON(sessionBlocks)
	case types.OutputFormatTable:
		return outputBlocksTable(sessionBlocks, tokenLimit)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
}

func outputBlocksJSON(sessionBlocks []types.SessionBlock) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(sessionBlocks)
}

func outputBlocksTable(sessionBlocks []types.SessionBlock, tokenLimit *int) error {
	if len(sessionBlocks) == 0 {
		fmt.Println(utils.Yellow("No blocks found"))
		return nil
	}
	
	utils.PrintTitleBox("Codex CLI Usage Blocks (5-hour periods)")
	
	table := utils.CreateCcusageStyleTable([]string{
		"Block Start",
		"Status",
		"Duration",
		"Requests",
		"Input",
		"Output",
		"Total Tokens",
		"Cost (USD)",
		"Models",
	})
	
	var totalCost float64
	var totalRequests, totalTokens int
	activeBlockFound := false
	
	for _, block := range sessionBlocks {
		status := "Complete"
		duration := formatBlockDuration(block)
		
		if block.IsActive {
			status = utils.Green("● Active")
			activeBlockFound = true
		} else if block.IsGap {
			status = utils.Gray("○ Gap")
		}
		
		// Format models list
		modelsStr := formatModelsList(block.Models)
		
		// Apply token limit warning
		tokensStr := utils.FormatNumber(block.TotalTokens)
		if tokenLimit != nil && block.TotalTokens > *tokenLimit {
			tokensStr = utils.Red(tokensStr + " ⚠")
		}
		
		// Color code cost
		costStr := utils.FormatCurrency(block.TotalCost)
		if block.TotalCost > 1.0 {
			costStr = utils.Red(costStr)
		} else if block.TotalCost > 0.1 {
			costStr = utils.Yellow(costStr)
		} else if block.TotalCost > 0 {
			costStr = utils.Green(costStr)
		}
		
		table.Append([]string{
			block.StartTime.Format("2006-01-02 15:04"),
			status,
			duration,
			utils.FormatNumber(block.RequestCount),
			utils.Gray(utils.FormatNumber(block.InputTokens)),
			utils.Gray(utils.FormatNumber(block.OutputTokens)),
			tokensStr,
			costStr,
			modelsStr,
		})
		
		if !block.IsGap {
			totalCost += block.TotalCost
			totalRequests += block.RequestCount
			totalTokens += block.TotalTokens
		}
	}
	
	// Add totals row
	table.SetFooter([]string{
		utils.BoldYellow("TOTAL"),
		"",
		"",
		utils.BoldYellow(utils.FormatNumber(totalRequests)),
		"",
		"",
		utils.BoldYellow(utils.FormatNumber(totalTokens)),
		utils.BoldYellow(utils.FormatCurrency(totalCost)),
		"",
	})
	
	table.Render()
	
	// Show active block projection if found
	if activeBlockFound {
		if activeBlock := blocks.GetActiveBlock(sessionBlocks); activeBlock != nil {
			fmt.Println()
			showActiveBlockProjection(activeBlock)
		}
	}
	
	return nil
}

// formatBlockDuration formats the duration of a block
func formatBlockDuration(block types.SessionBlock) string {
	if block.IsGap {
		duration := block.EndTime.Sub(block.StartTime)
		hours := int(duration.Hours())
		return utils.Gray(fmt.Sprintf("%dh gap", hours))
	}
	
	if block.ActualEndTime != nil {
		duration := block.ActualEndTime.Sub(block.StartTime)
		hours := int(duration.Hours())
		mins := int(duration.Minutes()) % 60
		if hours > 0 {
			return fmt.Sprintf("%dh %dm", hours, mins)
		}
		return fmt.Sprintf("%dm", mins)
	}
	
	if block.IsActive {
		now := time.Now()
		elapsed := now.Sub(block.StartTime)
		remaining := block.EndTime.Sub(now)
		elapsedHours := int(elapsed.Hours())
		elapsedMins := int(elapsed.Minutes()) % 60
		remainingHours := int(remaining.Hours())
		remainingMins := int(remaining.Minutes()) % 60
		
		return utils.Green(fmt.Sprintf("%dh %dm / %dh %dm", elapsedHours, elapsedMins, remainingHours, remainingMins))
	}
	
	return "5h"
}

// formatModelsList formats the list of models for display
func formatModelsList(models []string) string {
	if len(models) == 0 {
		return utils.Gray("-")
	}
	
	if len(models) == 1 {
		return models[0]
	}
	
	// Show first 2 models, then "..."
	if len(models) > 2 {
		return fmt.Sprintf("%s, %s...", models[0], models[1])
	}
	
	return strings.Join(models, ", ")
}

// showActiveBlockProjection shows projections for the active block
func showActiveBlockProjection(block *types.SessionBlock) {
	projection := blocks.CalculateProjections(block)
	if projection == nil {
		return
	}
	
	fmt.Printf("%s\n", utils.BoldWhite("🔮 Active Block Projections"))
	fmt.Printf("Current: %s tokens, %s\n", 
		utils.FormatNumber(block.TotalTokens),
		utils.FormatCurrency(block.TotalCost))
	fmt.Printf("Projected 5h total: %s tokens, %s\n",
		utils.Cyan(utils.FormatNumber(projection.ProjectedTokens)),
		utils.Cyan(utils.FormatCurrency(projection.ProjectedCost)))
	fmt.Printf("Burn rate: %s tokens/min\n", 
		utils.Gray(fmt.Sprintf("%.1f", projection.BurnRate)))
	
	// Show warning if projection is high
	if projection.ProjectedCost > 5.0 {
		fmt.Printf("%s High cost projected for this block!\n", utils.Red("⚠️"))
	}
}

func init() {
	rootCmd.AddCommand(blocksCmd)
	
	// Blocks-specific flags
	blocksCmd.Flags().Bool("live", false, "Enable live monitoring mode")
	blocksCmd.Flags().Bool("active", false, "Show only the currently active block")
	blocksCmd.Flags().Bool("recent", false, "Show only recent blocks")
	blocksCmd.Flags().Int("recent-days", 3, "Number of recent days to show (with --recent)")
	blocksCmd.Flags().Int("session-duration", 5, "Block duration in hours (default: 5)")
	blocksCmd.Flags().Int("refresh-interval", 1, "Refresh interval in seconds for live mode")
	blocksCmd.Flags().String("token-limit", "", "Token limit threshold for warnings (number or 'max')")
}