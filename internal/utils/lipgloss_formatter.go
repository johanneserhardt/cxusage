package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/johanneserhardt/cxusage/internal/types"
)

// Theme-adaptive colors that work across different terminal themes
var (
	// Use adaptive colors that work with any terminal theme
	primaryColor   = lipgloss.AdaptiveColor{Light: "#005577", Dark: "#5FAFFF"} // Blue tones
	successColor   = lipgloss.AdaptiveColor{Light: "#00AA00", Dark: "#55FF55"} // Green tones  
	warningColor   = lipgloss.AdaptiveColor{Light: "#DD6600", Dark: "#FFAA00"} // Orange/Yellow tones
	errorColor     = lipgloss.AdaptiveColor{Light: "#CC0000", Dark: "#FF5555"} // Red tones
	mutedColor     = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"} // Gray tones
	textColor      = lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"} // Text
	
	// Styles for different elements
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1)
	
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		Background(primaryColor).
		Padding(0, 1).
		Align(lipgloss.Center)
	
	cellStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Left)
	
	numberStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Right).
		Padding(0, 1)
	
	costStyle = lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Right).
		Padding(0, 1)
	
	totalStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor).
		Padding(0, 1)
	
	modelStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Padding(0, 1)
	
	dateStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Padding(0, 1)
)

// FormatDailyUsageTableLipgloss creates a beautiful daily usage table with Lipgloss
func FormatDailyUsageTableLipgloss(dailyUsage []types.DailyUsage) {
	if len(dailyUsage) == 0 {
		fmt.Println(warningStyle.Render("No usage data found"))
		return
	}
	
	// Create title
	title := titleStyle.Render("Codex CLI Token Usage Report - Daily")
	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	
	// Define column widths
	colWidths := []int{12, 25, 10, 10, 12, 12, 14, 12}
	headers := []string{"Date", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total Tokens", "Cost (USD)"}
	
	// Create header row
	var headerCells []string
	for i, header := range headers {
		style := headerStyle.Width(colWidths[i])
		headerCells = append(headerCells, style.Render(header))
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	
	// Create data rows
	var totalCost float64
	var totalInput, totalOutput, totalCacheCreate, totalCacheRead, totalTokens int
	
	for _, day := range dailyUsage {
		// Calculate token breakdowns
		var inputTokens, outputTokens int
		var modelsList []string
		
		for model, usage := range day.ModelUsage {
			inputTokens += usage.PromptTokens
			outputTokens += usage.CompletionTokens
			modelsList = append(modelsList, model)
		}
		
		// Format models list
		modelsStr := strings.Join(modelsList, ", ")
		if len(modelsStr) > 22 {
			modelsStr = modelsStr[:19] + "..."
		}
		modelsStr = "- " + modelsStr
		
		// Cache tokens (0 for now)
		cacheCreate := 0
		cacheRead := 0
		
		// Style each cell
		cells := []string{
			dateStyle.Width(colWidths[0]).Render(day.Date),
			modelStyle.Width(colWidths[1]).Render(modelsStr),
			numberStyle.Width(colWidths[2]).Render(FormatNumber(inputTokens)),
			numberStyle.Width(colWidths[3]).Render(FormatNumber(outputTokens)),
			numberStyle.Width(colWidths[4]).Render(FormatNumber(cacheCreate)),
			numberStyle.Width(colWidths[5]).Render(FormatNumber(cacheRead)),
			numberStyle.Width(colWidths[6]).Render(FormatNumber(day.TotalTokens)),
			formatCostCell(day.TotalCost, colWidths[7]),
		}
		
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		
		// Update totals
		totalCost += day.TotalCost
		totalInput += inputTokens
		totalOutput += outputTokens
		totalCacheCreate += cacheCreate
		totalCacheRead += cacheRead
		totalTokens += day.TotalTokens
	}
	
	// Create separator
	separator := strings.Repeat("─", sum(colWidths)+len(colWidths)-1)
	fmt.Println(lipgloss.NewStyle().Foreground(mutedColor).Render(separator))
	
	// Create totals row
	totalCells := []string{
		totalStyle.Width(colWidths[0]).Render("TOTAL"),
		cellStyle.Width(colWidths[1]).Render(""),
		totalStyle.Width(colWidths[2]).Render(FormatNumber(totalInput)),
		totalStyle.Width(colWidths[3]).Render(FormatNumber(totalOutput)),
		totalStyle.Width(colWidths[4]).Render(FormatNumber(totalCacheCreate)),
		totalStyle.Width(colWidths[5]).Render(FormatNumber(totalCacheRead)),
		totalStyle.Width(colWidths[6]).Render(FormatNumber(totalTokens)),
		totalStyle.Width(colWidths[7]).Render(FormatCurrency(totalCost)),
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, totalCells...))
	fmt.Println()
}

// FormatMonthlyUsageTableLipgloss creates a beautiful monthly usage table with Lipgloss
func FormatMonthlyUsageTableLipgloss(monthlyUsage []types.MonthlyUsage) {
	if len(monthlyUsage) == 0 {
		fmt.Println(warningStyle.Render("No usage data found"))
		return
	}
	
	// Create title
	title := titleStyle.Render("Codex CLI Token Usage Report - Monthly")
	fmt.Println()
	fmt.Println(title)
	fmt.Println()
	
	// Define column widths for monthly table
	colWidths := []int{10, 12, 15, 14, 15, 14, 16}
	headers := []string{"Month", "Days Active", "Total Requests", "Input Tokens", "Output Tokens", "Total Tokens", "Total Cost (USD)"}
	
	// Create header row
	var headerCells []string
	for i, header := range headers {
		style := headerStyle.Width(colWidths[i])
		headerCells = append(headerCells, style.Render(header))
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...))
	
	// Create data rows
	var totalCost float64
	var totalRequests, totalInput, totalOutput, totalTokens int
	
	for _, month := range monthlyUsage {
		// Calculate token breakdowns
		var inputTokens, outputTokens int
		for _, usage := range month.ModelUsage {
			inputTokens += usage.PromptTokens
			outputTokens += usage.CompletionTokens
		}
		
		activeDays := len(month.DailyBreakdown)
		if activeDays == 0 {
			activeDays = 1
		}
		
		// Style each cell
		cells := []string{
			dateStyle.Width(colWidths[0]).Render(month.Month),
			numberStyle.Width(colWidths[1]).Render(strconv.Itoa(activeDays)),
			numberStyle.Width(colWidths[2]).Render(FormatNumber(month.RequestCount)),
			numberStyle.Width(colWidths[3]).Render(FormatNumber(inputTokens)),
			numberStyle.Width(colWidths[4]).Render(FormatNumber(outputTokens)),
			numberStyle.Width(colWidths[5]).Render(FormatNumber(month.TotalTokens)),
			formatCostCell(month.TotalCost, colWidths[6]),
		}
		
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
		
		// Update totals
		totalCost += month.TotalCost
		totalRequests += month.RequestCount
		totalInput += inputTokens
		totalOutput += outputTokens
		totalTokens += month.TotalTokens
	}
	
	// Create separator
	separator := strings.Repeat("─", sum(colWidths)+len(colWidths)-1)
	fmt.Println(lipgloss.NewStyle().Foreground(mutedColor).Render(separator))
	
	// Create totals row
	totalCells := []string{
		totalStyle.Width(colWidths[0]).Render("TOTAL"),
		cellStyle.Width(colWidths[1]).Render(""),
		totalStyle.Width(colWidths[2]).Render(FormatNumber(totalRequests)),
		totalStyle.Width(colWidths[3]).Render(FormatNumber(totalInput)),
		totalStyle.Width(colWidths[4]).Render(FormatNumber(totalOutput)),
		totalStyle.Width(colWidths[5]).Render(FormatNumber(totalTokens)),
		totalStyle.Width(colWidths[6]).Render(FormatCurrency(totalCost)),
	}
	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, totalCells...))
	fmt.Println()
}

// formatCostCell formats cost with color based on amount
func formatCostCell(cost float64, width int) string {
	costStr := FormatCurrency(cost)
	var style lipgloss.Style
	
	if cost > 1.0 {
		style = costStyle.Foreground(errorColor).Width(width)
	} else if cost > 0.1 {
		style = costStyle.Foreground(warningColor).Width(width)
	} else {
		style = costStyle.Foreground(successColor).Width(width)
	}
	
	return style.Render(costStr)
}

// sum returns the sum of integers in a slice
func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// warningStyle for warnings and errors
var warningStyle = lipgloss.NewStyle().
	Foreground(warningColor).
	Bold(true)