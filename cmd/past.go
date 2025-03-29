/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pastCmd represents the past command
var pastCmd = &cobra.Command{
	Use:   "past",
	Short: "Fetch past football match results",
	Long: `The 'past' command retrieves scores of previously played football matches. 
It supports fetching all past matches or filtering by specific leagues and dates.

Examples:
  # Get all past football match results
  sharingan past

  # Get past results for the English Premier League (EPL)
  sharingan past --league EPL

  # Get past results from a specific date (YYYY-MM-DD)
  sharingan past --date 2024-03-20
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("past called")
	},
}

func init() {
	rootCmd.AddCommand(pastCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pastCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pastCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
