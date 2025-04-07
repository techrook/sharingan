package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "sharingan",
	Short: "A CLI tool for fetching live scores, past matches, and team stats for different sports.",
	Long:  `Sharingan is a CLI tool for retrieving real-time and past match data for football and other sports.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'sharingan help' to see available commands")
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Common variables used across commands
var (
	league    string
	date      string
	team      string
	dateRange int
	detailed  bool
	format    string
)

// Initialize commands
func init() {
	// Commands are added in their respective files
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func defaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
