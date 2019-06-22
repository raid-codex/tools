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
