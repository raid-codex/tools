package common

import (
	"fmt"
	"sort"
	"strings"
)

type Skill struct {
	Name           string          `json:"name"`
	RawDescription string          `json:"raw_description"`
	Slug           string          `json:"slug"`
	Effects        []*StatusEffect `json:"effects"`
	DamageBasedOn  []string        `json:"damaged_based_on"`
}

func (s *Skill) Sanitize() error {
	s.Slug = GetLinkNameFromSanitizedName(s.Name)
	if s.Effects == nil || len(s.Effects) == 0 {
		effects, basedOn, err := getEffectsFromDescription(s.Effects, s.DamageBasedOn, s.RawDescription)
		if err != nil {
			return err
		}
		s.Effects = effects
		s.DamageBasedOn = basedOn
	} else {
		for _, effect := range s.Effects {
			errSanitize := effect.Sanitize()
			if errSanitize != nil {
				return errSanitize
			}
		}
	}
	return nil
}

func getEffectsFromDescription(effects []*StatusEffect, damageBasedOn []string, rawDescription string) ([]*StatusEffect, []string, error) {
	currentEffects := map[string]*StatusEffect{}
	/*	for _, effect := range effects {
		currentEffects[effect.Type] = effect
	}*/
	basedOn := map[string]bool{}
	/*	for _, stat := range damageBasedOn {
		basedOn[stat] = true
	}*/
	extract := statusEffectFromSkillAuraMore.FindAllString(rawDescription, -1)
	for _, m := range extract {
		m2 := statusEffectFromSkillAura.FindAllStringSubmatch(m, -1)
		if len(m2) != 1 {
			panic(fmt.Sprintf("wtf %+v", m2))
		}
		val := m2[0][1]
		switch true {
		case strings.Contains(m, fmt.Sprintf("is under a [%s]", val)),
			strings.Contains(m, fmt.Sprintf("is under [%s]", val)),
			strings.Contains(m, fmt.Sprintf("under [%s]", val)):
			// this is pure garbage
			continue
		}
		if _, ok := debuffs[val]; ok {
			currentEffects[val] = &StatusEffect{
				EffectType: "debuff",
				Type:       val,
			}
		} else if _, ok := buffs[val]; ok {
			currentEffects[val] = &StatusEffect{
				EffectType: "buff",
				Type:       val,
			}
		} else if _, ok := stats[val]; ok {
			basedOn[val] = true
		} else {
			return nil, nil, fmt.Errorf("unknown thing: %+v", m2)
		}
		if v, ok := buffDebuffRateExtraSlug[val]; ok && strings.Contains(m, v) {
			currentEffects[val].Extra = true
		}
	}
	newEffects := make([]*StatusEffect, len(currentEffects))
	idx := 0
	for _, effect := range currentEffects {
		errSanitize := effect.Sanitize()
		if errSanitize != nil {
			return nil, nil, errSanitize
		}
		newEffects[idx] = effect
		idx++
	}
	newBasedOn := make([]string, len(basedOn))
	idx = 0
	for stat, ok := range basedOn {
		if !ok {
			continue
		}
		newBasedOn[idx] = stat
		idx++
	}
	sort.SliceStable(newEffects, func(i, j int) bool {
		return newEffects[i].Slug < newEffects[j].Slug
	})
	sort.SliceStable(newBasedOn, func(i, j int) bool {
		return newBasedOn[i] < newBasedOn[j]
	})
	return newEffects, newBasedOn, nil
}
