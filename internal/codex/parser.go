package codex

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/sirupsen/logrus"
)

// ParseUsageFiles parses all Codex usage log files and returns usage entries
func ParseUsageFiles(cfg *types.Config, startDate, endDate time.Time, logger *logrus.Logger) ([]types.CodexUsageEntry, error) {
	files, err := GetUsageLogFiles(cfg)
	if err != nil {
		return nil, err
	}

	logger.WithField("files_count", len(files)).Info("Found Codex usage log files")

	var allEntries []types.CodexUsageEntry

	for _, file := range files {
		entries, err := parseJSONLFile(file, startDate, endDate, logger)
		if err != nil {
			logger.WithError(err).WithField("file", filepath.Base(file)).Warn("Failed to parse log file")
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	logger.WithField("total_entries", len(allEntries)).Info("Parsed Codex usage entries")
	return allEntries, nil
}

// parseJSONLFile parses a JSONL file and returns Codex usage entries within the date range
func parseJSONLFile(filename string, startDate, endDate time.Time, logger *logrus.Logger) ([]types.CodexUsageEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []types.CodexUsageEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry, err := parseLogLine(line)
		if err != nil {
			logger.WithError(err).WithFields(logrus.Fields{
				"file": filepath.Base(filename),
				"line": lineNum,
			}).Debug("Failed to parse log line, skipping")
			continue
		}

		// Filter by date range
		if entry.Timestamp.After(startDate) && entry.Timestamp.Before(endDate) {
			// Cost calculation moved to utils package to avoid import cycle
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"file":    filepath.Base(filename),
		"entries": len(entries),
		"lines":   lineNum,
	}).Debug("Parsed Codex log file")

	return entries, nil
}

// parseLogLine parses a single JSONL line into a CodexUsageEntry
func parseLogLine(line string) (types.CodexUsageEntry, error) {
	var entry types.CodexUsageEntry
	
	// Try to parse as CodexUsageEntry first
	if err := json.Unmarshal([]byte(line), &entry); err == nil && entry.Timestamp.IsZero() == false {
		return entry, nil
	}

	// If that fails, try parsing as a more generic log entry and convert
	var genericEntry map[string]interface{}
	if err := json.Unmarshal([]byte(line), &genericEntry); err != nil {
		return entry, err
	}

	// Convert generic entry to CodexUsageEntry
	return convertGenericEntry(genericEntry)
}

// convertGenericEntry converts a generic log entry to CodexUsageEntry
func convertGenericEntry(data map[string]interface{}) (types.CodexUsageEntry, error) {
	var entry types.CodexUsageEntry

	// Extract timestamp
	if ts, ok := data["timestamp"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, ts); err == nil {
			entry.Timestamp = parsedTime
		}
	}

	// Extract other fields
	if sessionID, ok := data["session_id"].(string); ok {
		entry.SessionID = sessionID
	}
	if requestID, ok := data["request_id"].(string); ok {
		entry.RequestID = requestID
	}
	if model, ok := data["model"].(string); ok {
		entry.Model = model
	}
	if command, ok := data["command"].(string); ok {
		entry.Command = command
	}
	if projectPath, ok := data["project_path"].(string); ok {
		entry.ProjectPath = projectPath
	}
	if cost, ok := data["cost"].(float64); ok {
		entry.Cost = cost
	}
	if duration, ok := data["duration_ms"].(float64); ok {
		entry.Duration = int64(duration)
	}

	// Extract usage data
	if usageData, ok := data["usage"].(map[string]interface{}); ok {
		if promptTokens, ok := usageData["prompt_tokens"].(float64); ok {
			entry.Usage.PromptTokens = int(promptTokens)
		}
		if completionTokens, ok := usageData["completion_tokens"].(float64); ok {
			entry.Usage.CompletionTokens = int(completionTokens)
		}
		if totalTokens, ok := usageData["total_tokens"].(float64); ok {
			entry.Usage.TotalTokens = int(totalTokens)
		}
	}

	return entry, nil
}

// convertUsage converts types.Usage to api.APIUsage for cost calculation compatibility
func convertUsage(usage types.Usage) types.Usage {
	return usage // They're the same structure
}