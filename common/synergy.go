package common

import "sort"

type SynergyContext struct {
	Key SynergyContextKey `json:"key"`
}

type Synergy struct {
	Context   SynergyContext `json:"context"`
	Champions []string       `json:"champions"`
}

type SynergyContextKey string

const (
	SynergyContextKey_PoisonCounterattack SynergyContextKey = "poison-counterattack"
)

func (s *Synergy) Sanitize() error {
	championsUnique := map[string]bool{}
	for _, champion := range s.Champions {
		championsUnique[champion] = true
	}
	newChampions := make([]string, len(championsUnique))
	idx := 0
	for champion := range championsUnique {
		newChampions[idx] = champion
		idx++
	}
	sort.SliceStable(newChampions, func(i, j int) bool { return newChampions[i] < newChampions[j] })
	s.Champions = newChampions
	return nil
}
