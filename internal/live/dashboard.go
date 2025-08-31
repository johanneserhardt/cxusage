package live

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/johanneserhardt/cxusage/internal/blocks"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/johanneserhardt/cxusage/internal/utils"
)

const (
	// Dashboard dimensions
	DashboardWidth = 120
	ProgressBarWidth = 80
	
	// Colors and symbols
	SessionEmoji = "🟢"
	UsageEmoji = "🔥"
	ProjectionEmoji = "📈"
	ModelsEmoji = "⚙️"
	RefreshEmoji = "🔄"
)

// DashboardRenderer handles the beautiful live dashboard rendering
type DashboardRenderer struct {
	width int
	tokenLimit int
}

// NewDashboardRenderer creates a new dashboard renderer
func NewDashboardRenderer(tokenLimit int) *DashboardRenderer {
	return &DashboardRenderer{
		width: DashboardWidth,
		tokenLimit: tokenLimit,
	}
}

// RenderFullDashboard renders the complete ccusage-style dashboard
func (d *DashboardRenderer) RenderFullDashboard(block *types.SessionBlock, now time.Time) {
	// Clear screen and move to top
	fmt.Print("\033[2J\033[H")
	
	// Render header
	d.renderHeader()
	
	// Render session section
	d.renderSessionSection(block, now)
	
	// Render usage section
	d.renderUsageSection(block, now)
	
	// Render projection section
	d.renderProjectionSection(block)
	
	// Render models section
	d.renderModelsSection(block)
	
	// Render footer
	d.renderFooter()
}

// renderHeader renders the dashboard header
func (d *DashboardRenderer) renderHeader() {
	title := "CODEX CLI - LIVE TOKEN USAGE MONITOR"
	padding := (d.width - len(title)) / 2
	
	d.renderTopBorder()
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", padding),
		utils.BoldWhite(title),
		strings.Repeat(" ", d.width-len(title)-padding))
	d.renderSectionBorder()
}

// renderSessionSection renders the session progress section
func (d *DashboardRenderer) renderSessionSection(block *types.SessionBlock, now time.Time) {
	elapsed := now.Sub(block.StartTime)
	remaining := block.EndTime.Sub(now)
	totalDuration := block.EndTime.Sub(block.StartTime)
	
	elapsedHours := int(elapsed.Hours())
	remainingHours := int(remaining.Hours())
	
	// Calculate percentage
	progress := elapsed.Seconds() / totalDuration.Seconds()
	percentage := progress * 100
	
	// Session header
	sessionTitle := fmt.Sprintf("%s SESSION", SessionEmoji)
	progressText := fmt.Sprintf("%.1f%%", percentage)
	
	headerPadding := d.width - len(sessionTitle) - len(progressText) - 2
	fmt.Printf("│ %s%s%s │\n", 
		utils.BoldWhite(sessionTitle),
		strings.Repeat(" ", headerPadding),
		utils.BoldWhite(progressText))
	
	// Time details
	timeInfo := fmt.Sprintf("Started: %s  Elapsed: %dh  Remaining: %dh (%s)",
		utils.Cyan(block.StartTime.Format("15:04:05 AM")),
		elapsedHours,
		remainingHours,
		utils.Gray(block.EndTime.Format("15:04:05 AM")))
	
	timePadding := d.width - len(stripColors(timeInfo)) - 2
	fmt.Printf("│ %s%s │\n", timeInfo, strings.Repeat(" ", timePadding))
	
	// Progress bar
	d.renderProgressBar(progress, "green", "SESSION")
	
	d.renderSectionBorder()
}

// renderUsageSection renders the token usage section
func (d *DashboardRenderer) renderUsageSection(block *types.SessionBlock, now time.Time) {
	elapsed := now.Sub(block.StartTime)
	burnRate := float64(block.TotalTokens) / elapsed.Minutes()
	
	// Calculate usage percentage against limit
	var usagePercent float64
	var usageColorName string
	var status string
	
	if d.tokenLimit > 0 {
		usagePercent = float64(block.TotalTokens) / float64(d.tokenLimit) * 100
		if usagePercent > 100 {
			usageColorName = "red"
			status = "EXCEEDED"
		} else if usagePercent > 80 {
			usageColorName = "red"
			status = "HIGH"
		} else if usagePercent > 50 {
			usageColorName = "yellow"
			status = "MODERATE"
		} else {
			usageColorName = "green"
			status = "NORMAL"
		}
	} else {
		usagePercent = math.Min(float64(block.TotalTokens) / 50000 * 100, 100) // Assume 50k as reference
		usageColorName = "green"
		status = "TRACKING"
	}
	
	// Usage header
	usageTitle := fmt.Sprintf("%s USAGE", UsageEmoji)
	usagePercentText := fmt.Sprintf("%.1f%% (%s/%s)", 
		usagePercent,
		utils.FormatNumber(block.TotalTokens),
		utils.FormatNumber(d.tokenLimit))
	
	if d.tokenLimit == 0 {
		usagePercentText = fmt.Sprintf("%s tokens", utils.FormatNumber(block.TotalTokens))
	}
	
	headerPadding := d.width - len(stripColors(usageTitle)) - len(stripColors(usagePercentText)) - 2
	fmt.Printf("│ %s%s%s │\n", 
		utils.BoldWhite(usageTitle),
		strings.Repeat(" ", headerPadding),
		utils.BoldWhite(usagePercentText))
	
	// Usage details  
	var statusColored string
	switch usageColorName {
	case "red":
		statusColored = utils.Red(status)
	case "yellow":
		statusColored = utils.Yellow(status)
	case "green":
		statusColored = utils.Green(status)
	default:
		statusColored = status
	}
	
	usageDetails := fmt.Sprintf("Tokens: %s (Burn Rate: %s token/min ⚡ %s)  Limit: %s tokens  Cost: %s",
		utils.FormatNumber(block.TotalTokens),
		utils.Yellow(fmt.Sprintf("%.0f", burnRate)),
		statusColored,
		utils.FormatNumber(d.tokenLimit),
		d.formatCost(block.TotalCost))
	
	if d.tokenLimit == 0 {
		usageDetails = fmt.Sprintf("Tokens: %s (Burn Rate: %s token/min ⚡ %s)  Cost: %s",
			utils.FormatNumber(block.TotalTokens),
			utils.Yellow(fmt.Sprintf("%.0f", burnRate)),
			statusColored,
			d.formatCost(block.TotalCost))
	}
	
	detailsPadding := d.width - len(stripColors(usageDetails)) - 2
	fmt.Printf("│ %s%s │\n", usageDetails, strings.Repeat(" ", detailsPadding))
	
	// Usage progress bar
	d.renderProgressBar(usagePercent/100, usageColorName, "USAGE")
	
	d.renderSectionBorder()
}

// renderProjectionSection renders the projection section
func (d *DashboardRenderer) renderProjectionSection(block *types.SessionBlock) {
	projection := blocks.CalculateProjections(block)
	if projection == nil {
		return
	}
	
	// Calculate projection percentage
	var projectionPercent float64
	var projectionColorName string
	var status string
	
	if d.tokenLimit > 0 {
		projectionPercent = float64(projection.ProjectedTokens) / float64(d.tokenLimit) * 100
		if projectionPercent > 100 {
			projectionColorName = "red"
			status = "❌ WILL EXCEED LIMIT"
		} else if projectionPercent > 80 {
			projectionColorName = "red"
			status = "⚠️ APPROACHING LIMIT"
		} else if projectionPercent > 50 {
			projectionColorName = "yellow"
			status = "📊 MODERATE PROJECTION"
		} else {
			projectionColorName = "green"
			status = "✅ WITHIN LIMIT"
		}
	} else {
		projectionPercent = math.Min(float64(projection.ProjectedTokens) / 100000 * 100, 100)
		projectionColorName = "green"
		status = "📊 PROJECTED"
	}
	
	// Projection header
	projectionTitle := fmt.Sprintf("%s PROJECTION", ProjectionEmoji)
	projectionPercentText := fmt.Sprintf("%.1f%% (%s/%s)", 
		projectionPercent,
		utils.FormatNumber(projection.ProjectedTokens),
		utils.FormatNumber(d.tokenLimit))
	
	if d.tokenLimit == 0 {
		projectionPercentText = fmt.Sprintf("%s tokens", utils.FormatNumber(projection.ProjectedTokens))
	}
	
	headerPadding := d.width - len(stripColors(projectionTitle)) - len(stripColors(projectionPercentText)) - 2
	fmt.Printf("│ %s%s%s │\n", 
		utils.BoldWhite(projectionTitle),
		strings.Repeat(" ", headerPadding),
		utils.BoldWhite(projectionPercentText))
	
	// Projection details
	var statusColored string
	switch projectionColorName {
	case "red":
		statusColored = utils.Red(status)
	case "yellow":
		statusColored = utils.Yellow(status)
	case "green":
		statusColored = utils.Green(status)
	default:
		statusColored = status
	}
	
	projectionDetails := fmt.Sprintf("Status: %s  Tokens: %s  Cost: %s",
		statusColored,
		utils.FormatNumber(projection.ProjectedTokens),
		d.formatCost(projection.ProjectedCost))
	
	detailsPadding := d.width - len(stripColors(projectionDetails)) - 2
	fmt.Printf("│ %s%s │\n", projectionDetails, strings.Repeat(" ", detailsPadding))
	
	// Projection progress bar
	d.renderProgressBar(projectionPercent/100, projectionColorName, "PROJECTION")
	
	d.renderSectionBorder()
}

// renderModelsSection renders the models section
func (d *DashboardRenderer) renderModelsSection(block *types.SessionBlock) {
	if len(block.Models) == 0 {
		return
	}
	
	// Models header
	modelsTitle := fmt.Sprintf("%s Models: %s", ModelsEmoji, strings.Join(block.Models, ", "))
	modelsPadding := d.width - len(stripColors(modelsTitle)) - 2
	fmt.Printf("│ %s%s │\n", 
		utils.BoldWhite(modelsTitle),
		strings.Repeat(" ", modelsPadding))
	
	d.renderSectionBorder()
}

// renderFooter renders the dashboard footer
func (d *DashboardRenderer) renderFooter() {
	footerText := fmt.Sprintf("%s Refreshing every 1s  •  Press Ctrl+C to stop", RefreshEmoji)
	footerPadding := (d.width - len(stripColors(footerText))) / 2
	
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", footerPadding),
		utils.Gray(footerText),
		strings.Repeat(" ", d.width-len(stripColors(footerText))-footerPadding))
	
	d.renderBottomBorder()
}

// renderProgressBar renders a colored progress bar
func (d *DashboardRenderer) renderProgressBar(progress float64, colorName string, label string) {
	// Ensure progress is between 0 and 1
	progress = math.Max(0, math.Min(1, progress))
	
	filled := int(progress * float64(ProgressBarWidth))
	
	// Build progress bar
	bar := "["
	for i := 0; i < ProgressBarWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	
	// Color the bar based on colorName
	var coloredBar string
	switch colorName {
	case "red":
		coloredBar = utils.Red(bar)
	case "yellow": 
		coloredBar = utils.Yellow(bar)
	case "green":
		coloredBar = utils.Green(bar)
	default:
		coloredBar = bar
	}
	
	// Center the progress bar
	padding := (d.width - ProgressBarWidth - 4) / 2 // -4 for brackets and spaces
	fmt.Printf("│%s %s %s│\n", 
		strings.Repeat(" ", padding),
		coloredBar,
		strings.Repeat(" ", d.width-ProgressBarWidth-4-padding))
}

// renderTopBorder renders the top border
func (d *DashboardRenderer) renderTopBorder() {
	fmt.Printf("┌%s┐\n", strings.Repeat("─", d.width))
}

// renderSectionBorder renders a section separator
func (d *DashboardRenderer) renderSectionBorder() {
	fmt.Printf("├%s┤\n", strings.Repeat("─", d.width))
}

// renderBottomBorder renders the bottom border
func (d *DashboardRenderer) renderBottomBorder() {
	fmt.Printf("└%s┘\n", strings.Repeat("─", d.width))
}

// formatCost formats cost with color
func (d *DashboardRenderer) formatCost(cost float64) string {
	costStr := utils.FormatCurrency(cost)
	if cost > 5.0 {
		return utils.Red(costStr)
	} else if cost > 1.0 {
		return utils.Yellow(costStr)
	}
	return utils.Green(costStr)
}

// RenderWaitingState renders the waiting state dashboard
func (d *DashboardRenderer) RenderWaitingState(now time.Time) {
	// Render header
	d.renderHeader()
	
	// Waiting message
	waitingTitle := "⏳ WAITING FOR CODEX CLI ACTIVITY..."
	waitingPadding := (d.width - len(waitingTitle)) / 2
	
	fmt.Printf("│%s│\n", strings.Repeat(" ", d.width))
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", waitingPadding),
		utils.BoldYellow(waitingTitle),
		strings.Repeat(" ", d.width-len(waitingTitle)-waitingPadding))
	fmt.Printf("│%s│\n", strings.Repeat(" ", d.width))
	
	infoText := "No active 5-hour billing block found. Start using Codex CLI to see live usage tracking."
	infoPadding := (d.width - len(infoText)) / 2
	
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", infoPadding),
		utils.Gray(infoText),
		strings.Repeat(" ", d.width-len(infoText)-infoPadding))
	fmt.Printf("│%s│\n", strings.Repeat(" ", d.width))
	
	// Updated time
	timeText := fmt.Sprintf("Updated: %s", now.Format("2006-01-02 15:04:05"))
	timePadding := (d.width - len(timeText)) / 2
	
	fmt.Printf("│%s%s%s│\n", 
		strings.Repeat(" ", timePadding),
		utils.Gray(timeText),
		strings.Repeat(" ", d.width-len(timeText)-timePadding))
	fmt.Printf("│%s│\n", strings.Repeat(" ", d.width))
	
	// Render footer
	d.renderFooter()
}

// stripColors removes color codes for length calculation
func stripColors(s string) string {
	// This is a simplified version - in production you'd want a more robust implementation
	// For now, we'll just estimate
	return s
}