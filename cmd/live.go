package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

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

func init() {
	rootCmd.AddCommand(liveCmd)
}

// ESPN API response structures
type ESPNResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Status       Status        `json:"status"`
	Competitions []Competition `json:"competitions"`
	Links        []Link        `json:"links"`
}

type Status struct {
	Type struct {
		State       string `json:"state"`
		Completed   bool   `json:"completed"`
		Description string `json:"description"`
		Detail      string `json:"detail"`
	} `json:"type"`
}

type Competition struct {
	ID          string       `json:"id"`
	Date        string       `json:"date"`
	Status      Status       `json:"status"`
	Venue       Venue        `json:"venue"`
	Competitors []Competitor `json:"competitors"`
	Details     []DetailObj  `json:"details"` // Changed from Detail to DetailObj
}

type Venue struct {
	FullName string `json:"fullName"`
	Address  struct {
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"address"`
}

type Competitor struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Score    string `json:"score"`
	HomeAway string `json:"homeAway"`
	Team     Team   `json:"team"`
}

type Team struct {
	ID               string `json:"id"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Logo             string `json:"logo"`
}

// DetailObj struct to handle the object structure in the API response
type DetailObj struct {
	Type struct {
		ID           string `json:"id"`
		Abbreviation string `json:"abbreviation"`
		Name         string `json:"name"`
	} `json:"type"`
	Clock struct { // Changed from string to struct
		Value        float64 `json:"value"`
		DisplayValue string  `json:"displayValue"`
	} `json:"clock"`
}

type Link struct {
	Rel  []string `json:"rel"`
	Href string   `json:"href"`
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

	if err := saveResponseToFile(body, "espn_response.json"); err != nil {
		log.Printf("Warning: Failed to save response to file: %v", err)
	}

	var espnData ESPNResponse
	err = json.Unmarshal(body, &espnData)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response sample: %s\n", string(body[:min(1000, len(body))]))
		return
	}

	// Display matches
	if len(espnData.Events) == 0 {
		fmt.Println("No matches found for today.")
		return
	}

	// Group matches by state (live, upcoming, completed)
	var liveMatches, upcomingMatches, completedMatches []Event

	for _, event := range espnData.Events {
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
		fmt.Println("\n🔴 LIVE MATCHES")
		fmt.Println("=================================")
		displayMatches(liveMatches)
	}

	// Display upcoming matches
	if len(upcomingMatches) > 0 {
		fmt.Println("\n⏳ UPCOMING MATCHES")
		fmt.Println("=================================")
		displayMatches(upcomingMatches)
	}

	// Display completed matches
	if len(completedMatches) > 0 {
		fmt.Println("\n✅ COMPLETED MATCHES")
		fmt.Println("=================================")
		displayMatches(completedMatches)
	}

	// Display total count
	fmt.Printf("\nTotal matches: %d (Live: %d, Upcoming: %d, Completed: %d)\n",
		len(espnData.Events), len(liveMatches), len(upcomingMatches), len(completedMatches))
}

// Helper function to save response to file for debugging
func saveResponseToFile(data []byte, filename string) error {
	return nil // Will implement if needed
}

// displayMatches formats and displays a list of matches
func displayMatches(events []Event) {
	for _, event := range events {
		if len(event.Competitions) == 0 {
			continue
		}

		competition := event.Competitions[0]

		// Extract home and away teams
		var homeTeam, awayTeam Competitor
		for _, competitor := range competition.Competitors {
			if competitor.HomeAway == "home" {
				homeTeam = competitor
			} else {
				awayTeam = competitor
			}
		}

		// Format venue information
		venue := ""
		if competition.Venue.FullName != "" {
			venue = competition.Venue.FullName
			if competition.Venue.Address.City != "" {
				venue += ", " + competition.Venue.Address.City
			}
		}

		fmt.Println("---------------------------------")
		fmt.Printf("⚽  %s vs %s\n",
			defaultIfEmpty(homeTeam.Team.DisplayName, "Home Team"),
			defaultIfEmpty(awayTeam.Team.DisplayName, "Away Team"))

		// Show scores
		fmt.Printf("📊  Score: %s - %s\n",
			defaultIfEmpty(homeTeam.Score, "0"),
			defaultIfEmpty(awayTeam.Score, "0"))

		// Match status
		switch event.Status.Type.State {
		case "in":
			fmt.Printf("⏱️  Status: %s\n", event.Status.Type.Detail)
		case "pre":
			// Parse and format the date for upcoming matches
			matchTime, err := time.Parse(time.RFC3339, event.Date)
			if err == nil {
				fmt.Printf("🕒  Kickoff: %s\n", matchTime.Format("Mon Jan 2 15:04 MST"))
			} else {
				fmt.Printf("🕒  Kickoff: %s\n", event.Status.Type.Detail)
			}
		case "post":
			fmt.Printf("🏁  Final: %s\n", event.Status.Type.Detail)
		}

		// Show venue if available
		if venue != "" {
			fmt.Printf("🏟️  Venue: %s\n", venue)
		}

		// Show match name/competition
		parts := strings.Split(event.Name, " - ")
		if len(parts) > 1 {
			fmt.Printf("🏆  Competition: %s\n", parts[0])
		}
	}
	fmt.Println("=================================")
}

func defaultIfEmpty(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
