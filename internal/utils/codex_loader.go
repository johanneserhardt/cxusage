package utils

import (
	"fmt"
	"time"

	"github.com/johanneserhardt/cxusage/internal/codex"
	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/sirupsen/logrus"
)

// LoadDailyUsageFromCodex loads daily usage data from Codex CLI local files
func LoadDailyUsageFromCodex(cfg *types.Config, startDate, endDate time.Time, logger *logrus.Logger) ([]types.DailyUsage, error) {
	logger.Info("Loading usage data from Codex CLI local files")

	// Check if Codex directory exists
	exists, err := codex.CodexDirExists(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to check Codex directory: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("Codex CLI directory not found. Make sure Codex CLI is installed and has been used")
	}

	// Parse usage entries from local files
	entries, err := codex.ParseUsageFiles(cfg, startDate, endDate, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Codex usage files: %w", err)
	}

	if len(entries) == 0 {
		return []types.DailyUsage{}, nil // Return empty slice instead of error
	}

	// Convert to API format for aggregation and calculate costs
	apiEntries := convertCodexToAPIEntriesWithCosts(entries, logger)

	return AggregateDailyUsage(apiEntries), nil
}

// LoadMonthlyUsageFromCodex loads monthly usage data from Codex CLI local files
func LoadMonthlyUsageFromCodex(cfg *types.Config, startDate, endDate time.Time, logger *logrus.Logger) ([]types.MonthlyUsage, error) {
	logger.Info("Loading monthly usage data from Codex CLI local files")

	// First load all daily usage, then aggregate by month
	dailyUsage, err := LoadDailyUsageFromCodex(cfg, startDate, endDate, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load daily usage data: %w", err)
	}

	// Convert daily usage to monthly aggregation
	return aggregateDailyToMonthly(dailyUsage), nil
}

// convertCodexToAPIEntriesWithCosts converts Codex entries to API format and calculates costs
func convertCodexToAPIEntriesWithCosts(codexEntries []types.CodexUsageEntry, logger *logrus.Logger) []APIUsageEntry {
	var apiEntries []APIUsageEntry

	for _, entry := range codexEntries {
		cost := entry.Cost
		
		// Calculate cost if not provided
		if cost == 0 {
			calculatedCost, err := CalculateCost(entry.Model, entry.Usage)
			if err != nil {
				logger.WithError(err).WithField("model", entry.Model).Debug("Failed to calculate cost")
			} else {
				cost = calculatedCost
			}
		}
		
		apiEntry := APIUsageEntry{
			ID:      entry.RequestID,
			Model:   entry.Model,
			Created: entry.Timestamp.Unix(),
			Usage: APIUsage{
				PromptTokens:     entry.Usage.PromptTokens,
				CompletionTokens: entry.Usage.CompletionTokens,
				TotalTokens:      entry.Usage.TotalTokens,
			},
			Cost: cost,
		}
		apiEntries = append(apiEntries, apiEntry)
	}

	return apiEntries
}

// APIUsageEntry represents usage entry compatible with existing aggregation code
type APIUsageEntry struct {
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Created int64     `json:"created"`
	Usage   APIUsage  `json:"usage"`
	Cost    float64   `json:"cost"`
}

// APIUsage represents token usage compatible with existing aggregation code
type APIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}