package utils

import (
	"fmt"

	"github.com/johanneserhardt/cxusage/internal/types"
)

// Model pricing per 1M tokens (as of 2025)
var modelPricing = map[string]struct {
	InputPrice  float64 // per 1M tokens
	OutputPrice float64 // per 1M tokens
}{
	// GPT-4 models
	"gpt-4":                    {InputPrice: 30.0, OutputPrice: 60.0},
	"gpt-4-32k":               {InputPrice: 60.0, OutputPrice: 120.0},
	"gpt-4-turbo":             {InputPrice: 10.0, OutputPrice: 30.0},
	"gpt-4-turbo-preview":     {InputPrice: 10.0, OutputPrice: 30.0},
	"gpt-4-1106-preview":      {InputPrice: 10.0, OutputPrice: 30.0},
	"gpt-4-0125-preview":      {InputPrice: 10.0, OutputPrice: 30.0},
	"gpt-4-vision-preview":    {InputPrice: 10.0, OutputPrice: 30.0},
	"gpt-4o":                  {InputPrice: 5.0, OutputPrice: 15.0},
	"gpt-4o-mini":            {InputPrice: 0.15, OutputPrice: 0.6},

	// GPT-3.5 models
	"gpt-3.5-turbo":           {InputPrice: 0.5, OutputPrice: 1.5},
	"gpt-3.5-turbo-16k":       {InputPrice: 3.0, OutputPrice: 4.0},
	"gpt-3.5-turbo-0125":      {InputPrice: 0.5, OutputPrice: 1.5},
	"gpt-3.5-turbo-1106":      {InputPrice: 1.0, OutputPrice: 2.0},
	"gpt-3.5-turbo-instruct":  {InputPrice: 1.5, OutputPrice: 2.0},

	// Codex models (deprecated but might still be in logs)
	"code-davinci-002":        {InputPrice: 0.0, OutputPrice: 0.0}, // Free during beta
	"code-cushman-001":        {InputPrice: 0.0, OutputPrice: 0.0}, // Free during beta

	// Text completion models (legacy)
	"text-davinci-003":        {InputPrice: 20.0, OutputPrice: 20.0},
	"text-davinci-002":        {InputPrice: 20.0, OutputPrice: 20.0},
	"text-curie-001":          {InputPrice: 2.0, OutputPrice: 2.0},
	"text-babbage-001":        {InputPrice: 0.5, OutputPrice: 0.5},
	"text-ada-001":            {InputPrice: 0.4, OutputPrice: 0.4},

	// Embedding models
	"text-embedding-3-small": {InputPrice: 0.02, OutputPrice: 0.0},
	"text-embedding-3-large": {InputPrice: 0.13, OutputPrice: 0.0},
	"text-embedding-ada-002": {InputPrice: 0.10, OutputPrice: 0.0},

	// Fine-tuning models (base rates)
	"davinci:ft-personal":     {InputPrice: 120.0, OutputPrice: 120.0},
	"curie:ft-personal":       {InputPrice: 12.0, OutputPrice: 12.0},
	"babbage:ft-personal":     {InputPrice: 2.4, OutputPrice: 2.4},
	"ada:ft-personal":         {InputPrice: 1.6, OutputPrice: 1.6},
}

// CalculateCost calculates the cost for a given model and usage
func CalculateCost(model string, usage types.Usage) (float64, error) {
	pricing, exists := modelPricing[model]
	if !exists {
		// Try to match by prefix for fine-tuned models
		for modelPrefix, p := range modelPricing {
			if len(model) > len(modelPrefix) && model[:len(modelPrefix)] == modelPrefix {
				pricing = p
				exists = true
				break
			}
		}
		
		if !exists {
			return 0, fmt.Errorf("pricing not available for model: %s", model)
		}
	}

	// Calculate cost per million tokens
	inputCost := float64(usage.PromptTokens) * pricing.InputPrice / 1000000
	outputCost := float64(usage.CompletionTokens) * pricing.OutputPrice / 1000000
	
	return inputCost + outputCost, nil
}

// GetModelPricing returns the pricing for a specific model
func GetModelPricing(model string) (inputPrice, outputPrice float64, exists bool) {
	pricing, exists := modelPricing[model]
	if !exists {
		// Try to match by prefix for fine-tuned models
		for modelPrefix, p := range modelPricing {
			if len(model) > len(modelPrefix) && model[:len(modelPrefix)] == modelPrefix {
				pricing = p
				exists = true
				break
			}
		}
	}
	
	return pricing.InputPrice, pricing.OutputPrice, exists
}

// GetSupportedModels returns a list of all supported models
func GetSupportedModels() []string {
	models := make([]string, 0, len(modelPricing))
	for model := range modelPricing {
		models = append(models, model)
	}
	return models
}

// EstimateCost estimates the cost for a given number of tokens
func EstimateCost(model string, promptTokens, completionTokens int) (float64, error) {
	usage := types.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
	}
	
	return CalculateCost(model, usage)
}