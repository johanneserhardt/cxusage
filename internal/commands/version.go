package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set during build time
	Version = "dev"
	// BuildTime is set during build time
	BuildTime = "unknown"
	// GitCommit is set during build time
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, build time, and git commit information for cxusage.",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("cxusage version %s\n", Version)
	fmt.Printf("Built: %s\n", BuildTime)
	fmt.Printf("Commit: %s\n", GitCommit)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}