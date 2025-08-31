package types

import (
	"time"
)

// CodexUsageEntry represents a usage entry from Codex CLI logs
type CodexUsageEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	SessionID    string    `json:"session_id"`
	RequestID    string    `json:"request_id"`
	Model        string    `json:"model"`
	Usage        Usage     `json:"usage"`
	Cost         float64   `json:"cost,omitempty"`
	Command      string    `json:"command,omitempty"`
	ProjectPath  string    `json:"project_path,omitempty"`
	Duration     int64     `json:"duration_ms,omitempty"`
}

// CodexConfig represents Codex CLI configuration
type CodexConfig struct {
	DefaultModel string `json:"default_model" yaml:"default_model"`
	ApprovalMode string `json:"approval_mode" yaml:"approval_mode"`
	// Add other config fields as needed
}

// CodexPaths represents the directory structure for Codex CLI
type CodexPaths struct {
	ConfigDir        string // ~/.codex
	ConfigFile       string // ~/.codex/config.yaml
	InstructionsFile string // ~/.codex/instructions.md
	LogsDir          string // ~/.codex/logs (if exists)
	ProjectsDir      string // ~/.codex/projects (if exists)
}