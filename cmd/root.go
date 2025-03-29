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

// Initialize commands
func init() {
	rootCmd.AddCommand(liveCmd)
	rootCmd.AddCommand(pastCmd)
	rootCmd.AddCommand(teamCmd)
}
