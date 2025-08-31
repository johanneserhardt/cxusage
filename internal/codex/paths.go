package codex

import (
	"os"
	"path/filepath"

	"github.com/johanneserhardt/cxusage/internal/types"
)

const (
	DefaultCodexDir      = ".codex"
	ConfigFileName       = "config.yaml"
	InstructionsFileName = "instructions.md"
	LogsDir              = "logs"
	ProjectsDir          = "projects"
)

// GetCodexPaths returns the Codex CLI directory paths
func GetCodexPaths(cfg *types.Config) (*types.CodexPaths, error) {
	var codexDir string
	
	// Use custom codex path if provided, otherwise use default
	if cfg.CodexPath != "" {
		codexDir = cfg.CodexPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		codexDir = filepath.Join(homeDir, DefaultCodexDir)
	}

	return &types.CodexPaths{
		ConfigDir:        codexDir,
		ConfigFile:       filepath.Join(codexDir, ConfigFileName),
		InstructionsFile: filepath.Join(codexDir, InstructionsFileName),
		LogsDir:          filepath.Join(codexDir, LogsDir),
		ProjectsDir:      filepath.Join(codexDir, ProjectsDir),
	}, nil
}

// CodexDirExists checks if the Codex CLI directory exists
func CodexDirExists(cfg *types.Config) (bool, error) {
	paths, err := GetCodexPaths(cfg)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(paths.ConfigDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

// GetUsageLogFiles returns all usage log files from Codex CLI
func GetUsageLogFiles(cfg *types.Config) ([]string, error) {
	paths, err := GetCodexPaths(cfg)
	if err != nil {
		return nil, err
	}

	var files []string

	// Check logs directory
	if stat, err := os.Stat(paths.LogsDir); err == nil && stat.IsDir() {
		logFiles, err := filepath.Glob(filepath.Join(paths.LogsDir, "*.jsonl"))
		if err != nil {
			return nil, err
		}
		files = append(files, logFiles...)
	}

	// Check projects directory (similar to Claude Code structure)
	if stat, err := os.Stat(paths.ProjectsDir); err == nil && stat.IsDir() {
		projectFiles, err := filepath.Glob(filepath.Join(paths.ProjectsDir, "*", "*.jsonl"))
		if err != nil {
			return nil, err
		}
		files = append(files, projectFiles...)
	}

	return files, nil
}