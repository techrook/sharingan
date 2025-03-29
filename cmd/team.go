/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Fetch team-specific match results and stats",
	Long: `The 'team' command retrieves information about a specific football team, 
including past results, upcoming fixtures, and stats.

Examples:
  # Get all available information for a team
  sharingan team --name "Manchester United"

  # Get only past results for a team
  sharingan team --name "Real Madrid" --past

  # Get only upcoming fixtures for a team
  sharingan team --name "Barcelona" --upcoming
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("team called")
	},
}
