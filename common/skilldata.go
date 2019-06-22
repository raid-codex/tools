package common

type SkillData struct {
	Level     string          `json:"level"`
	Hits      int64           `json:"hits"`
	Target    *Target         `json:"target"`
	Effects   []*StatusEffect `json:"effects"`
	BasedOn   []string        `json:"based_on"`
	Cooldown  int64           `json:"cooldown"`
	RawDetail string          `json:"raw_detail"`
}

func (sd *SkillData) Sanitize() error {
	errTarget := sd.Target.Sanitize()
	if errTarget != nil {
		return errTarget
	}
	if sd.Effects == nil {
		sd.Effects = make([]*StatusEffect, 0)
	}
	for _, effect := range sd.Effects {
		errSanitize := effect.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}
	if sd.BasedOn == nil {
		sd.BasedOn = make([]string, 0)
	}
	return nil
}
