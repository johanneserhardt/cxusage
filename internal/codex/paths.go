package codex

import (
    "os"
    "path/filepath"
    "strings"

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

    // Helper: recursively collect .jsonl files under a root directory
    collectJSONL := func(root string) {
        _ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
            if err != nil {
                return nil // skip unreadable entries
            }
            if d.IsDir() {
                return nil
            }
            if strings.HasSuffix(d.Name(), ".jsonl") {
                files = append(files, path)
            }
            return nil
        })
    }

    // Check logs directory (recursively)
    if stat, err := os.Stat(paths.LogsDir); err == nil && stat.IsDir() {
        collectJSONL(paths.LogsDir)
    }

    // Check projects directory (recursively, supports session subfolders)
    if stat, err := os.Stat(paths.ProjectsDir); err == nil && stat.IsDir() {
        collectJSONL(paths.ProjectsDir)
    }

    // Check sessions directory (where Codex CLI actually stores files)
    sessionsDir := filepath.Join(paths.ConfigDir, "sessions")
    if stat, err := os.Stat(sessionsDir); err == nil && stat.IsDir() {
        collectJSONL(sessionsDir)
    }

    return files, nil
}
