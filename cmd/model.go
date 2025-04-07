package cmd

// ESPN API response structures
type ESPNResponse struct {
	Events  []Event  `json:"events"`
	Leagues []League `json:"leagues,omitempty"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         string        `json:"date"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Status       Status        `json:"status"`
	Competitions []Competition `json:"competitions"`
	Links        []Link        `json:"links,omitempty"`
	League       League        `json:"league,omitempty"`
}

type Status struct {
	Type StatusType `json:"type"`
}

type StatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
}

type Competition struct {
	ID          string       `json:"id"`
	Date        string       `json:"date"`
	Status      Status       `json:"status,omitempty"`
	Venue       Venue        `json:"venue"`
	Competitors []Competitor `json:"competitors"`
	Details     []DetailObj  `json:"details,omitempty"`
	Notes       []Note       `json:"notes,omitempty"`
}

type Venue struct {
	FullName string `json:"fullName"`
	Name     string `json:"name"`
	Address  struct {
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"address"`
	Capacity int `json:"capacity,omitempty"`
}

type Competitor struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Score    string `json:"score"`
	HomeAway string `json:"homeAway"`
	Team     Team   `json:"team"`
	Winner   bool   `json:"winner,omitempty"`
	Form     string `json:"form,omitempty"`
	Stats    []Stat `json:"stats,omitempty"`
}

type Team struct {
	ID               string `json:"id"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Logo             string `json:"logo"`
	Color            string `json:"color,omitempty"`
	Rank             int    `json:"rank,omitempty"`
}

// DetailObj struct to handle the object structure in the API response
type DetailObj struct {
	Type struct {
		ID           string `json:"id"`
		Abbreviation string `json:"abbreviation"`
		Name         string `json:"name"`
	} `json:"type"`
	Clock struct {
		Value        float64 `json:"value"`
		DisplayValue string  `json:"displayValue"`
	} `json:"clock"`
}

type Note struct {
	Type     string `json:"type"`
	Headline string `json:"headline"`
}

type Link struct {
	Rel  []string `json:"rel"`
	Href string   `json:"href"`
}

type League struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	ShortName    string `json:"shortName,omitempty"`
	Slug         string `json:"slug,omitempty"`
	LogoURL      string `json:"logos.dark.href,omitempty"`
}

type Stat struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TeamResponse for team API responses
type TeamResponse struct {
	Team       Team         `json:"team"`
	Roster     []Player     `json:"roster,omitempty"`
	NextMatch  Event        `json:"nextEvent,omitempty"`
	Record     TeamRecord   `json:"record,omitempty"`
	Statistics []Stat       `json:"statistics,omitempty"`
	Standings  TeamStanding `json:"standings,omitempty"`
}

type Player struct {
	ID           string `json:"id"`
	FullName     string `json:"fullName"`
	JerseyNumber string `json:"jersey"`
	Position     string `json:"position"`
	Age          int    `json:"age"`
	Nationality  string `json:"nationality"`
}

type TeamRecord struct {
	Wins         int `json:"wins"`
	Losses       int `json:"losses"`
	Draws        int `json:"draws"`
	GoalsFor     int `json:"goalsFor"`
	GoalsAgainst int `json:"goalsAgainst"`
}

type TeamStanding struct {
	Position int    `json:"position"`
	Points   int    `json:"points"`
	League   string `json:"league"`
	Form     string `json:"form"`
	GoalDiff int    `json:"goalDiff"`
}
