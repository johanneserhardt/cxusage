package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/johanneserhardt/cxusage/internal/blocks"
	"github.com/johanneserhardt/cxusage/internal/types"
)

// FormatBlocksTableLipgloss creates a beautiful blocks table with Lipgloss
func FormatBlocksTableLipgloss(sessionBlocks []types.SessionBlock, tokenLimit *int) {
	if len(sessionBlocks) == 0 {
		warningStyle := lipgloss.NewStyle().Foreground(warningColor)
		fmt.Println(warningStyle.Render("No blocks found"))
		return
	}
	
	// Create title
	title := titleStyle.Render("Codex CLI Usage Blocks (5-hour periods)")
	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	
	// Define column widths for blocks table
	colWidths := []int{16, 12, 10, 10, 10, 10, 14, 12, 20}
	headers := []string{"Block Start", "Status", "Duration", "Requests", "Input", "Output", "Total Tokens", "Cost (USD)", "Models"}
	
	// Create header row
	var headerCells []string
	for i, header := range headers {
		style := headerStyle.Width(colWidths[i])
		headerCells = append(headerCells, style.Render(header))
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	
	// Create data rows
	var totalCost float64
	var totalRequests, totalTokens int
	activeBlockFound := false
	
	for _, block := range sessionBlocks {
		// Format status
		var status string
		var statusStyle lipgloss.Style
		
		if block.IsActive {
			status = "● Active"
			statusStyle = lipgloss.NewStyle().Foreground(successColor).Width(colWidths[1])
			activeBlockFound = true
		} else if block.IsGap {
			status = "○ Gap"
			statusStyle = lipgloss.NewStyle().Foreground(mutedColor).Width(colWidths[1])
		} else {
			status = "Complete"
			statusStyle = cellStyle.Width(colWidths[1])
		}
		
		// Format duration
		duration := formatBlockDuration(block)
		
		// Format models list
		modelsStr := formatModelsListForBlocks(block.Models)
		
		// Apply token limit warning
		tokensStr := FormatNumber(block.TotalTokens)
		tokensStyle := numberStyle.Width(colWidths[6])
		if tokenLimit != nil && block.TotalTokens > *tokenLimit {
			tokensStr += " ⚠"
			tokensStyle = tokensStyle.Foreground(errorColor)
		}
		
		// Style cells
		cells := []string{
			dateStyle.Width(colWidths[0]).Render(block.StartTime.Format("2006-01-02 15:04")),
			statusStyle.Render(status),
			cellStyle.Width(colWidths[2]).Render(duration),
			numberStyle.Width(colWidths[3]).Render(FormatNumber(block.RequestCount)),
			numberStyle.Width(colWidths[4]).Render(FormatNumber(block.InputTokens)),
			numberStyle.Width(colWidths[5]).Render(FormatNumber(block.OutputTokens)),
			tokensStyle.Render(tokensStr),
			formatCostCell(block.TotalCost, colWidths[7]),
			modelStyle.Width(colWidths[8]).Render(modelsStr),
		}
		
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		
		if !block.IsGap {
			totalCost += block.TotalCost
			totalRequests += block.RequestCount
			totalTokens += block.TotalTokens
		}
	}
	
	// Create separator
	separator := strings.Repeat("─", sum(colWidths)+len(colWidths)-1)
	separatorStyle := lipgloss.NewStyle().Foreground(mutedColor)
	fmt.Println(separatorStyle.Render(separator))
	
	// Create totals row
	totalCells := []string{
		totalStyle.Width(colWidths[0]).Render("TOTAL"),
		cellStyle.Width(colWidths[1]).Render(""),
		cellStyle.Width(colWidths[2]).Render(""),
		totalStyle.Width(colWidths[3]).Render(FormatNumber(totalRequests)),
		cellStyle.Width(colWidths[4]).Render(""),
		cellStyle.Width(colWidths[5]).Render(""),
		totalStyle.Width(colWidths[6]).Render(FormatNumber(totalTokens)),
		totalStyle.Width(colWidths[7]).Render(FormatCurrency(totalCost)),
		cellStyle.Width(colWidths[8]).Render(""),
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, totalCells...))
	
	// Show active block projection if found
	if activeBlockFound {
		if activeBlock := blocks.GetActiveBlock(sessionBlocks); activeBlock != nil {
			fmt.Println()
			showActiveBlockProjectionLipgloss(activeBlock)
		}
	}
	
	fmt.Println()
}

// formatBlockDuration formats the duration of a block with proper styling
func formatBlockDuration(block types.SessionBlock) string {
    if block.IsGap {
        duration := block.EndTime.Sub(block.StartTime)
        hours := int(duration.Hours())
        return fmt.Sprintf("%dh gap", hours)
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
        if remaining < 0 {
            remaining = 0
        }
        elapsedHours := int(elapsed.Hours())
        elapsedMins := int(elapsed.Minutes()) % 60
        remainingHours := int(remaining.Hours())
        remainingMins := int(remaining.Minutes()) % 60

        return fmt.Sprintf("%dh %dm / %dh %dm", elapsedHours, elapsedMins, remainingHours, remainingMins)
    }
	
	return "5h"
}

// formatModelsListForBlocks formats the list of models for blocks display
func formatModelsListForBlocks(models []string) string {
	if len(models) == 0 {
		return "-"
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

// showActiveBlockProjectionLipgloss shows projections for the active block with beautiful styling
func showActiveBlockProjectionLipgloss(block *types.SessionBlock) {
	projection := blocks.CalculateProjections(block)
	if projection == nil {
		return
	}
	
	// Create projection title
	projectionTitle := titleStyle.Render("🔮 Active Block Projections")
	fmt.Println(projectionTitle)
	fmt.Println()
	
	// Current vs projected with beautiful styling
	currentStyle := lipgloss.NewStyle().Foreground(textColor).Bold(true)
	projectedStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	
	fmt.Printf("Current: %s tokens, %s\n", 
		currentStyle.Render(FormatNumber(block.TotalTokens)),
		currentStyle.Render(FormatCurrency(block.TotalCost)))
	
	fmt.Printf("Projected 5h total: %s tokens, %s\n",
		projectedStyle.Render(FormatNumber(projection.ProjectedTokens)),
		projectedStyle.Render(FormatCurrency(projection.ProjectedCost)))
	
	burnRateStyle := lipgloss.NewStyle().Foreground(mutedColor)
	fmt.Printf("Burn rate: %s\n", 
		burnRateStyle.Render(fmt.Sprintf("%.1f tokens/min", projection.BurnRate)))
	
	// Show warning if projection is high
	if projection.ProjectedCost > 5.0 {
		warningStyle := lipgloss.NewStyle().Foreground(errorColor).Bold(true)
		fmt.Printf("%s\n", warningStyle.Render("⚠️ High cost projected for this block!"))
	}
}
