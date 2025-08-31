package types

import (
	"time"
)

// UsageData represents OpenAI API usage information
type UsageData struct {
	ID           string    `json:"id"`
	Object       string    `json:"object"`
	Model        string    `json:"model"`
	Created      int64     `json:"created"`
	Usage        Usage     `json:"usage"`
	RequestType  string    `json:"request_type"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
	Organization string    `json:"organization,omitempty"`
	User         string    `json:"user,omitempty"`
}

// Usage represents token usage for a request
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// DailyUsage represents aggregated usage data for a single day
type DailyUsage struct {
	Date         string             `json:"date"`
	TotalCost    float64            `json:"total_cost"`
	TotalTokens  int                `json:"total_tokens"`
	RequestCount int                `json:"request_count"`
	ModelUsage   map[string]Usage   `json:"model_usage"`
	ModelCosts   map[string]float64 `json:"model_costs"`
}

// MonthlyUsage represents aggregated usage data for a month
type MonthlyUsage struct {
	Month        string             `json:"month"`
	TotalCost    float64            `json:"total_cost"`
	TotalTokens  int                `json:"total_tokens"`
	RequestCount int                `json:"request_count"`
	DailyBreakdown []DailyUsage     `json:"daily_breakdown"`
	ModelUsage   map[string]Usage   `json:"model_usage"`
	ModelCosts   map[string]float64 `json:"model_costs"`
}

// SessionUsage represents usage data grouped by session/project
type SessionUsage struct {
	SessionID    string             `json:"session_id"`
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	Duration     time.Duration      `json:"duration"`
	TotalCost    float64            `json:"total_cost"`
	TotalTokens  int                `json:"total_tokens"`
	RequestCount int                `json:"request_count"`
	ModelUsage   map[string]Usage   `json:"model_usage"`
	ModelCosts   map[string]float64 `json:"model_costs"`
}

// Config represents application configuration (updated for local file reading)
type Config struct {
	LogLevel     string `mapstructure:"log_level"`
	LocalLogging bool   `mapstructure:"local_logging"`
	LogsDir      string `mapstructure:"logs_dir"`
	CodexPath    string `mapstructure:"codex_path"` // Optional custom codex directory
}

// OutputFormat represents the output format for CLI commands
type OutputFormat string

const (
	OutputFormatTable OutputFormat = "table"
	OutputFormatJSON  OutputFormat = "json"
)

// CostMode represents how costs should be calculated
type CostMode string

const (
	CostModeAuto      CostMode = "auto"
	CostModeCalculate CostMode = "calculate"
	CostModeDisplay   CostMode = "display"
)