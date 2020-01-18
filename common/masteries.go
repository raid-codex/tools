package common

type Masteries struct {
	From      string   `json:"from"`
	Author    string   `json:"author"`
	Locations []string `json:"locations"`
	Offense   []string `json:"offense"`
	Defense   []string `json:"defense"`
	Support   []string `json:"support"`
}
