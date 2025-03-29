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

func init() {
	rootCmd.AddCommand(teamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// teamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// teamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
