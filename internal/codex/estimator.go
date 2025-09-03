package codex

import (
    "encoding/json"
    "strings"
    "unicode/utf8"

    "github.com/johanneserhardt/cxusage/internal/types"
)

// CodexMessage represents the structure of Codex CLI message data
type CodexMessage struct {
	Type      string    `json:"type"`
	ID        string    `json:"id,omitempty"`
	Role      string    `json:"role,omitempty"`
	Content   []Content `json:"content,omitempty"`
	Model     string    `json:"model,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
}

// Content represents message content blocks
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// TokenEstimator provides token counting functionality for Codex CLI content
type TokenEstimator struct {
    // Estimation constants based on OpenAI's guidelines
    CharactersPerToken float64
}

// NewTokenEstimator creates a new token estimator
func NewTokenEstimator() *TokenEstimator {
	return &TokenEstimator{
		CharactersPerToken: 1.5, // Adjusted based on real Codex CLI usage (8250/1290 ≈ 6.4x, so 4.0/6.4 ≈ 0.6, using 1.5 for safety)
	}
}

// EstimateTokens estimates token count from text content
func (e *TokenEstimator) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}
	
	// Count UTF-8 characters (more accurate than byte length)
	charCount := utf8.RuneCountInString(text)
	
	// Apply character-to-token ratio
	tokenEstimate := float64(charCount) / e.CharactersPerToken
	
	// Apply overhead factor based on real Codex CLI usage patterns
	// Accounts for system prompts, reasoning tokens, function calls, etc.
	overheadFactor := 2.8 // Based on 8250/1290 ≈ 6.4x observed, but conservative
	tokenEstimate *= overheadFactor
	
	// Round to nearest integer, minimum 1 for non-empty text
	if tokenEstimate < 1 && charCount > 0 {
		return 1
	}
	
	return int(tokenEstimate + 0.5) // Round to nearest
}

// EstimateTokensFromMessage estimates tokens from a Codex message
func (e *TokenEstimator) EstimateTokensFromMessage(msg CodexMessage) (inputTokens, outputTokens int) {
	// Extract all text content from the message
	var allText strings.Builder
	
    for _, content := range msg.Content {
        if strings.TrimSpace(content.Text) != "" {
            allText.WriteString(content.Text)
            allText.WriteString(" ")
        }
    }
	
	textContent := strings.TrimSpace(allText.String())
	totalTokens := e.EstimateTokens(textContent)
	
	// Determine if this is input (user) or output (assistant)
	switch msg.Role {
	case "user":
		inputTokens = totalTokens
		outputTokens = 0
	case "assistant":
		inputTokens = 0
		outputTokens = totalTokens
	default:
		// For messages without clear role, assume it's output
		inputTokens = 0
		outputTokens = totalTokens
	}
	
	return inputTokens, outputTokens
}

// EstimateCostFromTokens estimates cost based on model and token counts
func (e *TokenEstimator) EstimateCostFromTokens(model string, inputTokens, outputTokens int) float64 {
	// Default to gpt-4o if model not specified or not recognized
	if model == "" {
		model = "gpt-4o"
	}
	
	// Create usage object for cost calculation
	usage := types.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: outputTokens,
		TotalTokens:      inputTokens + outputTokens,
	}
	
	// Use existing cost calculation with estimated tokens
	cost, err := calculateCostSafely(model, usage)
	if err != nil {
		// Fallback: estimate with gpt-4o pricing if model not found
		cost, _ = calculateCostSafely("gpt-4o", usage)
	}
	
	return cost
}

// calculateCostSafely wraps cost calculation with error handling
func calculateCostSafely(model string, usage types.Usage) (float64, error) {
    // Minimal per-model pricing table (fallbacks to gpt-4o baseline)
    type rate struct{ inPerM, outPerM float64 }
    pricing := map[string]rate{
        // Known baseline (OpenAI public pricing, subject to change)
        "gpt-4o":        {inPerM: 5.0, outPerM: 15.0},
        "gpt-4o-mini":   {inPerM: 0.15, outPerM: 0.6},
        "gpt-4":         {inPerM: 10.0, outPerM: 30.0},
        "gpt-3.5-turbo": {inPerM: 0.5, outPerM: 1.5},
    }

    // Pick pricing by exact or prefix match, else fallback
    r, ok := pricing[model]
    if !ok {
        // Try loose prefix match for variants like gpt-4o-mini-...
        for k, v := range pricing {
            if strings.HasPrefix(model, k) {
                r = v
                ok = true
                break
            }
        }
    }

    if !ok {
        r = pricing["gpt-4o"]
    }

    inputCost := float64(usage.PromptTokens) * r.inPerM / 1000000.0
    outputCost := float64(usage.CompletionTokens) * r.outPerM / 1000000.0
    return inputCost + outputCost, nil
}

// ParseCodexMessage parses a JSONL line into a Codex message
func ParseCodexMessage(line string) (*CodexMessage, error) {
	var msg CodexMessage
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ExtractModelFromMessage attempts to extract model information from message
func ExtractModelFromMessage(msg CodexMessage) string {
	// Check if model is directly specified
	if msg.Model != "" {
		return msg.Model
	}
	
	// Try to extract from message content or other fields
	// For Codex CLI, we might need to infer the model
	
	// Default to gpt-4o for estimation
	return "gpt-4o"
}
