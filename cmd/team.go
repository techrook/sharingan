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

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Fetch team information and stats",
	Long: `The 'team' command retrieves information about a specific team, 
including upcoming matches, recent results, and team statistics.

Examples:
  # Get information about a team
  sharingan team --name "Manchester United"

  # Get information about a team with abbreviation
  sharingan team --name MUN
`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchTeamInfo()
	},
}

func init() {
	rootCmd.AddCommand(teamCmd)

	teamCmd.Flags().StringVarP(&team, "name", "n", "", "Team name or abbreviation")
	teamCmd.MarkFlagRequired("name")
	teamCmd.Flags().StringVarP(&format, "format", "f", "pretty", "Output format (pretty, json)")
}

func fetchTeamInfo() {
	if team == "" {
		fmt.Println("Please provide a team name or abbreviation using the --name flag")
		return
	}

	// First search for the team ID
	searchURL := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/teams?limit=1000")

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	fmt.Printf("Searching for team: %s...\n", team)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// For debugging
	if os.Getenv("DEBUG") == "true" {
		if err := os.WriteFile("espn_teams_response.json", body, 0644); err != nil {
			log.Printf("Warning: Failed to save response to file: %v", err)
		}
	}

	type TeamsResponse struct {
		Sports []struct {
			Leagues []struct {
				Teams []struct {
					Team Team `json:"team"`
				} `json:"teams"`
			} `json:"leagues"`
		} `json:"sports"`
	}

	var teamsData TeamsResponse
	err = json.Unmarshal(body, &teamsData)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
		fmt.Printf("Response sample: %s\n", string(body[:min(1000, len(body))]))
		return
	}

	// Find the team
	var foundTeam Team
	var teamFound bool

	searchTerm := strings.ToLower(team)
	for _, sport := range teamsData.Sports {
		for _, league := range sport.Leagues {
			for _, t := range league.Teams {
				teamName := strings.ToLower(t.Team.DisplayName)
				teamAbbrev := strings.ToLower(t.Team.Abbreviation)

				if strings.Contains(teamName, searchTerm) || teamAbbrev == searchTerm {
					foundTeam = t.Team
					teamFound = true
					break
				}
			}
			if teamFound {
				break
			}
		}
		if teamFound {
			break
		}
	}

	if !teamFound {
		fmt.Printf("Team '%s' not found. Please check the name or abbreviation.\n", team)
		return
	}

	// Now fetch detailed team info using the ID
	teamURL := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/teams/%s", foundTeam.ID)
	req, err = http.NewRequest("GET", teamURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching team data: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if format == "json" {
		fmt.Println(string(body))
		return
	}

	// For debugging
	if os.Getenv("DEBUG") == "true" {
		if err := os.WriteFile("espn_team_detail.json", body, 0644); err != nil {
			log.Printf("Warning: Failed to save response to file: %v", err)
		}
	}

	var teamData map[string]interface{}
	err = json.Unmarshal(body, &teamData)
	if err != nil {
		log.Printf("Error parsing team JSON: %v", err)
		return
	}

	// Display team info
	titleStyle := color.New(color.FgCyan, color.Bold).SprintFunc()
	subtitleStyle := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n%s\n", titleStyle(fmt.Sprintf("TEAM: %s", foundTeam.DisplayName)))
	fmt.Println("==================================")

	// Get team logo if available
	if foundTeam.Logo != "" {
		fmt.Printf("Logo URL: %s\n", foundTeam.Logo)
	}

	// Extract and display next match
	nextMatch, hasNextMatch := teamData["nextEvent"]
	if hasNextMatch && len(nextMatch.([]interface{})) > 0 {
		fmt.Printf("\n%s\n", subtitleStyle("NEXT MATCH"))
		nextMatchData := nextMatch.([]interface{})[0].(map[string]interface{})

		// Try to extract match date
		matchDate, hasDate := nextMatchData["date"].(string)
		if hasDate {
			matchTime, err := time.Parse(time.RFC3339, matchDate)
			if err == nil {
				fmt.Printf("Date: %s\n", matchTime.Format("Mon Jan 2, 2006 15:04 MST"))
			}
		}

		// Try to extract competition info
		competitions, hasCompetitions := nextMatchData["competitions"].([]interface{})
		if hasCompetitions && len(competitions) > 0 {
			competition := competitions[0].(map[string]interface{})

			// Extract competitors
			competitors, hasCompetitors := competition["competitors"].([]interface{})
			if hasCompetitors && len(competitors) >= 2 {
				var homeTeam, awayTeam map[string]interface{}

				for _, comp := range competitors {
					competitor := comp.(map[string]interface{})
					homeAway, ok := competitor["homeAway"].(string)
					if !ok {
						continue
					}

					if homeAway == "home" {
						homeTeam = competitor
					} else {
						awayTeam = competitor
					}
				}

				if homeTeam != nil && awayTeam != nil {
					homeTeamData := homeTeam["team"].(map[string]interface{})
					awayTeamData := awayTeam["team"].(map[string]interface{})

					fmt.Printf("Match: %s vs %s\n",
						homeTeamData["displayName"].(string),
						awayTeamData["displayName"].(string))
				}
			}

			// Try to extract venue
			venue, hasVenue := competition["venue"].(map[string]interface{})
			if hasVenue {
				venueName, hasName := venue["fullName"].(string)
				if hasName {
					fmt.Printf("Venue: %s\n", venueName)
				}
			}
		}
	}

	// Display team statistics if available
	teamStats, hasStats := teamData["statistics"]
	if hasStats {
		fmt.Printf("\n%s\n", subtitleStyle("TEAM STATISTICS"))

		stats := teamStats.([]interface{})
		if len(stats) > 0 {
			for _, stat := range stats {
				statObj := stat.(map[string]interface{})
				name, hasName := statObj["name"].(string)
				value, hasValue := statObj["value"].(interface{})

				if hasName && hasValue {
					fmt.Printf("%s: %v\n", name, value)
				}
			}
		} else {
			fmt.Println("No statistics available")
		}
	}

	// Fetch recent results
	fmt.Printf("\n%s\n", subtitleStyle("RECENT RESULTS"))

	// Calculate date range for recent results (last 30 days)
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	resultsURL := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/soccer/all/teams/%s/schedule?dates=%s-%s",
		foundTeam.ID, startDate, endDate)

	req, err = http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		fmt.Println("Error fetching recent results")
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error fetching recent results")
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading results response")
		return
	}

	var scheduleData ESPNResponse
	err = json.Unmarshal(body, &scheduleData)
	if err != nil {
		fmt.Println("No recent matches found")
		return
	}

	var recentMatches []Event
	for _, event := range scheduleData.Events {
		if event.Status.Type.State == "post" {
			recentMatches = append(recentMatches, event)
		}
	}

	if len(recentMatches) == 0 {
		fmt.Println("No recent matches found")
	} else {
		displayMatches(recentMatches, "completed")
	}
}
