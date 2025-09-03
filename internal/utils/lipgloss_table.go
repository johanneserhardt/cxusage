package utils

import (
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/charmbracelet/lipgloss"
    xterm "github.com/charmbracelet/x/term"
    runewidth "github.com/mattn/go-runewidth"
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

// Global compact mode toggle, set via CLI flag
var compactMode bool
var widthOverride int

// SetCompactMode sets whether tables should use compact minimum widths
func SetCompactMode(b bool) { compactMode = b }

// isCompact reports whether compact mode is enabled
func isCompact() bool { return compactMode }

// SetWidthOverride sets a fixed table width (overrides terminal detection)
func SetWidthOverride(w int) { widthOverride = w }

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

// getTerminalWidth returns the terminal width or a reasonable default
func getTerminalWidth() int {
    if widthOverride > 0 {
        return widthOverride
    }
    if w, _, err := xterm.GetSize(uintptr(os.Stdout.Fd())); err == nil && w > 0 {
        return w
    }
    // Fallback to common width if not attached to a TTY
    return 120
}

// computeAutoWidths calculates column widths based on content and terminal width
func computeAutoWidths(headers []string, rows [][]string, min []int) []int {
    n := len(headers)
    widths := make([]int, n)
    // Seed with header width
    for i := 0; i < n; i++ {
        w := runewidth.StringWidth(headers[i])
        if w < min[i] {
            w = min[i]
        }
        widths[i] = w
    }
    // Grow with row content
    for _, row := range rows {
        for i := 0; i < n && i < len(row); i++ {
            w := runewidth.StringWidth(row[i])
            if w > widths[i] {
                widths[i] = w
            }
        }
    }

    // Calculate target max internal width between table corners
    termWidth := getTerminalWidth()
    // Between corners we have sum(width+2) + (n-1) separators
    target := termWidth - 2
    current := 0
    for _, w := range widths {
        current += w + 2
    }
    current += (n - 1)

    // Shrink proportionally from widest columns until it fits
    if current > target {
        overflow := current - target
        // Build list of indices sorted by width desc each iteration
        for overflow > 0 {
            // Find widest shrinkable column
            widest := -1
            idx := -1
            for i := 0; i < n; i++ {
                if widths[i] > min[i] && widths[i] > widest {
                    widest = widths[i]
                    idx = i
                }
            }
            if idx == -1 {
                break // cannot shrink further
            }
            widths[idx]--
            overflow--
        }
    }
    return widths
}

// FormatDailyUsageTableProper creates a proper table like ccusage
func FormatDailyUsageTableProper(dailyUsage []types.DailyUsage) {
	if len(dailyUsage) == 0 {
		fmt.Println("No usage data found")
		return
	}
	
	// Print title with border like ccusage
	title := "Codex CLI Token Usage Report - Daily (~estimated)"
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
    headers := []string{"Date", "Models", "Input", "Output", "Cache Create", "Cache Read", "Total Tokens", "Cost (USD)"}
	
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
	
    // Autosize column widths based on content and terminal width
    min := []int{10, 12, 6, 6, 6, 6, 10, 8}
    if isCompact() {
        min = []int{8, 10, 5, 5, 5, 5, 9, 8}
    }
    widths := computeAutoWidths(headers, rows, min)
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
	title := "Codex CLI Token Usage Report - Monthly (~estimated)"
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
    headers := []string{"Month", "Days Active", "Total Requests", "Input Tokens", "Output Tokens", "Total Tokens", "Total Cost (USD)"}
	
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
	
    // Autosize widths
    min := []int{7, 6, 10, 10, 10, 10, 12}
    if isCompact() {
        min = []int{6, 5, 8, 8, 8, 8, 10}
    }
    widths := computeAutoWidths(headers, rows, min)
    // Render the table
    table := CreateTable(headers, rows, widths)
    fmt.Println(table)
}

// padString pads a string to a specific width
func padString(s string, width int) string {
    // Ensure we measure and trim by display width (handles wide runes)
    w := runewidth.StringWidth(s)
    if w > width {
        // Trim to width-3 and add ellipsis if possible
        target := width
        suffix := ""
        if width > 3 {
            target = width - 3
            suffix = "..."
        }
        // Accumulate runes until reaching target width
        var b strings.Builder
        cur := 0
        for _, r := range s {
            rw := runewidth.RuneWidth(r)
            if cur+rw > target {
                break
            }
            b.WriteRune(r)
            cur += rw
        }
        return b.String() + suffix
    }
    // Pad with spaces based on visual width
    return s + strings.Repeat(" ", width-w)
}
