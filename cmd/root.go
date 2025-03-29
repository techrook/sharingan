package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sharingan",
	Short: "Sharingan CLI fetches live sports scores",
	Long:  `A CLI tool for fetching live scores, past matches, and team stats for different sports.`,
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

func init() {
	rootCmd.AddCommand(liveCmd)
}
