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
    // De-dup across files using a composite key
    seen := make(map[string]struct{})

    for _, file := range files {
        entries, err := parseCodexSessionFile(file, startDate, endDate, logger)
        if err != nil {
            logger.WithError(err).WithField("file", filepath.Base(file)).Warn("Failed to parse log file")
            continue
        }
        for _, e := range entries {
            key := e.SessionID + "|" + e.RequestID + "|" + e.Timestamp.Format(time.RFC3339Nano)
            if _, ok := seen[key]; ok {
                continue
            }
            seen[key] = struct{}{}
            allEntries = append(allEntries, e)
        }
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

        // Filter by date range (inclusive) so boundary events are not dropped
        if (entry.Timestamp.Equal(startDate) || entry.Timestamp.After(startDate)) &&
            (entry.Timestamp.Equal(endDate) || entry.Timestamp.Before(endDate)) {
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

// parseMessageEntry parses a message entry and prefers logged usage/cost, with estimation as fallback
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

    // Prefer logged usage/cost if present in the raw JSON
    var (
        inputTokens  int
        outputTokens int
        cost         float64
        model        string
    )

    // Default model extraction from structured message
    model = ExtractModelFromMessage(*msg)

    // Parse raw JSON for usage/cost fields (supports both top-level and message.usage)
    var raw map[string]interface{}
    if err := json.Unmarshal([]byte(line), &raw); err == nil {
        // Try to get model from raw if msg.Model missing
        if model == "" {
            if m, ok := raw["model"].(string); ok {
                model = m
            } else if m2, ok := getNestedString(raw, "message", "model"); ok {
                model = m2
            }
        }

        // Extract usage
        if usageMap, ok := raw["usage"].(map[string]interface{}); ok {
            inputTokens, outputTokens = extractUsageTokens(usageMap)
        } else if msgMap, ok := raw["message"].(map[string]interface{}); ok {
            if usageMap2, ok := msgMap["usage"].(map[string]interface{}); ok {
                inputTokens, outputTokens = extractUsageTokens(usageMap2)
            }
        }

        // Extract cost (accept costUSD or cost fields)
        if c, ok := raw["costUSD"].(float64); ok {
            cost = c
        } else if cNum, ok := raw["cost"].(float64); ok {
            cost = cNum
        } else if msgMap, ok := raw["message"].(map[string]interface{}); ok {
            if c2, ok := msgMap["costUSD"].(float64); ok {
                cost = c2
            }
        }
    }

    // If usage not present, estimate from content
    if inputTokens == 0 && outputTokens == 0 {
        inputTokens, outputTokens = estimator.EstimateTokensFromMessage(*msg)
    }

    // If cost not present, estimate from tokens and model
    if cost == 0 {
        cost = estimator.EstimateCostFromTokens(model, inputTokens, outputTokens)
    }

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

// Helper: extract tokens from a usage map with various possible keys
func extractUsageTokens(m map[string]interface{}) (in int, out int) {
    // Support both snake_case and OpenAI-style names
    if v, ok := getNumber(m, "input_tokens"); ok {
        in = v
    }
    if v, ok := getNumber(m, "output_tokens"); ok {
        out = v
    }
    if v, ok := getNumber(m, "prompt_tokens"); ok && in == 0 {
        in = v
    }
    if v, ok := getNumber(m, "completion_tokens"); ok && out == 0 {
        out = v
    }
    // If only total provided, split conservatively (assign to output)
    if in == 0 && out == 0 {
        if v, ok := getNumber(m, "total_tokens"); ok {
            out = v
        }
    }
    return in, out
}

func getNumber(m map[string]interface{}, key string) (int, bool) {
    if v, ok := m[key]; ok {
        switch t := v.(type) {
        case float64:
            return int(t), true
        case int:
            return t, true
        }
    }
    return 0, false
}

func getNestedString(m map[string]interface{}, keys ...string) (string, bool) {
    cur := any(m)
    for i, k := range keys {
        mm, ok := cur.(map[string]interface{})
        if !ok {
            return "", false
        }
        v, ok := mm[k]
        if !ok {
            return "", false
        }
        if i == len(keys)-1 {
            if s, ok := v.(string); ok {
                return s, true
            }
            return "", false
        }
        cur = v
    }
    return "", false
}
