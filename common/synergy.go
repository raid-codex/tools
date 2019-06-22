package common

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
	newChampions := make([]string, 0)
	for _, champion := range s.Champions {
		if !championsUnique[champion] {
			championsUnique[champion] = true
			newChampions = append(newChampions, champion)
		}
	}
	s.Champions = newChampions
	return nil
}

var (
	SynergyData = map[SynergyContextKey]struct {
		Title          string
		RawDescription string
	}{
		SynergyContextKey_PoisonCounterattack: {
			Title:          "Poison and Counterattack",
			RawDescription: "Mixing a champion having A1 applying a Poison debuff, and a champion able to place a counterattack buff on him, is a very good situational synergy that can be impressive during Clan Boss battles.",
		},
	}
)
