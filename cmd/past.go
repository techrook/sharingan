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

  # Get detailed match information
  sharingan past --detailed
`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchPastMatches()
	},
}

func init() {
	rootCmd.AddCommand(pastCmd)

	pastCmd.Flags().StringVarP(&league, "league", "l", "", "Filter by league (e.g. EPL, La Liga)")
	pastCmd.Flags().StringVarP(&date, "date", "d", "", "Filter by date (YYYY-MM-DD)")
	pastCmd.Flags().BoolVarP(&detailed, "detailed", "D", false, "Show detailed match information")
	pastCmd.Flags().IntVarP(&dateRange, "range", "r", 1, "Date range in days (for multiple days)")
	pastCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format (pretty, json)")
}

// fetchPastMatches retrieves match results from the ESPN API
func fetchPastMatches() {
	// Set date(s) for the query
	startDate := date
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	// Calculate end date if range is specified
	endDate := ""
	if dateRange > 1 {
		endDateObj, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			log.Fatalf("Invalid date format: %v", err)
		}
		endDate = endDateObj.AddDate(0, 0, dateRange-1).Format("2006-01-02")
	}

	// Build URL based on date(s)
	var url string
	if endDate != "" {
		url = fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/scoreboard?dates=%s", startDate)
	} else {
		url = fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/scoreboard?dates=%s", startDate)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Logging the request URL for debugging
	fmt.Printf("Request URL: %s\n", url)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Log the raw response if debugging
	if os.Getenv("DEBUG") == "true" {
		if err := os.WriteFile("espn_past_response.json", body, 0644); err != nil {
			log.Printf("Warning: Failed to save response to file: %v", err)
		} else {
			fmt.Println("Saved raw response to espn_past_response.json")
		}
	}

	// If the format is JSON, output the raw response
	if format == "json" {
		fmt.Println(string(body))
		return
	}

	// Parse the JSON response
	var espnData ESPNResponse
	err = json.Unmarshal(body, &espnData)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response sample: %s\n", string(body[:min(1000, len(body))]))
		return
	}

	// Debug: print the raw events to inspect their structure
	for _, event := range espnData.Events {
		fmt.Printf("Event: %s | Status: %+v\n", event.Name, event.Status.Type.State)
	}

	// Filter completed matches
	var completedMatches []Event
	for _, event := range espnData.Events {
		if event.Status.Type.State == "post" {
			// Apply league filter if specified
			if league == "" || strings.Contains(strings.ToLower(event.Name), strings.ToLower(league)) ||
				(event.League.Abbreviation != "" && (strings.EqualFold(event.League.Abbreviation, league) ||
					strings.EqualFold(event.League.Name, league))) {
				completedMatches = append(completedMatches, event)
			}
		}
	}

	// If no completed matches are found, print a message
	if len(completedMatches) == 0 {
		fmt.Println("No completed matches found for the selected filters.")
		return
	}

	// Print header for completed matches
	completedHeader := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Println("\n" + completedHeader("✅ COMPLETED MATCHES"))
	fmt.Println("=================================")

	// Display the matches
	displayMatches(completedMatches, "completed")
	fmt.Printf("\nTotal completed matches: %d\n", len(completedMatches))
}
