package utils

import (
	"github.com/johanneserhardt/cxusage/internal/types"
)

// aggregateDailyToMonthly converts daily usage data to monthly summaries
func aggregateDailyToMonthly(dailyUsage []types.DailyUsage) []types.MonthlyUsage {
	monthlyMap := make(map[string]*types.MonthlyUsage)

	for _, day := range dailyUsage {
		// Extract month from date (YYYY-MM-DD -> YYYY-MM)
		month := day.Date[:7]
		
		if _, exists := monthlyMap[month]; !exists {
			monthlyMap[month] = &types.MonthlyUsage{
				Month:          month,
				ModelUsage:     make(map[string]types.Usage),
				ModelCosts:     make(map[string]float64),
				DailyBreakdown: []types.DailyUsage{},
			}
		}

		monthly := monthlyMap[month]
		monthly.RequestCount += day.RequestCount
		monthly.TotalTokens += day.TotalTokens
		monthly.TotalCost += day.TotalCost
		monthly.DailyBreakdown = append(monthly.DailyBreakdown, day)

		// Aggregate model usage
		for model, usage := range day.ModelUsage {
			if existingUsage, exists := monthly.ModelUsage[model]; exists {
				existingUsage.PromptTokens += usage.PromptTokens
				existingUsage.CompletionTokens += usage.CompletionTokens
				existingUsage.TotalTokens += usage.TotalTokens
				monthly.ModelUsage[model] = existingUsage
			} else {
				monthly.ModelUsage[model] = usage
			}
		}

		// Aggregate model costs
		for model, cost := range day.ModelCosts {
			monthly.ModelCosts[model] += cost
		}
	}

	// Convert map to sorted slice
	var monthlyUsage []types.MonthlyUsage
	for _, monthly := range monthlyMap {
		monthlyUsage = append(monthlyUsage, *monthly)
	}

	// Sort by month
	for i := 0; i < len(monthlyUsage)-1; i++ {
		for j := i + 1; j < len(monthlyUsage); j++ {
			if monthlyUsage[i].Month > monthlyUsage[j].Month {
				monthlyUsage[i], monthlyUsage[j] = monthlyUsage[j], monthlyUsage[i]
			}
		}
	}

	return monthlyUsage
}