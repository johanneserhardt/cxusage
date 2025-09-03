package live

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/johanneserhardt/cxusage/internal/types"
)

// ProjectUsageData represents usage data split by global vs current project
type ProjectUsageData struct {
	GlobalBlock  *types.SessionBlock
	ProjectBlock *types.SessionBlock
	ProjectName  string
	ProjectPath  string
}

// ExtractProjectUsageData splits global block data into global vs project-specific views
func ExtractProjectUsageData(globalBlock *types.SessionBlock) *ProjectUsageData {
	currentDir, err := os.Getwd()
	if err != nil {
		// If we can't get current directory, just return global data
		return &ProjectUsageData{
			GlobalBlock:  globalBlock,
			ProjectBlock: &types.SessionBlock{}, // Empty project block
			ProjectName:  "unknown",
			ProjectPath:  "unknown",
		}
	}

	projectName := filepath.Base(currentDir)
	projectPath := currentDir

	// Create project-filtered block by filtering global block data
	projectBlock := &types.SessionBlock{
		StartTime:     globalBlock.StartTime,
		EndTime:       globalBlock.EndTime,
		IsActive:      globalBlock.IsActive,
		ActualEndTime: globalBlock.ActualEndTime,
		IsGap:         false,
		ModelUsage:    make(map[string]types.Usage),
		ModelCosts:    make(map[string]float64),
		Models:        []string{},
	}

	// For now, approximate project usage as a percentage of global usage
	// In a real implementation, we'd need to track which sessions belong to which projects
	// This is a simplified approach for demonstration
	projectRatio := estimateProjectRatio(projectPath)
	
	// Apply project ratio to global usage
	for model, usage := range globalBlock.ModelUsage {
		projectUsage := types.Usage{
			PromptTokens:     int(float64(usage.PromptTokens) * projectRatio),
			CompletionTokens: int(float64(usage.CompletionTokens) * projectRatio),
			TotalTokens:      int(float64(usage.TotalTokens) * projectRatio),
		}
		
		if projectUsage.TotalTokens > 0 {
			projectBlock.ModelUsage[model] = projectUsage
			projectBlock.Models = append(projectBlock.Models, model)
			projectBlock.ModelCosts[model] = globalBlock.ModelCosts[model] * projectRatio
		}
	}

	// Update project totals
	for _, usage := range projectBlock.ModelUsage {
		projectBlock.RequestCount += int(float64(globalBlock.RequestCount) * projectRatio / float64(len(globalBlock.ModelUsage)))
		projectBlock.TotalTokens += usage.TotalTokens
		projectBlock.InputTokens += usage.PromptTokens
		projectBlock.OutputTokens += usage.CompletionTokens
	}
	
	for _, cost := range projectBlock.ModelCosts {
		projectBlock.TotalCost += cost
	}

	return &ProjectUsageData{
		GlobalBlock:  globalBlock,
		ProjectBlock: projectBlock,
		ProjectName:  projectName,
		ProjectPath:  projectPath,
	}
}

// estimateProjectRatio estimates what percentage of global usage belongs to current project
func estimateProjectRatio(projectPath string) float64 {
	// Simple heuristic based on project characteristics
	// In a real implementation, this would analyze session data for project context
	
	// Check if this looks like a development project
	if isDevProject(projectPath) {
		return 0.4 // Assume 40% of usage when in active dev projects
	}
	
	// Default to smaller ratio for non-dev directories
	return 0.2 // 20% for general usage
}

// isDevProject checks if the current directory appears to be a development project
func isDevProject(projectPath string) bool {
	// Look for common development indicators
	indicators := []string{
		"package.json", "go.mod", "Cargo.toml", "pyproject.toml",
		"requirements.txt", ".git", "README.md", "src/", "cmd/",
	}
	
	for _, indicator := range indicators {
		fullPath := filepath.Join(projectPath, indicator)
		if _, err := os.Stat(fullPath); err == nil {
			return true
		}
	}
	
	return false
}

// GetProjectDisplayName returns a short, displayable project name
func GetProjectDisplayName(projectName, projectPath string) string {
	// Truncate long project names for display
	if len(projectName) > 15 {
		return projectName[:12] + "..."
	}
	
	// Show relative path for nested projects
	if strings.Contains(projectPath, "/") {
		parts := strings.Split(projectPath, "/")
		if len(parts) >= 2 {
			return parts[len(parts)-2] + "/" + projectName
		}
	}
	
	return projectName
}