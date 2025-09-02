package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/johanneserhardt/cxusage/internal/types"
)

// Table styles for proper table formatting like ccusage
var (
	// Table border styles
	tableBorderStyle = lipgloss.Border{
		Top:    "─",
		Bottom: "─", 
		Left:   "│",
		Right:  "│",
		TopLeft: "┌",
		TopRight: "┐",
		BottomLeft: "└",
		BottomRight: "┘",
	}
	
	// Cell styles with proper table formatting
	tableCellStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Left)
	
	tableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Align(lipgloss.Left)
	
	tableNumberStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Right)
	
	tableTotalStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor).
		Padding(0, 1).
		Align(lipgloss.Left)
)

// CreateTable creates a proper table structure like ccusage
func CreateTable(headers []string, rows [][]string, widths []int) string {
	var result strings.Builder
	
	// Create top border
	result.WriteString("┌")
	for i, width := range widths {
		result.WriteString(strings.Repeat("─", width+2)) // +2 for padding
		if i < len(widths)-1 {
			result.WriteString("┬")
		}
	}
	result.WriteString("┐\n")
	
	// Create header row
	result.WriteString("│")
	for i, header := range headers {
		content := padString(header, widths[i])
		result.WriteString(fmt.Sprintf(" %s ", content))
		result.WriteString("│")
	}
	result.WriteString("\n")
	
	// Create header separator
	result.WriteString("├")
	for i, width := range widths {
		result.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			result.WriteString("┼")
		}
	}
	result.WriteString("┤\n")
	
	// Create data rows
	for _, row := range rows {
		result.WriteString("│")
		for i, cell := range row {
			if i < len(widths) {
				content := padString(cell, widths[i])
				result.WriteString(fmt.Sprintf(" %s ", content))
				result.WriteString("│")
			}
		}
		result.WriteString("\n")
	}
	
	// Create bottom border
	result.WriteString("└")
	for i, width := range widths {
		result.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			result.WriteString("┴")
		}
	}
	result.WriteString("┘")
	
	return result.String()
}

// FormatDailyUsageTableProper creates a proper table like ccusage
func FormatDailyUsageTableProper(dailyUsage []types.DailyUsage) {
	if len(dailyUsage) == 0 {
		fmt.Println("No usage data found")
		return
	}
	
	// Print title with border like ccusage
	title := "Codex CLI Token Usage Report - Daily"
	titleBorder := lipgloss.NewStyle().
		BorderStyle(tableBorderStyle).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true)
	
	fmt.Println()
	fmt.Println(titleBorder.Render(title))
	fmt.Println()
	
	// Define headers and calculate data
	headers := []string{"Date", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total Tokens", "Cost (USD)"}
	widths := []int{12, 20, 10, 10, 12, 12, 12, 10}
	
	var rows [][]string
	var totalCost float64
	var totalInput, totalOutput, totalTokens int
	
	// Process each day
	for _, day := range dailyUsage {
		var inputTokens, outputTokens int
		var modelsList []string
		
		for model, usage := range day.ModelUsage {
			inputTokens += usage.PromptTokens
			outputTokens += usage.CompletionTokens
			modelsList = append(modelsList, model)
		}
		
		modelsStr := strings.Join(modelsList, ", ")
		if len(modelsStr) > 18 {
			modelsStr = modelsStr[:15] + "..."
		}
		modelsStr = "- " + modelsStr
		
		// Create row
		row := []string{
			day.Date,
			modelsStr,
			FormatNumber(inputTokens),
			FormatNumber(outputTokens),
			"0", // Cache create
			"0", // Cache read
			FormatNumber(day.TotalTokens),
			FormatCurrency(day.TotalCost),
		}
		rows = append(rows, row)
		
		totalCost += day.TotalCost
		totalInput += inputTokens
		totalOutput += outputTokens
		totalTokens += day.TotalTokens
	}
	
	// Add totals row
	totalRow := []string{
		"Total",
		"",
		FormatNumber(totalInput),
		FormatNumber(totalOutput),
		"0",
		"0", 
		FormatNumber(totalTokens),
		FormatCurrency(totalCost),
	}
	rows = append(rows, totalRow)
	
	// Render the table
	table := CreateTable(headers, rows, widths)
	fmt.Println(table)
}

// FormatMonthlyUsageTableProper creates a proper monthly table like ccusage
func FormatMonthlyUsageTableProper(monthlyUsage []types.MonthlyUsage) {
	if len(monthlyUsage) == 0 {
		fmt.Println("No usage data found")
		return
	}
	
	// Print title with border
	title := "Codex CLI Token Usage Report - Monthly"
	titleBorder := lipgloss.NewStyle().
		BorderStyle(tableBorderStyle).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Foreground(primaryColor).
		Bold(true)
	
	fmt.Println()
	fmt.Println(titleBorder.Render(title))
	fmt.Println()
	
	// Define headers
	headers := []string{"Month", "Days Active", "Total Requests", "Input Tokens", "Output Tokens", "Total Tokens", "Total Cost (USD)"}
	widths := []int{10, 12, 15, 14, 15, 14, 16}
	
	var rows [][]string
	var totalCost float64
	var totalRequests, totalInput, totalOutput, totalTokens int
	
	// Process each month
	for _, month := range monthlyUsage {
		var inputTokens, outputTokens int
		for _, usage := range month.ModelUsage {
			inputTokens += usage.PromptTokens
			outputTokens += usage.CompletionTokens
		}
		
		activeDays := len(month.DailyBreakdown)
		if activeDays == 0 {
			activeDays = 1
		}
		
		row := []string{
			month.Month,
			strconv.Itoa(activeDays),
			FormatNumber(month.RequestCount),
			FormatNumber(inputTokens),
			FormatNumber(outputTokens),
			FormatNumber(month.TotalTokens),
			FormatCurrency(month.TotalCost),
		}
		rows = append(rows, row)
		
		totalCost += month.TotalCost
		totalRequests += month.RequestCount
		totalInput += inputTokens
		totalOutput += outputTokens
		totalTokens += month.TotalTokens
	}
	
	// Add totals row
	totalRow := []string{
		"Total",
		"",
		FormatNumber(totalRequests),
		FormatNumber(totalInput),
		FormatNumber(totalOutput),
		FormatNumber(totalTokens),
		FormatCurrency(totalCost),
	}
	rows = append(rows, totalRow)
	
	// Render the table
	table := CreateTable(headers, rows, widths)
	fmt.Println(table)
}

// padString pads a string to a specific width
func padString(s string, width int) string {
	if len(s) >= width {
		if width > 3 {
			return s[:width-3] + "..."
		}
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}