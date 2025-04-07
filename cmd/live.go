package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
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

  # Get live scores for a specific league
  sharingan live --league EPL

  # Get detailed match information
  sharingan live --detailed
`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchLiveMatches()
	},
}

// displayMatches displays a list of matches based on their state
func displayMatches(matches []Event, state string) {
	for _, match := range matches {
		status := strings.ToUpper(match.Status.Type.Detail)
		homeTeam := match.Competitions[0].Competitors[0]
		awayTeam := match.Competitions[0].Competitors[1]

		fmt.Printf("%s vs %s\n", homeTeam.Team.DisplayName, awayTeam.Team.DisplayName)
		fmt.Printf("Score: %s - %s\n", homeTeam.Score, awayTeam.Score)
		fmt.Printf("Status: %s\n", status)
		fmt.Println("---------------------------------")

		if detailed {
			fmt.Printf("Venue: %s\n", match.Competitions[0].Venue.FullName)
			fmt.Printf("Start Time: %s\n", match.Date)
			fmt.Printf("League: %s\n", match.League.Name)
			fmt.Println()
		}
	}
}

func init() {
	rootCmd.AddCommand(liveCmd)

	// Add flags
	liveCmd.Flags().StringVarP(&league, "league", "l", "", "Filter by league (e.g. EPL, La Liga)")
	liveCmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed match information")
	liveCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format (pretty, json)")
}

func fetchLiveMatches() {
	url := "https://site.api.espn.com/apis/site/v2/sports/soccer/all/scoreboard"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	fmt.Println("Fetching live football matches from ESPN API...")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if format == "json" {
		fmt.Println(string(body))
		return
	}

	// Save response for debugging if DEBUG env var is set
	if os.Getenv("DEBUG") == "true" {
		if err := os.WriteFile("espn_response.json", body, 0644); err != nil {
			log.Printf("Warning: Failed to save response to file: %v", err)
		} else {
			fmt.Println("Saved raw response to espn_response.json")
		}
	}

	var espnData ESPNResponse
	err = json.Unmarshal(body, &espnData)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response sample: %s\n", string(body[:min(1000, len(body))]))
		return
	}

	// Apply league filter if specified
	var filteredEvents []Event
	if league != "" {
		for _, event := range espnData.Events {
			// Check if event matches the league filter
			isMatch := strings.Contains(strings.ToLower(event.Name), strings.ToLower(league))

			// Also check league object if available
			if event.League.Abbreviation != "" {
				isMatch = isMatch || strings.EqualFold(event.League.Abbreviation, league) ||
					strings.EqualFold(event.League.Name, league)
			}

			if isMatch {
				filteredEvents = append(filteredEvents, event)
			}
		}
	} else {
		filteredEvents = espnData.Events
	}

	// Display matches
	if len(filteredEvents) == 0 {
		fmt.Println("No matches found for today.")
		return
	}

	// Group matches by state (live, upcoming, completed)
	var liveMatches, upcomingMatches, completedMatches []Event

	for _, event := range filteredEvents {
		switch event.Status.Type.State {
		case "in":
			liveMatches = append(liveMatches, event)
		case "pre":
			upcomingMatches = append(upcomingMatches, event)
		case "post":
			completedMatches = append(completedMatches, event)
		}
	}

	// Display live matches first
	if len(liveMatches) > 0 {
		liveHeader := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Println("\n" + liveHeader("🔴 LIVE MATCHES"))
		fmt.Println("=================================")
		displayMatches(liveMatches, "live")
	}

	// Display upcoming matches
	if len(upcomingMatches) > 0 {
		upcomingHeader := color.New(color.FgYellow, color.Bold).SprintFunc()
		fmt.Println("\n" + upcomingHeader("⏳ UPCOMING MATCHES"))
		fmt.Println("=================================")
		displayMatches(upcomingMatches, "upcoming")
	}

	// Display completed matches
	if len(completedMatches) > 0 {
		completedHeader := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Println("\n" + completedHeader("✅ COMPLETED MATCHES"))
		fmt.Println("=================================")
		displayMatches(completedMatches, "completed")
	}

	// Display total count
	fmt.Printf("\nTotal matches: %d (Live: %d, Upcoming: %d, Completed: %d)\n",
		len(filteredEvents), len(liveMatches), len(upcomingMatches), len(completedMatches))
}
