package blocks

import (
	"sort"
	"time"

	"github.com/johanneserhardt/cxusage/internal/types"
)

const (
	// DefaultSessionDurationHours is the default 5-hour billing block duration
	DefaultSessionDurationHours = 5
)

// AggregateIntoBlocks converts usage entries into 5-hour billing blocks
func AggregateIntoBlocks(entries []types.CodexUsageEntry, sessionDurationHours int) []types.SessionBlock {
	if len(entries) == 0 {
		return []types.SessionBlock{}
	}

	// Sort entries by timestamp
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	blockDuration := time.Duration(sessionDurationHours) * time.Hour
	var blocks []types.SessionBlock
	blockMap := make(map[int64]*types.SessionBlock)

	for _, entry := range entries {
		// Floor to the hour and then find the appropriate 5-hour block
		blockStartTime := floorToBlockStart(entry.Timestamp, sessionDurationHours)
		blockKey := blockStartTime.Unix()

		// Get or create block
		block, exists := blockMap[blockKey]
		if !exists {
			block = &types.SessionBlock{
				StartTime:     blockStartTime,
				EndTime:       blockStartTime.Add(blockDuration),
				ActualEndTime: nil,
				IsActive:      false,
				IsGap:         false,
				ModelUsage:    make(map[string]types.Usage),
				ModelCosts:    make(map[string]float64),
				Models:        []string{},
			}
			blockMap[blockKey] = block
		}

		// Update block data
		block.RequestCount++
		block.TotalTokens += entry.Usage.TotalTokens
		block.TotalCost += entry.Cost
		block.InputTokens += entry.Usage.PromptTokens
		block.OutputTokens += entry.Usage.CompletionTokens

		// Update model usage
		if modelUsage, exists := block.ModelUsage[entry.Model]; exists {
			modelUsage.PromptTokens += entry.Usage.PromptTokens
			modelUsage.CompletionTokens += entry.Usage.CompletionTokens
			modelUsage.TotalTokens += entry.Usage.TotalTokens
			block.ModelUsage[entry.Model] = modelUsage
		} else {
			block.ModelUsage[entry.Model] = entry.Usage
			block.Models = append(block.Models, entry.Model)
		}

		block.ModelCosts[entry.Model] += entry.Cost

		// Update actual end time
		if block.ActualEndTime == nil || entry.Timestamp.After(*block.ActualEndTime) {
			block.ActualEndTime = &entry.Timestamp
		}
	}

	// Convert map to sorted slice and determine active blocks
	now := time.Now()
	for _, block := range blockMap {
		// Check if block is currently active (within 5-hour window)
		if now.After(block.StartTime) && now.Before(block.EndTime) {
			block.IsActive = true
		}

		blocks = append(blocks, *block)
	}
	
	// If no block is currently active but we have recent data, create/mark current block as active
	hasActiveBlock := false
	for _, block := range blocks {
		if block.IsActive {
			hasActiveBlock = true
			break
		}
	}
	
	// If no active block but we have recent usage, check if we should create a current active block
	if !hasActiveBlock && len(blocks) > 0 {
		latestBlock := blocks[len(blocks)-1]
		// If the latest activity was less than 1 hour ago, consider creating a new active block
		if latestBlock.ActualEndTime != nil && now.Sub(*latestBlock.ActualEndTime) < time.Hour {
			// Create a new active block for the current time
			currentBlockStart := floorToBlockStart(now, sessionDurationHours)
			currentBlockEnd := currentBlockStart.Add(time.Duration(sessionDurationHours) * time.Hour)
			
			// Check if this would be a new block
			newBlock := true
			for i, existing := range blocks {
				if existing.StartTime.Equal(currentBlockStart) {
					// Update existing block to be active
					blocks[i].IsActive = true
					blocks[i].EndTime = currentBlockEnd
					newBlock = false
					break
				}
			}
			
			if newBlock {
				activeBlock := types.SessionBlock{
					StartTime:     currentBlockStart,
					EndTime:       currentBlockEnd,
					IsActive:      true,
					ModelUsage:    make(map[string]types.Usage),
					ModelCosts:    make(map[string]float64),
					Models:        []string{},
				}
				blocks = append(blocks, activeBlock)
			}
		}
	}

	// Sort blocks by start time
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].StartTime.Before(blocks[j].StartTime)
	})

	// Fill gaps between blocks if needed
	blocks = fillGaps(blocks, sessionDurationHours)

	return blocks
}

// floorToBlockStart floors a timestamp to the appropriate block start time
func floorToBlockStart(timestamp time.Time, sessionDurationHours int) time.Time {
	// Floor to the hour first
	floored := time.Date(
		timestamp.Year(),
		timestamp.Month(),
		timestamp.Day(),
		timestamp.Hour(),
		0, 0, 0,
		timestamp.Location(),
	)

	// Find the appropriate block boundary (every 5 hours starting from midnight)
	hourOfDay := floored.Hour()
	blockIndex := hourOfDay / sessionDurationHours
	blockStartHour := blockIndex * sessionDurationHours

	return time.Date(
		floored.Year(),
		floored.Month(),
		floored.Day(),
		blockStartHour,
		0, 0, 0,
		floored.Location(),
	)
}

// fillGaps adds gap blocks between usage blocks to show inactive periods
func fillGaps(blocks []types.SessionBlock, sessionDurationHours int) []types.SessionBlock {
	if len(blocks) <= 1 {
		return blocks
	}

	var result []types.SessionBlock
	blockDuration := time.Duration(sessionDurationHours) * time.Hour

	for i, block := range blocks {
		result = append(result, block)

		// Check if there's a gap to the next block
		if i < len(blocks)-1 {
			nextBlock := blocks[i+1]
			expectedNextStart := block.EndTime

			// If there's a gap larger than one block duration, add gap blocks
			if nextBlock.StartTime.After(expectedNextStart) {
				gapStart := expectedNextStart
				for gapStart.Before(nextBlock.StartTime) {
					gapEnd := gapStart.Add(blockDuration)
					if gapEnd.After(nextBlock.StartTime) {
						gapEnd = nextBlock.StartTime
					}

					gapBlock := types.SessionBlock{
						StartTime:     gapStart,
						EndTime:       gapEnd,
						ActualEndTime: &gapEnd,
						IsActive:      false,
						IsGap:         true,
						ModelUsage:    make(map[string]types.Usage),
						ModelCosts:    make(map[string]float64),
						Models:        []string{},
					}
					result = append(result, gapBlock)
					gapStart = gapEnd
				}
			}
		}
	}

	return result
}

// GetActiveBlock returns the currently active block, if any
func GetActiveBlock(blocks []types.SessionBlock) *types.SessionBlock {
	for i := range blocks {
		if blocks[i].IsActive && !blocks[i].IsGap {
			return &blocks[i]
		}
	}
	return nil
}

// FilterRecentBlocks returns only blocks from the last N days
func FilterRecentBlocks(blocks []types.SessionBlock, days int) []types.SessionBlock {
	if days <= 0 {
		return blocks
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	var filtered []types.SessionBlock

	for _, block := range blocks {
		if block.StartTime.After(cutoff) {
			filtered = append(filtered, block)
		}
	}

	return filtered
}

// CalculateProjections calculates projections for an active block
func CalculateProjections(block *types.SessionBlock) *types.BlockProjection {
	if !block.IsActive {
		return nil
	}

	now := time.Now()
	timeElapsed := now.Sub(block.StartTime).Minutes()
	timeRemaining := block.EndTime.Sub(now)
	
	if timeElapsed <= 0 || block.TotalTokens == 0 {
		return &types.BlockProjection{
			ProjectedTokens: block.TotalTokens,
			ProjectedCost:   block.TotalCost,
			BurnRate:        0,
			TimeRemaining:   timeRemaining,
		}
	}

	// Calculate burn rate (tokens per minute)
	burnRate := float64(block.TotalTokens) / timeElapsed
	
	// Project total usage for the full 5-hour block
	totalMinutesInBlock := block.EndTime.Sub(block.StartTime).Minutes()
	projectedTokens := int(burnRate * totalMinutesInBlock)
	
	// Project cost based on current cost rate
	costRate := block.TotalCost / timeElapsed
	projectedCost := costRate * totalMinutesInBlock

	return &types.BlockProjection{
		ProjectedTokens: projectedTokens,
		ProjectedCost:   projectedCost,
		BurnRate:        burnRate,
		TimeRemaining:   timeRemaining,
	}
}