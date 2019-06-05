package common

import (
	"fmt"
	"sort"
	"strings"
)

type Skill struct {
	Passive        bool            `json:"passive"`
	Name           string          `json:"name"`
	RawDescription string          `json:"raw_description"`
	Slug           string          `json:"slug"`
	Effects        []*StatusEffect `json:"effects"`
	DamageBasedOn  []string        `json:"damaged_based_on"`
	GIID           string          `json:"giid"`
	Upgrades       []*SkillData    `json:"upgrades"`
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
	for _, upgrade := range s.Upgrades {
		errUpgrade := upgrade.Sanitize()
		if errUpgrade != nil {
			return errUpgrade
		}
	}
	sort.SliceStable(s.Upgrades, func(i, j int) bool {
		return s.Upgrades[i].Level < s.Upgrades[j].Level
	})
	return nil
}

func (s *Skill) SetSkillData(sd *SkillData) {
	ns := make([]*SkillData, 0)
	ns = append(ns, sd)
	for _, skillData := range s.Upgrades {
		if skillData.Level != sd.Level {
			ns = append(ns, skillData)
		}
	}
	s.Upgrades = ns
}

type SkillData struct {
	Level   string          `json:"level"`
	Hits    int64           `json:"hits"`
	Target  *Target         `json:"target"`
	Effects []*StatusEffect `json:"effects"`
	BasedOn []string        `json:"based_on"`
}

func (sd *SkillData) Sanitize() error {
	errTarget := sd.Target.Sanitize()
	if errTarget != nil {
		return errTarget
	}
	for _, effect := range sd.Effects {
		errSanitize := effect.Sanitize()
		if errSanitize != nil {
			return errSanitize
		}
	}
	return nil
}

func translateEffect(str string) string {
	if _, ok := effectTranslate[str]; ok {
		return effectTranslate[str]
	} else if strings.HasPrefix(str, "Dec ") {
		return strings.Replace(str, "Dec ", "Decrease ", -1)
	} else if strings.HasPrefix(str, "Inc ") {
		return strings.Replace(str, "Inc ", "Increase ", -1)
	}
	return str
}

var (
	effectTranslate = map[string]string{
		"Heal Red":           "Heal Reduction",
		"Crit":               "Critical Strike",
		"Inc Meter":          "Increase Turn Meter",
		"Dec Meter":          "Decrease Turn Meter",
		"Reflect DMG":        "Reflect Damage",
		"Inc Crit Rate":      "Increase C. RATE",
		"Skills on Cooldown": "Block Cooldown Skills",
		"Skill on Cooldown":  "Block Cooldown Skills",
		"Inc Skill Cooldown": "Block Cooldown Skills",
		"Cont Heal":          "Continuous Heal",
		"Counter":            "Counterattack",
		"Ally Prot":          "Ally Protection",
		"Block DMG":          "Block Damage",
		"Fill Meter":         "Increase Turn Meter",
		"cont heal":          "Continuous Heal",
		"Empty Meter":        "Decrease Turn Meter",
		"block debuffs":      "Block Debuffs",
		"Remove All Debuffs": "Remove ALL Debuffs",
	}
)

func (sd *SkillData) AddEffect(effect string, who string, turns int64, chance float64, placesIf string, value float64) {
	realEffect := translateEffect(effect)
	placesIf = translateEffect(placesIf)
	statusEffect := &StatusEffect{
		Type:     realEffect,
		Turns:    turns,
		Target:   &Target{Who: who},
		Chance:   chance,
		Value:    value,
		PlacesIf: placesIf,
	}
	if _, ok := debuffs[statusEffect.Type]; ok {
		statusEffect.EffectType = "debuff"
	} else if _, ok := buffs[statusEffect.Type]; ok {
		statusEffect.EffectType = "debuff"
	} else if _, ok := battleEnhancements[statusEffect.Type]; ok {
		statusEffect.EffectType = "battle_enhancement"
	} else {
		panic(fmt.Errorf("unknown effect %s", statusEffect.Type))
	}
	ne := make([]*StatusEffect, 0)
	ne = append(ne, statusEffect)
	for _, se := range sd.Effects {
		if se.Type != statusEffect.Type {
			ne = append(ne, se)
		}
	}
	sd.Effects = ne
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

type Target struct {
	Who     string `json:"who"`
	Targets string `json:"targets"`
}

func (t *Target) Sanitize() error {
	t.Who = strings.ToLower(t.Who)
	t.Targets = strings.ToLower(t.Targets)
	return nil
}
