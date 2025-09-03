package live

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/johanneserhardt/cxusage/internal/blocks"
	"github.com/johanneserhardt/cxusage/internal/codex"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/johanneserhardt/cxusage/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultRefreshInterval is the default refresh rate for live monitoring
	DefaultRefreshInterval = 1 * time.Second // Match ccusage's 1s refresh
	
	// MinRefreshInterval is the minimum allowed refresh interval
	MinRefreshInterval = 1 * time.Second
	
	// MaxRefreshInterval is the maximum allowed refresh interval
	MaxRefreshInterval = 60 * time.Second
)

// LiveMonitor handles real-time monitoring of Codex usage
type LiveMonitor struct {
	config     *types.LiveMonitoringConfig
	cfg        *types.Config
	logger     *logrus.Logger
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewLiveMonitor creates a new live monitor instance
func NewLiveMonitor(config *types.LiveMonitoringConfig, cfg *types.Config, logger *logrus.Logger) *LiveMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &LiveMonitor{
		config: config,
		cfg:    cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins live monitoring with real-time updates
func (m *LiveMonitor) Start() error {
	m.logger.Info("Starting live monitoring...")
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Setup terminal
	m.setupTerminal()
	defer m.cleanupTerminal()
	
	// Start monitoring loop
	ticker := time.NewTicker(m.config.RefreshInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-m.ctx.Done():
			return nil
		case <-sigChan:
			m.logger.Info("Received shutdown signal")
			return nil
		case <-ticker.C:
			if err := m.renderUpdate(); err != nil {
				m.logger.WithError(err).Error("Failed to render update")
			}
		}
	}
}

// Stop stops the live monitoring
func (m *LiveMonitor) Stop() {
	m.cancel()
}

// renderUpdate renders a single update of the live monitoring display
func (m *LiveMonitor) renderUpdate() error {
	// Move to top without clearing (reduces flicker)
	fmt.Print("\033[H")
	
	// Get current time for display
	now := time.Now()
	
	// Load latest usage data
	endTime := now
	startTime := now.AddDate(0, 0, -1) // Last 24 hours for active blocks
	
	entries, err := codex.ParseUsageFiles(m.cfg, startTime, endTime, m.logger)
	if err != nil {
		return fmt.Errorf("failed to load usage data: %w", err)
	}
	
	// Aggregate into blocks
	sessionBlocks := blocks.AggregateIntoBlocks(entries, m.config.SessionDurationHours)
	
	// Filter to recent blocks only
	recentBlocks := blocks.FilterRecentBlocks(sessionBlocks, 1) // Last 24 hours
	
	// Find active block
	activeBlock := blocks.GetActiveBlock(recentBlocks)
	
	if activeBlock == nil {
		m.renderWaitingState(now)
		return nil
	}
	
	// Render active block with projections
	m.renderActiveBlock(activeBlock, now)
	
	return nil
}

// renderWaitingState renders the display when no active block exists
func (m *LiveMonitor) renderWaitingState(now time.Time) {
	// Clear screen and move to top
	fmt.Print("\033[2J\033[H")
	
	dashboard := NewDashboardRenderer(50000) // Default limit for display
	dashboard.RenderWaitingState(now)
}

// renderActiveBlock renders the active block with live data and projections
func (m *LiveMonitor) renderActiveBlock(block *types.SessionBlock, now time.Time) {
	// Use token limit from config or default
	tokenLimit := 50000 // Default reference limit
	if m.config.TokenLimit != nil {
		tokenLimit = *m.config.TokenLimit
	}
	
	// Create and use the beautiful dashboard renderer
	dashboard := NewDashboardRenderer(tokenLimit)
	dashboard.RenderFullDashboard(block, now)
}

// renderProgressBar renders a visual progress bar for the 5-hour block
func (m *LiveMonitor) renderProgressBar(progress float64) {
	barWidth := 50
	filled := int(progress * float64(barWidth))
	
	fmt.Printf("Progress: [")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			fmt.Printf("%s", utils.Green("█"))
		} else {
			fmt.Printf("%s", utils.Gray("░"))
		}
	}
	fmt.Printf("] %.1f%%\n", progress*100)
	fmt.Println()
}

// formatCostWithColor formats cost with color coding
func (m *LiveMonitor) formatCostWithColor(cost float64) string {
	costStr := utils.FormatCurrency(cost)
	if cost > 1.0 {
		return utils.Red(costStr)
	} else if cost > 0.1 {
		return utils.Yellow(costStr)
	}
	return utils.Green(costStr)
}

// formatProjectionWithColor formats projection with comparison to current
func (m *LiveMonitor) formatProjectionWithColor(projected, current int) string {
	projectedStr := utils.FormatNumber(projected)
	if projected > current*2 {
		return utils.Red(projectedStr)
	} else if float64(projected) > float64(current)*1.5 {
		return utils.Yellow(projectedStr)
	}
	return utils.Green(projectedStr)
}

// formatCostProjection formats cost projection with comparison
func (m *LiveMonitor) formatCostProjection(projected, current float64) string {
	projectedStr := utils.FormatCurrency(projected)
	if projected > current*2 {
		return utils.Red(projectedStr)
	} else if projected > current*1.5 {
		return utils.Yellow(projectedStr)
	}
	return utils.Green(projectedStr)
}

// formatModelName formats model name with color
func (m *LiveMonitor) formatModelName(model string) string {
	switch {
	case model == "gpt-4o":
		return utils.Magenta(model)
	case model == "gpt-4":
		return utils.Blue(model)
	case model == "gpt-3.5-turbo":
		return utils.Green(model)
	case model == "gpt-4o-mini":
		return utils.Cyan(model)
	default:
		return model
	}
}

// setupTerminal prepares the terminal for live monitoring
func (m *LiveMonitor) setupTerminal() {
	// Hide cursor
	fmt.Print("\033[?25l")
	
	// Clear screen once at start
	fmt.Print("\033[2J\033[H")
}

// cleanupTerminal restores terminal state
func (m *LiveMonitor) cleanupTerminal() {
	// Show cursor
	fmt.Print("\033[?25h")
	
	// Clear screen one more time
	fmt.Print("\033[2J\033[H")
	
	fmt.Println("Live monitoring stopped.")
}