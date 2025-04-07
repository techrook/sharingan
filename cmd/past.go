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

var (
	league string
	date   string
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
`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchPastMatches()
	},
}

func init() {
	rootCmd.AddCommand(pastCmd)

	pastCmd.Flags().StringVarP(&league, "league", "l", "", "Filter by league (e.g. EPL, La Liga)")
	pastCmd.Flags().StringVarP(&date, "date", "d", "", "Filter by date (YYYY-MM-DD)")
}

// fetchPastMatches retrieves match results from the ESPN API
func fetchPastMatches() {
	queryDate := date
	if queryDate == "" {
		queryDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/scoreboard?dates=%s", queryDate)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	fmt.Printf("Fetching past football matches from ESPN API for %s...\n", queryDate)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var espnData ESPNResponse
	err = json.Unmarshal(body, &espnData)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response sample: %s\n", string(body[:min(1000, len(body))]))
		return
	}

	var completedMatches []Event

	for _, event := range espnData.Events {
		if event.Status.Type.State == "post" {
			if league == "" || strings.Contains(strings.ToLower(event.Name), strings.ToLower(league)) {
				completedMatches = append(completedMatches, event)
			}
		}
	}

	if len(completedMatches) == 0 {
		fmt.Println("No completed matches found for the selected filters.")
		return
	}

	fmt.Println("\n✅ COMPLETED MATCHES")
	fmt.Println("=================================")
	displayMatches(completedMatches)
	fmt.Printf("\nTotal completed matches: %d\n", len(completedMatches))
}
