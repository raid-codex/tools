package common

type Aura struct {
	RawDescription string          `json:"raw_description"`
	Effects        []*StatusEffect `json:"effects"`
}

func (a *Aura) Sanitize() error {
	effects, _, err := getEffectsFromDescription(a.Effects, []string{}, a.RawDescription)
	if err != nil {
		return err
	}
	a.Effects = effects
	return nil
}
