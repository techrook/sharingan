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

func init() {
	rootCmd.AddCommand(liveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// liveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// liveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
