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
		CharactersPerToken: 4.0, // OpenAI's rough estimate: ~4 chars per token
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
		if content.Type == "text" || content.Type == "output_text" || content.Type == "input_text" {
			allText.WriteString(content.Text)
			allText.WriteString(" ") // Add space between content blocks
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
	// This will use the existing CalculateCost function from utils
	// We'll need to import utils or move the calculation here
	
	// For now, use simple estimation based on gpt-4o pricing
	// Input: $5 per 1M tokens, Output: $15 per 1M tokens
	inputCost := float64(usage.PromptTokens) * 5.0 / 1000000
	outputCost := float64(usage.CompletionTokens) * 15.0 / 1000000
	
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