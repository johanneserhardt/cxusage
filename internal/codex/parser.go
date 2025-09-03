package codex

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/johanneserhardt/cxusage/internal/types"
	"github.com/sirupsen/logrus"
)

// ParseUsageFiles parses all Codex usage log files and returns usage entries with token estimation
func ParseUsageFiles(cfg *types.Config, startDate, endDate time.Time, logger *logrus.Logger) ([]types.CodexUsageEntry, error) {
	files, err := GetUsageLogFiles(cfg)
	if err != nil {
		return nil, err
	}

	logger.WithField("files_count", len(files)).Info("Found Codex usage log files")

	var allEntries []types.CodexUsageEntry

	for _, file := range files {
		entries, err := parseCodexSessionFile(file, startDate, endDate, logger)
		if err != nil {
			logger.WithError(err).WithField("file", filepath.Base(file)).Warn("Failed to parse log file")
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	logger.WithField("total_entries", len(allEntries)).Info("Parsed Codex usage entries with token estimation")
	return allEntries, nil
}

// parseCodexSessionFile parses a complete Codex session file
func parseCodexSessionFile(filename string, startDate, endDate time.Time, logger *logrus.Logger) ([]types.CodexUsageEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []types.CodexUsageEntry
	scanner := bufio.NewScanner(file)
	
	// Increase buffer size for very large Codex CLI messages
	const maxCapacity = 10 * 1024 * 1024 // 10MB buffer for huge messages
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	
	estimator := NewTokenEstimator()
	var sessionTimestamp time.Time
	var sessionID string
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse each JSONL line
		if lineNum == 1 {
			// First line contains session metadata with timestamp
			sessionData, err := parseSessionMetadata(line)
			if err == nil {
				sessionTimestamp = sessionData.Timestamp
				sessionID = sessionData.ID
			}
		}

		// Try to parse as a message
		entry, err := parseMessageEntry(line, sessionTimestamp, sessionID, estimator)
		if err != nil {
			// Skip non-message entries (metadata, state, etc.)
			continue
		}

		// Filter by date range
		if entry.Timestamp.After(startDate) && entry.Timestamp.Before(endDate) {
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
	}).Debug("Parsed Codex session file")

	return entries, nil
}

// SessionMetadata represents the first line of a Codex session file
type SessionMetadata struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"-"`
	RawTime   string    `json:"timestamp"`
}

// parseSessionMetadata parses session metadata from first line
func parseSessionMetadata(line string) (*SessionMetadata, error) {
	var metadata SessionMetadata
	
	// Parse the JSON to extract basic fields
	var rawData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &rawData); err != nil {
		return nil, err
	}
	
	// Extract ID
	if id, ok := rawData["id"].(string); ok {
		metadata.ID = id
	}
	
	// Extract and parse timestamp
	if timeStr, ok := rawData["timestamp"].(string); ok {
		if timestamp, err := time.Parse(time.RFC3339, timeStr); err == nil {
			metadata.Timestamp = timestamp
		}
	}
	
	return &metadata, nil
}

// parseMessageEntry parses a message entry and estimates tokens
func parseMessageEntry(line string, sessionTimestamp time.Time, sessionID string, estimator *TokenEstimator) (types.CodexUsageEntry, error) {
	var entry types.CodexUsageEntry

	// Parse the message
	msg, err := ParseCodexMessage(line)
	if err != nil {
		return entry, err
	}

	// Only process actual messages with content
	if msg.Type != "message" || len(msg.Content) == 0 {
		return entry, errors.New("not a message entry")
	}

	// Use session timestamp if message doesn't have one
	if msg.Timestamp != "" {
		if timestamp, err := time.Parse(time.RFC3339, msg.Timestamp); err == nil {
			entry.Timestamp = timestamp
		} else {
			entry.Timestamp = sessionTimestamp
		}
	} else {
		entry.Timestamp = sessionTimestamp
	}

	// Estimate tokens from message content
	inputTokens, outputTokens := estimator.EstimateTokensFromMessage(*msg)
	
	// Extract model information
	model := ExtractModelFromMessage(*msg)
	
	// Calculate estimated cost
	cost := estimator.EstimateCostFromTokens(model, inputTokens, outputTokens)

	// Populate the usage entry
	entry.RequestID = msg.ID
	if entry.RequestID == "" {
		entry.RequestID = sessionID + "-msg"
	}
	entry.SessionID = sessionID
	entry.Model = model
	entry.Usage = types.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: outputTokens,
		TotalTokens:      inputTokens + outputTokens,
	}
	entry.Cost = cost

	return entry, nil
}

// Legacy functions (simplified)
func parseLogLine(line string) (types.CodexUsageEntry, error) {
	return parseMessageEntry(line, time.Now(), "unknown", NewTokenEstimator())
}

func convertGenericEntry(data map[string]interface{}) (types.CodexUsageEntry, error) {
	var entry types.CodexUsageEntry
	return entry, errors.New("legacy function")
}

func convertUsage(usage types.Usage) types.Usage {
	return usage
}