package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/johanneserhardt/cxusage/internal/types"
)

// Colors matching ccusage style
var (
	// Primary colors
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	
	// Muted colors  
	Gray        = color.New(color.FgHiBlack).SprintFunc()
	LightGray   = color.New(color.FgBlack).SprintFunc()
	
	// Bold colors
	BoldYellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	BoldGreen   = color.New(color.FgGreen, color.Bold).SprintFunc()
	BoldRed     = color.New(color.FgRed, color.Bold).SprintFunc()
	BoldCyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	BoldWhite   = color.New(color.FgWhite, color.Bold).SprintFunc()
)

// FormatNumber formats a number with thousand separators
func FormatNumber(n int) string {
	if n == 0 {
		return "0"
	}
	
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	
	// Add commas every 3 digits from right
	var result strings.Builder
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}
	
	return result.String()
}

// FormatCurrency formats a float64 as currency with $ sign
func FormatCurrency(amount float64) string {
    if amount == 0 {
        return "$0.00"
    }
    // Show finer precision for small amounts
    if amount < 0.10 {
        return fmt.Sprintf("$%.4f", amount)
    }
    return fmt.Sprintf("$%.2f", amount)
}

// CreateCcusageStyleTable creates a table matching ccusage's style
func CreateCcusageStyleTable(headers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	
	// ccusage style formatting - simple borders
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|") 
	table.SetRowSeparator("-")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderLine(true)
	table.SetTablePadding(" ")
	table.SetNoWhiteSpace(false)
	
	return table
}

// PrintTitleBox prints a title in a box like ccusage
func PrintTitleBox(title string) {
	titleLen := len(title)
	border := "┌" + strings.Repeat("─", titleLen+2) + "┐"
	content := "│ " + title + " │"
	bottomBorder := "└" + strings.Repeat("─", titleLen+2) + "┘"
	
	fmt.Println()
	fmt.Println(border)
	fmt.Println(content)
	fmt.Println(bottomBorder)
	fmt.Println()
}

// FormatDailyUsageTable creates a daily usage table matching ccusage style
func FormatDailyUsageTable(dailyUsage []types.DailyUsage) {
	if len(dailyUsage) == 0 {
		fmt.Println(Yellow("No usage data found"))
		return
	}
	
	PrintTitleBox("Codex CLI Token Usage Report - Daily")
	
	table := CreateCcusageStyleTable([]string{
		"Date", 
		"Models",
		"Input", 
		"Output",
		"Cache Create",
		"Cache Read", 
		"Total Tokens", 
		"Cost (USD)",
	})
	
	var totalCost float64
	var totalInput, totalOutput, totalCacheCreate, totalCacheRead, totalTokens int
	
	for _, day := range dailyUsage {
		// Calculate token breakdowns
		var inputTokens, outputTokens int
		var modelsList []string
		
		for model, usage := range day.ModelUsage {
			inputTokens += usage.PromptTokens
			outputTokens += usage.CompletionTokens
			modelsList = append(modelsList, formatModelNameSimple(model))
		}
		
		// For now, we don't have cache data, so set to 0
		cacheCreate := 0
		cacheRead := 0
		
		modelsStr := strings.Join(modelsList, ", ")
		if len(modelsStr) > 20 {
			modelsStr = modelsStr[:17] + "..."
		}
		
		table.Append([]string{
			day.Date,
			"- " + modelsStr,
			FormatNumber(inputTokens),
			FormatNumber(outputTokens),
			FormatNumber(cacheCreate),
			FormatNumber(cacheRead),
			FormatNumber(day.TotalTokens),
			FormatCurrency(day.TotalCost),
		})
		
		totalCost += day.TotalCost
		totalInput += inputTokens
		totalOutput += outputTokens
		totalCacheCreate += cacheCreate
		totalCacheRead += cacheRead
		totalTokens += day.TotalTokens
	}
	
	// Add totals row with yellow highlighting like ccusage
	table.SetFooter([]string{
		BoldYellow("Total"),
		"",
		BoldYellow(FormatNumber(totalInput)),
		BoldYellow(FormatNumber(totalOutput)),
		BoldYellow(FormatNumber(totalCacheCreate)),
		BoldYellow(FormatNumber(totalCacheRead)),
		BoldYellow(FormatNumber(totalTokens)),
		BoldYellow(FormatCurrency(totalCost)),
	})
	
	table.Render()
}

// FormatMonthlyUsageTable creates a monthly usage table
func FormatMonthlyUsageTable(monthlyUsage []types.MonthlyUsage) {
	if len(monthlyUsage) == 0 {
		fmt.Println(Yellow("No usage data found"))
		return
	}
	
	PrintTitleBox("Codex CLI Token Usage Report - Monthly")
	
	table := CreateCcusageStyleTable([]string{
		"Month", 
		"Days Active",
		"Total Requests",
		"Input Tokens",
		"Output Tokens", 
		"Total Tokens", 
		"Total Cost (USD)",
	})
	
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
			activeDays = 1 // Avoid division by zero
		}
		
		table.Append([]string{
			month.Month,
			strconv.Itoa(activeDays),
			FormatNumber(month.RequestCount),
			FormatNumber(inputTokens),
			FormatNumber(outputTokens),
			FormatNumber(month.TotalTokens),
			FormatCurrency(month.TotalCost),
		})
		
		totalCost += month.TotalCost
		totalRequests += month.RequestCount
		totalInput += inputTokens
		totalOutput += outputTokens
		totalTokens += month.TotalTokens
	}
	
	// Add totals row
	table.SetFooter([]string{
		BoldYellow("Total"),
		"",
		BoldYellow(FormatNumber(totalRequests)),
		BoldYellow(FormatNumber(totalInput)),
		BoldYellow(FormatNumber(totalOutput)),
		BoldYellow(FormatNumber(totalTokens)),
		BoldYellow(FormatCurrency(totalCost)),
	})
	
	table.Render()
}

// formatModelNameSimple returns a simple model name without colors for table display
func formatModelNameSimple(model string) string {
	// Shorten common model names for table display
	switch {
	case strings.Contains(model, "gpt-4o-mini"):
		return "gpt-4o-mini"
	case strings.Contains(model, "gpt-4o"):
		return "gpt-4o"
	case strings.Contains(model, "gpt-4"):
		return "gpt-4"
	case strings.Contains(model, "gpt-3.5-turbo"):
		return "gpt-3.5-turbo"
	case strings.Contains(model, "codex"):
		return "codex"
	default:
		return model
	}
}
