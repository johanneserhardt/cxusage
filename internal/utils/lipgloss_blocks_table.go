package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/johanneserhardt/cxusage/internal/blocks"
	"github.com/johanneserhardt/cxusage/internal/types"
)

// FormatBlocksTableProper creates a proper blocks table like ccusage
func FormatBlocksTableProper(sessionBlocks []types.SessionBlock, tokenLimit *int) {
	if len(sessionBlocks) == 0 {
		fmt.Println("No blocks found")
		return
	}
	
	// Print title with border
	title := "Codex CLI Usage Blocks (~estimated, 5-hour periods)"
	titleBorder := lipgloss.NewStyle().
		BorderStyle(tableBorderStyle).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true)
	
	fmt.Println()
	fmt.Println(titleBorder.Render(title))
	fmt.Println()
	
    // Define headers and build rows
    headers := []string{"Block Start", "Status", "Duration", "Requests", "Input", "Output", "Total Tokens", "Cost (USD)", "Models"}
	
	var rows [][]string
	var totalCost float64
	var totalRequests, totalTokens int
	activeBlockFound := false
	
	for _, block := range sessionBlocks {
		// Format status
		var status string
		if block.IsActive {
			status = "● Active"
			activeBlockFound = true
		} else if block.IsGap {
			status = "○ Gap"
		} else {
			status = "Complete"
		}
		
		// Format duration
		duration := formatBlockDurationSimple(block)
		
		// Format models
		modelsStr := formatModelsListSimple(block.Models)
		
		// Format tokens with warnings
		tokensStr := FormatNumber(block.TotalTokens)
		if tokenLimit != nil && block.TotalTokens > *tokenLimit {
			tokensStr += " ⚠"
		}
		
		// Create row
		row := []string{
			block.StartTime.Format("2006-01-02 15:04"),
			status,
			duration,
			FormatNumber(block.RequestCount),
			FormatNumber(block.InputTokens),
			FormatNumber(block.OutputTokens),
			tokensStr,
			FormatCurrency(block.TotalCost),
			modelsStr,
		}
		rows = append(rows, row)
		
		if !block.IsGap {
			totalCost += block.TotalCost
			totalRequests += block.RequestCount
			totalTokens += block.TotalTokens
		}
	}
	
	// Add totals row
	totalRow := []string{
		"TOTAL",
		"",
		"",
		FormatNumber(totalRequests),
		"",
		"",
		FormatNumber(totalTokens),
		FormatCurrency(totalCost),
		"",
	}
	rows = append(rows, totalRow)
	
    // Autosize widths
    min := []int{14, 8, 8, 7, 7, 7, 10, 10, 10}
    if isCompact() {
        min = []int{12, 6, 7, 6, 6, 6, 9, 9, 8}
    }
    widths := computeAutoWidths(headers, rows, min)
    // Render the table
    table := CreateTable(headers, rows, widths)
    fmt.Println(table)
	
	// Show active block projection if found
	if activeBlockFound {
		if activeBlock := blocks.GetActiveBlock(sessionBlocks); activeBlock != nil {
			fmt.Println()
			showActiveBlockProjectionProper(activeBlock)
		}
	}
}

// formatBlockDurationSimple formats block duration for table display
func formatBlockDurationSimple(block types.SessionBlock) string {
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
        
        return fmt.Sprintf("%dh%dm/%dh%dm", elapsedHours, elapsedMins, remainingHours, remainingMins)
    }
	
	return "5h"
}

// formatModelsListSimple formats models for table display
func formatModelsListSimple(models []string) string {
	if len(models) == 0 {
		return "-"
	}
	
	modelsStr := strings.Join(models, ", ")
	if len(modelsStr) > 18 {
		modelsStr = modelsStr[:15] + "..."
	}
	return modelsStr
}

// showActiveBlockProjectionProper shows projections with proper table formatting
func showActiveBlockProjectionProper(block *types.SessionBlock) {
	projection := blocks.CalculateProjections(block)
	if projection == nil {
		return
	}
	
	// Create projection title box
	title := "🔮 Active Block Projections"
	titleBorder := lipgloss.NewStyle().
		BorderStyle(tableBorderStyle).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true)
	
	fmt.Println(titleBorder.Render(title))
	fmt.Println()
	
	fmt.Printf("Current: %s tokens, %s\n", 
		FormatNumber(block.TotalTokens),
		FormatCurrency(block.TotalCost))
		
	fmt.Printf("Projected 5h total: %s tokens, %s\n",
		FormatNumber(projection.ProjectedTokens),
		FormatCurrency(projection.ProjectedCost))
		
	fmt.Printf("Burn rate: %.1f tokens/min\n", projection.BurnRate)
	
	// Show warning if projection is high
	if projection.ProjectedCost > 5.0 {
		fmt.Printf("⚠️ High cost projected for this block!\n")
	}
}
