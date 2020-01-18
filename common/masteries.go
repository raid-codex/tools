package common

import "sort"

type Masteries struct {
	From      string   `json:"from"`
	Author    string   `json:"author"`
	Locations []string `json:"locations"`
	Offense   []string `json:"offense"`
	Defense   []string `json:"defense"`
	Support   []string `json:"support"`
}

func (m *Masteries) Sanitize() error {
	sort.SliceStable(m.Offense, func(i, j int) bool { return m.Offense[i] < m.Offense[j] })
	sort.SliceStable(m.Defense, func(i, j int) bool { return m.Defense[i] < m.Defense[j] })
	sort.SliceStable(m.Support, func(i, j int) bool { return m.Support[i] < m.Support[j] })
	sort.SliceStable(m.Locations, func(i, j int) bool { return m.Locations[i] < m.Locations[j] })
	return nil
}
