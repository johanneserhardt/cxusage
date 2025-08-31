package utils

import (
	"sort"
	"time"

	"github.com/johanneserhardt/cxusage/internal/types"
)

// AggregateDailyUsage aggregates usage entries into daily summaries
func AggregateDailyUsage(entries []APIUsageEntry) []types.DailyUsage {
	dailyMap := make(map[string]*types.DailyUsage)

	for _, entry := range entries {
		date := time.Unix(entry.Created, 0).Format("2006-01-02")
		
		if _, exists := dailyMap[date]; !exists {
			dailyMap[date] = &types.DailyUsage{
				Date:        date,
				ModelUsage:  make(map[string]types.Usage),
				ModelCosts:  make(map[string]float64),
			}
		}

		daily := dailyMap[date]
		daily.RequestCount++
		daily.TotalTokens += entry.Usage.TotalTokens
		daily.TotalCost += entry.Cost

		// Update model-specific usage
		if modelUsage, exists := daily.ModelUsage[entry.Model]; exists {
			modelUsage.PromptTokens += entry.Usage.PromptTokens
			modelUsage.CompletionTokens += entry.Usage.CompletionTokens
			modelUsage.TotalTokens += entry.Usage.TotalTokens
			daily.ModelUsage[entry.Model] = modelUsage
		} else {
			daily.ModelUsage[entry.Model] = types.Usage{
				PromptTokens:     entry.Usage.PromptTokens,
				CompletionTokens: entry.Usage.CompletionTokens,
				TotalTokens:      entry.Usage.TotalTokens,
			}
		}

		daily.ModelCosts[entry.Model] += entry.Cost
	}

	// Convert map to sorted slice
	var dailyUsage []types.DailyUsage
	for _, daily := range dailyMap {
		dailyUsage = append(dailyUsage, *daily)
	}

	// Sort by date
	sort.Slice(dailyUsage, func(i, j int) bool {
		return dailyUsage[i].Date < dailyUsage[j].Date
	})

	return dailyUsage
}

// AggregateMonthlyUsage aggregates usage entries into monthly summaries
func AggregateMonthlyUsage(entries []APIUsageEntry) []types.MonthlyUsage {
	monthlyMap := make(map[string]*types.MonthlyUsage)

	for _, entry := range entries {
		entryTime := time.Unix(entry.Created, 0)
		month := entryTime.Format("2006-01")
		
		if _, exists := monthlyMap[month]; !exists {
			monthlyMap[month] = &types.MonthlyUsage{
				Month:          month,
				ModelUsage:     make(map[string]types.Usage),
				ModelCosts:     make(map[string]float64),
				DailyBreakdown: []types.DailyUsage{},
			}
		}

		monthly := monthlyMap[month]
		monthly.RequestCount++
		monthly.TotalTokens += entry.Usage.TotalTokens
		monthly.TotalCost += entry.Cost

		// Update model-specific usage
		if modelUsage, exists := monthly.ModelUsage[entry.Model]; exists {
			modelUsage.PromptTokens += entry.Usage.PromptTokens
			modelUsage.CompletionTokens += entry.Usage.CompletionTokens
			modelUsage.TotalTokens += entry.Usage.TotalTokens
			monthly.ModelUsage[entry.Model] = modelUsage
		} else {
			monthly.ModelUsage[entry.Model] = types.Usage{
				PromptTokens:     entry.Usage.PromptTokens,
				CompletionTokens: entry.Usage.CompletionTokens,
				TotalTokens:      entry.Usage.TotalTokens,
			}
		}

		monthly.ModelCosts[entry.Model] += entry.Cost
	}

	// Generate daily breakdown for each month
	for month, monthly := range monthlyMap {
		startTime, _ := time.Parse("2006-01", month)
		endTime := startTime.AddDate(0, 1, 0).Add(-time.Second)
		
		monthEntries := filterEntriesByDateRange(entries, startTime, endTime)
		monthly.DailyBreakdown = AggregateDailyUsage(monthEntries)
	}

	// Convert map to sorted slice
	var monthlyUsage []types.MonthlyUsage
	for _, monthly := range monthlyMap {
		monthlyUsage = append(monthlyUsage, *monthly)
	}

	// Sort by month
	sort.Slice(monthlyUsage, func(i, j int) bool {
		return monthlyUsage[i].Month < monthlyUsage[j].Month
	})

	return monthlyUsage
}

// AggregateSessionUsage aggregates usage entries into session summaries
func AggregateSessionUsage(entries []types.CodexUsageEntry) []types.SessionUsage {
	sessionMap := make(map[string]*types.SessionUsage)

	for _, entry := range entries {
		// Generate session ID from timestamp (group by hour)
		sessionID := entry.Timestamp.Truncate(time.Hour).Format("2006-01-02T15")
		
		if _, exists := sessionMap[sessionID]; !exists {
			sessionMap[sessionID] = &types.SessionUsage{
				SessionID:   sessionID,
				StartTime:   entry.Timestamp,
				EndTime:     entry.Timestamp,
				ModelUsage:  make(map[string]types.Usage),
				ModelCosts:  make(map[string]float64),
			}
		}

		session := sessionMap[sessionID]
		session.RequestCount++
		session.TotalTokens += entry.Usage.TotalTokens
		session.TotalCost += entry.Cost

		// Update time range
		if entry.Timestamp.Before(session.StartTime) {
			session.StartTime = entry.Timestamp
		}
		if entry.Timestamp.After(session.EndTime) {
			session.EndTime = entry.Timestamp
		}

		// Update model-specific usage
		if modelUsage, exists := session.ModelUsage[entry.Model]; exists {
			modelUsage.PromptTokens += entry.Usage.PromptTokens
			modelUsage.CompletionTokens += entry.Usage.CompletionTokens
			modelUsage.TotalTokens += entry.Usage.TotalTokens
			session.ModelUsage[entry.Model] = modelUsage
		} else {
			session.ModelUsage[entry.Model] = types.Usage{
				PromptTokens:     entry.Usage.PromptTokens,
				CompletionTokens: entry.Usage.CompletionTokens,
				TotalTokens:      entry.Usage.TotalTokens,
			}
		}

		session.ModelCosts[entry.Model] += entry.Cost
	}

	// Calculate durations and convert to slice
	var sessionUsage []types.SessionUsage
	for _, session := range sessionMap {
		session.Duration = session.EndTime.Sub(session.StartTime)
		sessionUsage = append(sessionUsage, *session)
	}

	// Sort by start time
	sort.Slice(sessionUsage, func(i, j int) bool {
		return sessionUsage[i].StartTime.Before(sessionUsage[j].StartTime)
	})

	return sessionUsage
}

// filterEntriesByDateRange filters usage entries by date range
func filterEntriesByDateRange(entries []APIUsageEntry, startTime, endTime time.Time) []APIUsageEntry {
	var filtered []APIUsageEntry
	
	for _, entry := range entries {
		entryTime := time.Unix(entry.Created, 0)
		if entryTime.After(startTime) && entryTime.Before(endTime) {
			filtered = append(filtered, entry)
		}
	}
	
	return filtered
}