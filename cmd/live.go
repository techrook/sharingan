/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Fetch live football scores",
	Long: `The 'live' command retrieves real-time football scores. 
It supports fetching all live matches or filtering by specific leagues.

Examples:
  # Get all live football scores
  sharingan live

  # Get live scores for the English Premier League (EPL)
  sharingan live --league EPL
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("live called")
	},
}
