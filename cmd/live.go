package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
)

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Fetch live football scores",
	Long: `The 'live' command retrieves real-time football scores. 
It fetches all matches scheduled for today, including live games and upcoming matches.

Examples:
  # Get all live football scores
  sharingan live
`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchLiveMatches()
	},
}

// fetchLiveMatches scrapes live and upcoming match data
func fetchLiveMatches() {
	url := "https://www.livescore.com/en/football/live/"

	// Initialize a new collector
	c := colly.NewCollector(
		colly.Async(true), // Enable async scraping
	)

	// Increase timeout settings
	c.SetRequestTimeout(20 * time.Second)

	// Set a custom user-agent to mimic a real browser
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"

	// Handle HTTP errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to fetch URL: %s\nHTTP Status: %d\nError: %s\n", r.Request.URL, r.StatusCode, err)
	})

	// Scrape matches
	c.OnHTML(".match-row", func(e *colly.HTMLElement) {
		homeTeam := e.ChildText(".team-name:first-child")
		awayTeam := e.ChildText(".team-name:last-child")
		matchTime := e.ChildText(".match-time")
		score := e.ChildText(".match-score")
		minutes := e.ChildText(".match-status")

		// Handle missing values
		if homeTeam == "" || awayTeam == "" {
			return
		}
		if matchTime == "" {
			matchTime = "TBD"
		}
		if score == "" {
			score = "0 - 0"
		}

		// Display match details
		fmt.Println("=================================")
		fmt.Printf("⚽  %s vs %s\n", homeTeam, awayTeam)

		if strings.Contains(minutes, "'") {
			// Match in progress
			fmt.Printf("⏱️  Minutes Played: %s\n", minutes)
			fmt.Printf("📊  Score: %s\n", score)
		} else {
			// Upcoming match
			fmt.Printf("⏳  Kickoff Time: %s\n", matchTime)
		}
		fmt.Println("=================================\n")
	})

	// Start scraping
	fmt.Println("Fetching live football matches...")
	err := c.Visit(url)
	if err != nil {
		log.Println("Error visiting URL:", err)
	}
	c.Wait() // Wait for async requests to complete
}
