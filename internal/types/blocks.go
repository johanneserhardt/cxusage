package types

import (
	"time"
)

// SessionBlock represents a 5-hour billing period (like Claude's blocks)
type SessionBlock struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	ActualEndTime *time.Time `json:"actual_end_time,omitempty"`
	IsActive     bool      `json:"is_active"`
	IsGap        bool      `json:"is_gap"`
	
	// Usage data
	RequestCount int                    `json:"request_count"`
	TotalTokens  int                    `json:"total_tokens"`
	TotalCost    float64               `json:"total_cost"`
	ModelUsage   map[string]Usage      `json:"model_usage"`
	ModelCosts   map[string]float64    `json:"model_costs"`
	Models       []string              `json:"models"`
	
	// Token breakdown
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationTokens      int `json:"cache_creation_tokens"`
	CacheReadTokens          int `json:"cache_read_tokens"`
}

// BlocksConfig represents configuration for blocks command
type BlocksConfig struct {
	SessionDurationHours int
	RefreshInterval      time.Duration
	TokenLimit           *int
	ShowActive           bool
	ShowRecent           bool
	RecentDays           int
}

// LiveMonitoringConfig represents live monitoring configuration
type LiveMonitoringConfig struct {
	RefreshInterval      time.Duration
	SessionDurationHours int
	TokenLimit           *int
	ShowProjections      bool
}

// BlockProjection represents projected usage for an active block
type BlockProjection struct {
	ProjectedTokens int     `json:"projected_tokens"`
	ProjectedCost   float64 `json:"projected_cost"`
	BurnRate        float64 `json:"burn_rate"` // tokens per minute
	TimeRemaining   time.Duration `json:"time_remaining"`
}